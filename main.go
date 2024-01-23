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

type Results struct {
	Results []Result `json:"results"`
}

type Result struct {
	Download Download `json:"download"`
}

type Download struct {
	Bandwidth int     `json:"bandwidth"`
	Bytes     int     `json:"bytes"`
	Elapsed   int     `json:"elapsed"`
	Latency   Latency `json:"latency"`
}

type Latency struct {
	Iqm    float32 `json:"iqm"`
	Low    float32 `json:"low"`
	High   float32 `json:"high"`
	Jitter float32 `json:"jitter"`
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
	err = os.MkdirAll(outFileDir, 0777)
	if err != nil {
		panicMessage := fmt.Sprintf("Error creating out file directory: %v", err)
		panic(panicMessage)
	}

	// Write results to out file
	err = os.WriteFile(outFilePath, finalResult, 0777)
	if err != nil {
		panicMessage := fmt.Sprintf("Error writing to file: %v", err)
		panic(panicMessage)
	}

	fmt.Println("Done!")
}
