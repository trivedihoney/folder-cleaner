package main

import (
	"flag"
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
	log.SetOutput(file)

	log.Println("Service triggered")

	max_age_hours_ptr := flag.Int("max-age-hours", -1, "Files older than this will be deleted")
	flag.Parse()

	if *max_age_hours_ptr == -1 {
		log.Fatal("Please provide the max-age-hours flag")
	} else if *max_age_hours_ptr <= 0 {
		log.Fatal("max-age-hours should be greater than 0")
	}

	max_age_hours := *max_age_hours_ptr

	argsWithoutProg := flag.Args()

	if len(argsWithoutProg) == 0 {
		log.Fatal("Please provide the path to the folder to clean up")
	}

	// Validate the path
	path := argsWithoutProg[0]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("Path does not exist: ", path)
	}

	t := time.Now()
	delete_counter := 0

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

		if diff > float64(max_age_hours) {
			e := os.RemoveAll(absPath)
			if e != nil {
				log.Fatal("Error while removing", absPath, e)
			} else {
				delete_counter++
			}
		}
	}
	if delete_counter > 0 {
		log.Println("Files & Folders deleted: ", delete_counter)
	}
}
