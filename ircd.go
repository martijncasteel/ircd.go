package main

import (
	"flag"
	"log"

	"./src/ircd"
)

func main() {
	file := flag.String("config", "config.yaml", "configuration file")
	flag.Parse()

	log.Printf("Loading %s ..", *file)
	config, err := ircd.LoadConfiguration(*file)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", config)

	server := ircd.Server(config)
	server.Run()
}
