package main

import (
	"github.com/IstarVin/manga-reader-go/intializer"
	"github.com/IstarVin/manga-reader-go/server"
	"github.com/IstarVin/manga-reader-go/syncmanager"
	"log"
)

func init() {
	syncmanager.Init(1)
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
