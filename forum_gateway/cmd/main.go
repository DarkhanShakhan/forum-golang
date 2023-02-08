package main

import (
	"forum_gateway/internal/app"
)

func init() {
	app.SetEnv()
}
func main() {
	app.Run()
}
