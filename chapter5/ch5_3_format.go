package main

type GenericData struct {
	Options GenericDataBlock
}

type GenericDataBlock struct {
	Server  string
	Address string
}

func main() {
	Data := GenericData{Options: GenericDataBlock{Server: "server01", Address: "127.0.0.1"}}

}
