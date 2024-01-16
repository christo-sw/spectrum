package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ConfigFile struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func main() {
	fmt.Println("Doing speed test...")

	// Read config file
	contents, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Parse config file
	var config ConfigFile
	err = json.Unmarshal(contents, &config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	fmt.Println("Config file parsed:", config)

	var outFileDir = "./out/" + time.Now().Format(time.DateOnly)
	var outFilePath = outFileDir + "/test.json"
	var finalResult []byte

	// Do a test for each server
	for _, server := range config.Servers {
		fmt.Println("Testing", server.Name)
		testResult, err := exec.Command("speedtest", "-f", "json", "-s", fmt.Sprintf("%v", server.Id)).Output()
		if err != nil {
			fmt.Printf("Error running speed test for %v: %v\n", server.Name, err)
			continue
		}

		finalResult = append(finalResult, testResult...)
	}

	// Create out file directory if it does not exist
	err = os.Mkdir(outFileDir, 0777)
	if err != nil {
		fmt.Println("Error creating out file directory:", err)
	}

	// Write results to out file
	err = os.WriteFile(outFilePath, finalResult, 0777)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Println("Done!")
}
