package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	file, err := os.OpenFile("folder-cleaner.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error while opening log file", err)
	}
	defer file.Close()

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		log.Fatal("Please provide the path to the folder to clean up")
	}

	log.SetOutput(file)

	// Validate the path
	path := argsWithoutProg[0]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("Path does not exist: ", path)
	}

	t := time.Now()
	delete_counter := 0

	log.Println("Service triggered")

	c, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range c {
		info, err := entry.Info()
		if err != nil {
			log.Fatal(err)
		}

		absPath, err := filepath.Abs(filepath.Join(path, entry.Name()))
		if err != nil {
			log.Fatal(err)
		}

		diff := t.Sub(info.ModTime()).Hours()

		if diff > 24 {
			log.Println("Dir to Remove: ", absPath)
			e := os.RemoveAll(absPath)
			if e != nil {
				log.Fatal("Error while removing dir: ", e)
			} else {
				delete_counter++
			}
		}
	}
	if delete_counter > 0 {
		log.Println("Directories deleted: ", delete_counter)
	}
}
