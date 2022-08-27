package admRouter

import (
	"bridge/micros/core/config"
	manager "bridge/service-managers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

var cli *manager.HttpClient

func TestMain(m *testing.M) {
	var err error
	config.Load()
	//log.Init(config.Get().Structured)
	cli, err = manager.MkHttpClient("https://localhost:8001", "")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	m.Run()
}

type loginResp struct {
	Token string
}

func login(t *testing.T) {
	resp, err := cli.PostJSON("/v1/u/login", `{"username": "root", "password": "root"}`)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	//cli.SetToken(res)
	fmt.Printf("status: %s\n", resp.Status)
	fmt.Printf("response: %s\n", string(res))

	liResp := loginResp{}
	if err = json.Unmarshal(res, &liResp); err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Printf("token: %s\n", liResp.Token)
	cli.SetToken(liResp.Token)

	var jar = cli.GetCookieJar()
	fmt.Printf("cookie jar: %v\n\n", jar)
}

func logout(t *testing.T) {
	resp, err := cli.PostJSON("/v1/u/logout", "")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)

	cli.SetToken("")
}

// actual test
func TestPing(t *testing.T) {
	login(t)
	fmt.Printf("\n\n Test ping...\n")
	resp, err := cli.Get("/v1/a/ping")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
	fmt.Printf("\n Test ping done\n\n")
	logout(t)
}

// actual test
func TestGetUsers(t *testing.T) {
	login(t)
	fmt.Printf("\n\n Test GetUsers...\n")
	resp, err := cli.Get("/v1/a/m/u/users/1?limit=20")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	bod, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", bod)
	fmt.Printf("\n Test GetUsers done\n\n")
	logout(t)
}

// actual test
func TestUnsetAuthenticator(t *testing.T) {
	login(t)
	fmt.Printf("\n\n Test GetUsers...\n")
	resp, err := cli.PostJSON("/v1/a/m/wel/unset/authenticator-prikey", "")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	bod, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", bod)
	fmt.Printf("\n Test GetUsers done\n\n")
	logout(t)
}
