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

	s, err := vaultClient.Logical().Write("/auth/kubernetes/login", map[string]interface{}{
		"role": "demo",
		"jwt":  string(content[:]),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("==> WARNING: Don't ever write secrets to logs.")
	log.Println("==>          This is for demonstration only.")
	log.Printf("Vault token: %s\n", s.Auth.ClientToken)

	vaultClient.SetToken(s.Auth.ClientToken)
	s, err = vaultClient.Logical().Read("/aws/creds/readonly")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Give some time for IAM cred creation to propagate since this action is
	// eventually consistent
	time.Sleep(30 * time.Second)

	accessKey := s.Data["access_key"].(string)
	secretKey := s.Data["secret_key"].(string)
	log.Println("==> WARNING: Don't ever write secrets to logs.")
	log.Println("==>          This is for demonstration only.")
	log.Printf("AWS Access Key: %s\n", accessKey)
	log.Printf("AWS Secret Key: %s\n", secretKey)

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

	log.Println("Starting renewal loop")
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
