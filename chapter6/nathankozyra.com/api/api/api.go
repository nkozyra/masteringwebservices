package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"

	"log"
	"net/http"
	"net/url"

	"strconv"
	"strings"

	Password "nathankozyra.com/api/password"
	Pseudoauth "nathankozyra.com/api/pseudoauth"
	OauthServices "nathankozyra.com/api/services"
	Documentation "nathankozyra.com/api/specification"
	"sync"
	"time"
)

var Database *sql.DB
var Routes *mux.Router
var Format string

type UserSession struct {
	ID              string
	GorillaSesssion *sessions.Session
	UID             int
	Expire          time.Time
}

var Session UserSession

func (us *UserSession) Create() {
	us.ID = Password.GenerateSessionID(32)
}

const serverName = "localhost"
const SSLport = ":443"
const HTTPport = ":8081"
const SSLprotocol = "https://"
const HTTPprotocol = "http://"

var PermittedDomains []string

type Count struct {
	DBCount int
}

type UpdateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type CreateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID       int    "json:id"
	Name     string "json:username"
	Email    string "json:email"
	First    string "json:first"
	Last     string "json:last"
	Password string "json:password"
	Salt     string "json:salt"
	Hash     string "json:hash"
}

type UserDocumentation struct {
}

type OauthAccessResponse struct {
	AccessToken string `json:"access_key"`
}

func Init(allowedDomains []string) {
	for _, domain := range allowedDomains {
		PermittedDomains = append(PermittedDomains, domain)
	}
	Routes = mux.NewRouter()
	Routes.HandleFunc("/interface", APIInterface).Methods("GET", "POST", "PUT", "UPDATE")
	Routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	Routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	Routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
	Routes.HandleFunc("/api/users", UsersInfo).Methods("OPTIONS")
	Routes.HandleFunc("/api/statuses", StatusCreate).Methods("POST")
	Routes.HandleFunc("/api/statuses", StatusRetrieve).Methods("GET")
	Routes.HandleFunc("/api/statuses/{id:[0-9]+}", StatusUpdate).Methods("PUT")
	Routes.HandleFunc("/api/statuses/{id:[0-9]+}", StatusDelete).Methods("DELETE")
	Routes.HandleFunc("/authorize", ApplicationAuthorize).Methods("POST")
	Routes.HandleFunc("/authorize", ApplicationAuthenticate).Methods("GET")
	Routes.HandleFunc("/authorize/{service:[a-z]+}", ServiceAuthorize).Methods("GET")
	Routes.HandleFunc("/connect/{service:[a-z]+}", ServiceConnect).Methods("GET")
	Routes.HandleFunc("/oauth/token", CheckCredentials).Methods("POST")
}

type Page struct {
	Title        string
	Authorize    bool
	Authenticate bool
	Application  string
	Action       string
	ConsumerKey  string
	Redirect     string
	PageType     string
}

func CheckLogin(w http.ResponseWriter, r *http.Request) bool {
	cookieSession, err := r.Cookie("sessionid")
	if err != nil {
		fmt.Println("no such cookie")
		Session.Create()
		fmt.Println(Session.ID)
		currTime := time.Now()
		Session.Expire = currTime.Local()
		Session.Expire.Add(time.Hour)

		return false
	} else {
		fmt.Println("found cookki")
		tmpSession := UserSession{UID: 0}
		loggedIn := Database.QueryRow("select user_id from sessions where session_id=?", cookieSession).Scan(&tmpSession.UID)
		if loggedIn != nil {
			return false
		} else {
			if tmpSession.UID == 0 {
				return false
			} else {

				return true
			}
		}
	}
	return false
}

func ServiceAuthorize(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	service := params["service"]

	loggedIn := CheckLogin(w, r)
	if loggedIn == false {
		Cookie := http.Cookie{Name: "sessionid", Value: Session.ID, Expires: Session.Expire}
		fmt.Println("Setting cookie!")
		http.SetCookie(w, &Cookie)
		redirect := url.QueryEscape("/authorize/" + service)
		http.Redirect(w, r, "/authorize?redirect="+redirect, http.StatusUnauthorized)
		return
	}

	redURL := OauthServices.GetAccessTokenURL(service, "")
	http.Redirect(w, r, redURL, http.StatusFound)

}

func ServiceConnect(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	fmt.Println(code)
}

