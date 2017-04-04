package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Vault struct {
	URL       string
	Path      string
	LoginPath string
	User      string
	Password  string
}

type Account struct {
	Name     string
	User     string `json:"user"`
	Password string `json:"pass"`
}

func login(vault Vault) (string, error) {
	url := fmt.Sprintf("%s%s", vault.URL, vault.LoginPath)
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	token := res.Header.Get("Cyberarklogonresult")
	if token == "" {
		log.Fatal("Could not get token")
	}

	return token, nil
}

func getAccount(vault Vault, token string, account Account) (Account, error) {
	url := fmt.Sprintf("%s%s/%s", vault.URL, vault.Path, account.Name)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return account, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return account, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&account); err != nil {
		log.Println(err)
	}

	return account, nil
}

func outputAccount(account Account) {
	fmt.Println("User:     ", account.User)
	fmt.Println("Password: ", account.Password)
}
