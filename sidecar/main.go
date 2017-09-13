package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const configLocation string = "/etc/example-app/config.json"

type config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("==> WARNING: Don't ever write secrets to logs.")
	log.Println("==>          This is for demonstration only.")
	log.Printf("Username: %s", config.Username)
	log.Printf("Password: %s", config.Password)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received shutting down gracefully...")
}

func getConfig() (*config, error) {
	for {
		content, err := ioutil.ReadFile(configLocation)
		switch {
		case os.IsNotExist(err):
			time.Sleep(5 * time.Millisecond)
			continue
		case err != nil:
			return nil, err
		}

		var config *config = &config{}
		err = json.Unmarshal(content, config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}
}
