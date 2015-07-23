package v1

import (
	"nathankozyra.com/api/api"
)

func API() {
	api.Init([]string{"http://www.example.com", "http://www.mastergoco.com"})
	api.StartServer()
}
