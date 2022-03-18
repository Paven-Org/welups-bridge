package publicRouter

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
	cli, err = manager.MkHttpClient("http://localhost:8001", "")
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
	resp, err := cli.PostJSON("/v1/p/login", `{"username": "root", "password": "root"}`)
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
	resp, err := cli.PostJSON("/v1/p/logout", "")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)

	cli.SetToken("")
}

// re-logout logged out session should fail
func TestDoubleLogout(t *testing.T) {
	login(t)
	logout(t)
	logout(t)
}

func TestLoginLogout(t *testing.T) {
	login(t)
	logout(t)
}
