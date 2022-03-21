package main

import (
	"landing-page/api"
)

func main() {
	apiServer := api.NewServer()
	apiServer.RunServer()
}
