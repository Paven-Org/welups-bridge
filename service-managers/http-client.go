package manager

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type HttpClient struct {
	client             *http.Client
	baseURL            string
	authorizationToken string
}

func MkHttpClient(baseURL string, authToken string) (*HttpClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	return &HttpClient{
			client:             client,
			baseURL:            baseURL,
			authorizationToken: authToken,
		},
		nil
}

func (cli *HttpClient) SetToken(token string) {
	cli.authorizationToken = token
}

func (cli *HttpClient) CloseIdleConnections() {
	cli.client.CloseIdleConnections()
}

func (cli *HttpClient) Get(route string) (resp *http.Response, err error) {
	reqURL := fmt.Sprintf("%s%s", cli.baseURL, route)
	request, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cli.authorizationToken))

	return cli.client.Do(request)
}

func (cli *HttpClient) Post(route, contentType string, body string) (resp *http.Response, err error) {
	reqURL := fmt.Sprintf("%s%s", cli.baseURL, route)
	request, err := http.NewRequest("POST", reqURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.authorizationToken))
	request.Header.Set("Content-Type", contentType)

	return cli.client.Do(request)
}

func (cli *HttpClient) PostJSON(route, body string) (*http.Response, error) {
	return cli.Post(route, "application/json", body)
}
