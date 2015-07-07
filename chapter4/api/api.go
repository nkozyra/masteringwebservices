package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var Database *sql.DB
var Routes *mux.Router
var Format string

type Count struct {
	DBCount int
}

type UpdateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

func Init() {
	Routes = mux.NewRouter()
	Routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	Routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	Routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
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

	Format = r.URL.Query()["format"][0]

}

func SetFormat(data interface{}) []byte {

	var apiOutput []byte
	if Format == "json" {
		output, _ := json.Marshal(data)
		apiOutput = output
	} else if Format == "xml" {
		output, _ := xml.Marshal(data)
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

	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	// SQL Injection here, keep reading!
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'"
	q, err := Database.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(q)
}
func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting retrieval")
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

func StartServer() {

	db, err := sql.Open("mysql", "root@/social_network")
	if err != nil {

	}
	Database = db

	http.Handle("/", Routes)
	http.ListenAndServe(":8080", nil)
}
