package main

import (
	"autograd/server/router"
)

func main() {
	e := router.New()

	e.Logger.Fatal(e.Start(":1323"))
}
