package client

import "github.com/unluckythoughts/go-microservice/tools/web"

var (
	c web.Client
)

func Setup() {
	c = web.NewClient("http://localhost:5678")
}