func ApplicationAuthorize(w http.ResponseWriter, r *http.Request) {

	CheckLogin(w, r)

	username := r.FormValue("username")
	password := r.FormValue("password")
	allow := r.FormValue("authorize")
	authType := r.FormValue("auth_type")
	redirect := r.FormValue("redirect")

	var dbPassword string
	var dbSalt string
	var dbUID string

	uerr := Database.QueryRow("SELECT user_password, user_salt, user_id from users where user_nickname=?", username).Scan(&dbPassword, &dbSalt, &dbUID)
	if uerr != nil {

	}

	consumerKey := r.FormValue("consumer_key")

	var CallbackURL string
	var appUID string
	if authType == "client" {
		err := Database.QueryRow("SELECT user_id,callback_url from api_credentials where consumer_key=?", consumerKey).Scan(&appUID, &CallbackURL)
		if err != nil {
			fmt.Println("SELECT user_id,callback_url from api_credentials where consumer_key=?", consumerKey)
			fmt.Println(err.Error())
			return
		}
	}

	expectedPassword := Password.GenerateHash(dbSalt, password)
	fmt.Println("allow:", allow)
	fmt.Println("authtype:", authType)
	fmt.Println(dbPassword, "=", expectedPassword)
	if (dbPassword == expectedPassword) && (allow == "1") && (authType == "consumer") {
		fmt.Println("Yes!")
		requestToken := Pseudoauth.GenerateToken()

		authorizeSQL := "INSERT INTO api_tokens set application_user_id=?, user_id=?, api_token_key=?"

		q, connectErr := Database.Exec(authorizeSQL, appUID, dbUID, requestToken)
		if connectErr != nil {
			fmt.Println(connectErr.Error())
		} else {
			fmt.Println(q)
		}
		redirectURL := CallbackURL + "?request_token=" + requestToken
		fmt.Println(redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusAccepted)

	} else if (dbPassword == expectedPassword) && authType == "user" {

		_, err := Database.Exec("insert into sessions set session_id=?,user_id=?", Session.ID, dbUID)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("redirecting")
		http.Redirect(w, r, redirect, http.StatusOK)
	} else {
		fmt.Println(authType)
		fmt.Println(dbPassword, expectedPassword)
		http.Redirect(w, r, "/authorize", http.StatusUnauthorized)
	}

}

func ApplicationAuthenticate(w http.ResponseWriter, r *http.Request) {

	Authorize := Page{}
	Authorize.Authenticate = true
	Authorize.Title = "Login"
	Authorize.Application = ""
	Authorize.Action = "/authorize"
	if len(r.URL.Query()["consumer_key"]) > 0 {
		Authorize.ConsumerKey = r.URL.Query()["consumer_key"][0]
	} else {
		Authorize.ConsumerKey = ""
	}
	if len(r.URL.Query()["redirect"]) > 0 {
		Authorize.Redirect = r.URL.Query()["redirect"][0]
	} else {
		Authorize.Redirect = ""
	}

	if Authorize.ConsumerKey == "" && Authorize.Redirect != "" {
		Authorize.PageType = "user"
	} else {
		Authorize.PageType = "consumer"
	}

	tpl := template.Must(template.New("main").ParseFiles("authorize.html"))
	tpl.ExecuteTemplate(w, "authorize.html", Authorize)
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirecting.")
	redirectURL := SSLprotocol + serverName + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func StatusDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Nothing to see here")
}

func StatusUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Coming soon to an API near you!")
}

func ValidateUserRequest(cKey string, cToken string) string {
	var UID string
	var aUID string
	var appUID string
	Database.QueryRow("SELECT at.user_id,at.application_user_id,ac.user_id as appuser from api_tokens at left join api_credentials ac on ac.user_id=at.application_user_id where api_token_key=?", cToken).Scan(&UID, &aUID, &appUID)

	return appUID
}

func StatusCreate(w http.ResponseWriter, r *http.Request) {

	Response := CreateResponse{}
	UserID := r.FormValue("user")
	Status := r.FormValue("status")
	Token := r.FormValue("token")
	ConsumerKey := r.FormValue("consumer_key")

	vUID := ValidateUserRequest(ConsumerKey, Token)
	if vUID != UserID {
		Response.Error = "Invalid user"
		http.Error(w, Response.Error, 401)
		//fmt.Println(w, SetFormat(Response))
	} else {
		_, inErr := Database.Exec("INSERT INTO users_status set user_status_text=?, user_id=?", Status, UserID)
		if inErr != nil {
			fmt.Println(inErr.Error())
			Response.Error = "Error creating status"
			http.Error(w, Response.Error, 500)
			fmt.Fprintln(w, Response)
		} else {
			Response.Error = "Status created"
			fmt.Fprintln(w, Response)
		}
	}

}

func StatusRetrieve(w http.ResponseWriter, r *http.Request) {

}

func ErrorMessages(err int64) (int, int, string) {
	errorMessage := ""
	statusCode := 200
	errorCode := 0
	switch err {
	case 1062:
		errorMessage = http.StatusText(409)
		errorCode = 10
		statusCode = http.StatusConflict
	default:
		errorMessage = http.StatusText(int(err))
		errorCode = 0
		statusCode = int(err)
	}

	return errorCode, statusCode, errorMessage

}

func GetFormat(r *http.Request) {

	if len(r.URL.Query()["format"]) > 0 {
		Format = r.URL.Query()["format"][0]
	} else {
		Format = "json"
	}
}

