package test

import (
	"encoding/json"
	"fmt"
	"hr-api/pkg/keyvault"
	"log"
	"testing"
)

func TestKeyvault(t *testing.T) {
	manager, err := keyvault.NewConfigLoader()
	if err != nil {
		log.Fatal(err)
	}

	secret, err := manager.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	s, _ := json.Marshal(secret)
	fmt.Println(string(s))
}
