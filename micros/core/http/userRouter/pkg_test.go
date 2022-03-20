package userRouter

import (
	"bridge/micros/core/config"
	"bridge/micros/core/model"
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

func login(t *testing.T, username string, password string) {
	resp, err := cli.PostJSON("/v1/u/login", fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
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

// re-logout logged out session should fail
func TestDoubleLogout(t *testing.T) {
	login(t, "root", "root")
	logout(t)
	logout(t)
}

func TestLoginLogout(t *testing.T) {
	login(t, "root", "root")
	logout(t)
}

func TestChangePasswd(t *testing.T) {
	login(t, "root", "root")
	passwd(t, "root", "moot")
	passwd(t, "moot", "root")
	logout(t)
}

func passwd(t *testing.T, oldpass string, newpass string) {
	resp, err := cli.PostJSON("/v1/u/passwd", fmt.Sprintf(`{"old_passwd": "%s", "new_passwd": "%s"}`, oldpass, newpass))
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
}

func TestGetUser(t *testing.T) {
	getuser(t, "root")
}

func getuser(t *testing.T, username string) {
	resp, err := cli.Get("/v1/u/root")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
	var user model.User
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	if err := json.Unmarshal(res, &user); err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	fmt.Printf("Got user: %+v\n", user)
}

func TestGetRoles(t *testing.T) {
	login(t, "root", "root")
	getroles(t)
	getusers(t)
	getuserswithrole(t)
	logout(t)
}

func getroles(t *testing.T) {
	resp, err := cli.Get("/v1/a/m/u/roles")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
	var roles []string
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	if err := json.Unmarshal(res, &roles); err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	fmt.Printf("Got roles: %+v\n", roles)
}

func getusers(t *testing.T) {
	resp, err := cli.Get("/v1/a/m/u/users/1?limit=1")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
	var users []model.User
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	if err := json.Unmarshal(res, &users); err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	fmt.Printf("Got users: %v\n", users)
}

func getuserswithrole(t *testing.T) {
	resp, err := cli.Get("/v1/a/m/u/haverole/service/1?limit=1")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("Response: ", resp)
	var users []model.User
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	if err := json.Unmarshal(res, &users); err != nil {
		t.Fatal("Invalid response, error: ", err.Error())
	}

	fmt.Printf("Got users: %v\n", users)
}