func SetFormat(data interface{}) []byte {

	var apiOutput []byte
	if Format == "json" {
		output, _ := json.Marshal(data)
		apiOutput = output
	} else if Format == "xml" {
		output, _ := xml.Marshal(data)
		apiOutput = output
	} else {
		output, _ := json.Marshal(data)
		apiOutput = output
	}
	return apiOutput
}

func dbErrorParse(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}

type DocMethod interface {
}

func UsersInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "DELETE,GET,HEAD,OPTIONS,POST,PUT")

	UserDocumentation := []DocMethod{}
	UserDocumentation = append(UserDocumentation, Documentation.UserPOST)
	UserDocumentation = append(UserDocumentation, Documentation.UserOPTIONS)

	output := SetFormat(UserDocumentation)
	fmt.Fprintln(w, string(output))
}

func APIInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "interface.html")
}

func UsersUpdate(w http.ResponseWriter, r *http.Request) {
	Response := UpdateResponse{}
	params := mux.Vars(r)
	uid := params["id"]
	email := r.FormValue("email")

	var userCount int

	err := Database.QueryRow("SELECT count(user_id) from users where user_id=?", uid).Scan(&userCount)
	if userCount == 0 {

		error, httpCode, msg := ErrorMessages(404)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)

	} else if err != nil {

	} else {

		_, uperr := Database.Exec("UPDATE users set user_email=? where user_id=?", email, uid)
		if uperr != nil {
			_, errorCode := dbErrorParse(uperr.Error())
			_, httpCode, msg := ErrorMessages(errorCode)

			Response.Error = msg
			Response.ErrorCode = httpCode
			http.Error(w, msg, httpCode)
		} else {
			Response.Error = "success"
			Response.ErrorCode = 0
			output := SetFormat(Response)
			fmt.Fprintln(w, string(output))
		}
	}

}

func UserCreate(w http.ResponseWriter, r *http.Request) {

	for _, domain := range PermittedDomains {
		fmt.Println("allowing", domain)
		w.Header().Set("Access-Control-Allow-Origin", domain)
	}

	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	NewUser.Password = r.FormValue("password")
	salt, hash := Password.ReturnPassword(NewUser.Password)
	fmt.Println(salt, hash)
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	Response := CreateResponse{}
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'" + ", user_password='" + hash + "', user_salt='" + salt + "'"
	q, err := Database.Exec(sql)
	if err != nil {
		errorMessage, errorCode := dbErrorParse(err.Error())
		fmt.Println(errorMessage)
		error, httpCode, msg := ErrorMessages(errorCode)
		Response.Error = msg
		Response.ErrorCode = error
		http.Error(w, "Conflict", httpCode)
	}

	fmt.Println(q)
}
func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting retrieval")

	accessToken := r.FormValue("access_token")
	if accessToken == "" || CheckToken(accessToken) == false {

	}

	GetFormat(r)
	start := 0
	limit := 10

	next := start + limit

	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:8080/api/users?start="+string(next)+"; rel=\"next\"")

	rows, _ := Database.Query("select user_id, user_nickname, user_first, user_last, user_email from users LIMIT 10")
	Response := Users{}

	for rows.Next() {

		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)
		fmt.Println(user)
		Response.Users = append(Response.Users, user)
	}

	output := SetFormat(Response)
	fmt.Fprintln(w, string(output))
}

func CheckToken(token string) bool {
	return true
}

func CheckCredentials(w http.ResponseWriter, r *http.Request) {
	var Credentials string
	Response := CreateResponse{}
	consumerKey := r.FormValue("consumer_key")
	fmt.Println(consumerKey)
	timestamp := r.FormValue("timestamp")
	signature := r.FormValue("signature")
	nonce := r.FormValue("nonce")
	err := Database.QueryRow("SELECT consumer_secret from api_credentials where consumer_key=?", consumerKey).Scan(&Credentials)
	if err != nil {
		error, httpCode, msg := ErrorMessages(404)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)
		return
	}

	token, err := Pseudoauth.ValidateSignature(consumerKey, Credentials, timestamp, nonce, signature, 0)
	if err != nil {
		error, httpCode, msg := ErrorMessages(401)
		log.Println(error)
		log.Println(w, msg, httpCode)
		Response.Error = msg
		Response.ErrorCode = httpCode
		http.Error(w, msg, httpCode)
		return
	}
	fmt.Println(token)
	AccessRequest := OauthAccessResponse{}
	AccessRequest.AccessToken = token.AccessToken
	output := SetFormat(AccessRequest)
	fmt.Fprintln(w, string(output))
}

func StartServer() {
	OauthServices.InitServices()
	fmt.Println(Password.GenerateSalt(22))
	fmt.Println(Password.GenerateSalt(41))

	db, err := sql.Open("mysql", "root@/social_network")
	if err != nil {

	}
	Database = db

	wg := sync.WaitGroup{}

	log.Println("Starting redirection server, try to access @ http:")

	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem", Routes)
		//http.ListenAndServe(SSLport,http.HandlerFunc(secureRequest))
		wg.Done()
	}()

	wg.Wait()
}
