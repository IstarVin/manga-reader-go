package intializer

import (
	"flag"
	"github.com/IstarVin/manga-reader-go/global"
	"log"
	"os"
	"path/filepath"
)

func LoadCLIArguments() {
	// Get Current Working Directory
	var err error
	global.DefaultConfig.DataDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Error getting the current working directory")
	}

	var config global.Configuration

	// Set CLI Arguments
	flag.StringVar(&config.DataDir, "dir", global.DefaultConfig.DataDir, "Set Config Directory")
	flag.StringVar(&config.Host, "host", global.DefaultConfig.Host, "Set the host of server")
	flag.IntVar(&config.Port, "port", global.DefaultConfig.Port, "Set the port of the server")
	flag.Parse()
	global.Config = &config

	// Set Global Variables
	global.ConfigFilePath = filepath.Join(global.Config.DataDir, "global.json")
	global.MangasDirectory = filepath.Join(global.Config.DataDir, "local")
	global.MangaDatabasePath = filepath.Join(global.Config.DataDir, "database", "mangaDB.json")
	global.CategoryDatabasePath = filepath.Join(global.Config.DataDir, "database", "categoryDB.json")

	// Create Directories
	err = os.MkdirAll(global.MangasDirectory, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating the mangas directory")
	}
	err = os.MkdirAll(filepath.Join(global.Config.DataDir, "database"), os.ModePerm)
	if err != nil {
		log.Fatal("Error creating the database directory")
	}
}
