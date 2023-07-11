package main

import (
	"github.com/IstarVin/manga-reader-go/intializer"
	"github.com/IstarVin/manga-reader-go/server"
	"log"
)

func init() {
	intializer.LoadCLIArguments()
	intializer.LoadConfigFile()
	intializer.LoadDatabase()
}

func main() {
	err := server.NewServer().Run()
	if err != nil {
		log.Fatal("Error running the server")
	}
}
