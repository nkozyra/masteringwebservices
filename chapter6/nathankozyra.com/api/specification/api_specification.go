package specification

type MethodPOST struct {
	POST EndPoint
}

type MethodGET struct {
	GET EndPoint
}

type MethodPUT struct {
	PUT EndPoint
}

type MethodOPTIONS struct {
	OPTIONS EndPoint
}

type EndPoint struct {
	Description string  `json:"description"`
	Parameters  []Param `json:"parameters"`
}

type Param struct {
	Name             string "json:name"
	ParameterDetails Detail `json:"details"`
}

type Detail struct {
	Type        string "json:type"
	Description string `json:"description"`
	Required    bool   "json:required"
}

var UserOPTIONS = MethodOPTIONS{OPTIONS: EndPoint{Description: "This page"}}
var UserPostParameters = []Param{{Name: "Email", ParameterDetails: Detail{Type: "string", Description: "A new user's email address", Required: false}}}

var UserPOST = MethodPOST{POST: EndPoint{Description: "Create a user", Parameters: UserPostParameters}}
var UserGET = MethodGET{GET: EndPoint{Description: "Access a user"}}
