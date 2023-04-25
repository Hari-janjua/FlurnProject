package main

import (
	"FlurnProject/Routes"
)

func main() {
	router := Routes.SetupRouter()
	router.Run()
}
