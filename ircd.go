package main

import (
	"log"
	"os"

	"./src/ircd"
)

func main() {
	log.Printf("Loading configuration ..")

	// load configuration from file args[1]
	config, err := ircd.LoadConfiguration(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", config)

	server := ircd.Server(config)
	server.Run()
}
