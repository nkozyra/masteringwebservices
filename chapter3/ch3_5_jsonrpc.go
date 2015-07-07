package main

/*test with curl
curl -X POST -H "Content-Type: application/json" \
-d '{"method":"stringService.Length","params":[{"Who":"Test"}], "id":"1"}' \
http://localhost:10000/rpc
*/

import (
	"fmt"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type RPCAPIArguments struct {
	Message string
}

type RPCAPIResponse struct {
	Message string
}

func getLength(message string) string {
	length := utf8.RuneCountInString(message)
	return strconv.FormatInt(int64(length), 10)
}

type StringService struct{}

func (h *StringService) Length(r *http.Request, arguments *RPCAPIArguments, reply *RPCAPIResponse) error {
	reply.Message = "Your string is " + getLength(arguments.Message) + " characters long"
	return nil
}

func main() {
	fmt.Println("Starting service")
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(StringService), "")
	http.Handle("/rpc", s)
	http.ListenAndServe(":10000", nil)
}
