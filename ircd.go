package main

import (
	"log"

	"./src/ircd"
)

func main() {
	log.Printf("Hello, World!")

	//config file
	config := ircd.LoadConfiguration("/location/of/config.json")

	server := ircd.Server(config)
	server.Run()
}
