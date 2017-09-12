package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/vault/api"
)

func main() {
	config := api.DefaultConfig()
	vaultClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Fatal(err)
	}

	s, err := vaultClient.Logical().Write("/auth/kube/login", map[string]interface{}{
		"role": "test",
		"jwt":  string(content[:]),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(s.Auth.ClientToken)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Keep token renewed
	renewer, err := vaultClient.NewRenewer(&api.RenewerInput{
		Secret: s,
		Grace:  1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	go renewer.Renew()
	defer renewer.Stop()

	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				log.Fatal(err)
			}
		case renewal := <-renewer.RenewCh():
			log.Printf("Successfully renewed: %#v", renewal)
		case <-quit:
			log.Fatal("Shutdown signal received, exiting...")
		}
	}

}
