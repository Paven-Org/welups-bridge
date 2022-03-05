package manager

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestHttpClient(t *testing.T) {
	cli, _ := MkHttpClient("https://google.com", "")

	query := "/search?q=rick+rolled"
	resp, err := cli.Get(query)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Printf("status: %s\n", resp.Status)
	fmt.Printf("response: %s\n", string(res))

}
