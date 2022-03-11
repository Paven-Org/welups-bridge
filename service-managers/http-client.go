package manager

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type SimpleJar struct {
	lock    sync.RWMutex
	cookies map[string](*http.Cookie)
}

func (j *SimpleJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if len(cookies) < 1 {
		return
	}
	fmt.Println("[SetCookies] url: ", *u)
	host, _, _ := net.SplitHostPort(u.Host)
	fmt.Println("[SetCookies] host: ", host)

	j.lock.Lock()
	defer j.lock.Unlock()

	// only keep the last cookie
	j.cookies[host] = cookies[len(cookies)-1]

	fmt.Println("[SetCookies] ok, result;", j.cookies[host])
}

func (j *SimpleJar) Cookies(u *url.URL) []*http.Cookie {
	fmt.Println("[GetCookies] url: ", *u)
	host, _, _ := net.SplitHostPort(u.Host)
	fmt.Println("[GetCookies] host: ", host)

	j.lock.RLock()
	defer j.lock.RUnlock()
	res, ok := j.cookies[host]
	if !ok {
		fmt.Println("[GetCookies] not ok, result;", res)
		return []*http.Cookie{}
	}
	fmt.Println("[GetCookies] ok, result;", res)

	return []*http.Cookie{res}
}

func (j *SimpleJar) GetAllCookies() map[string]*http.Cookie {
	return j.cookies
}

func MkSimpleJar() *SimpleJar {
	return &SimpleJar{
		lock:    sync.RWMutex{},
		cookies: make(map[string]*http.Cookie, 20),
	}
}

type HttpClient struct {
	client             *http.Client
	baseURL            string
	authorizationToken string
}

func MkHttpClient(baseURL string, authToken string) (*HttpClient, error) {
	jar := MkSimpleJar()

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

func (cli *HttpClient) GetCookieJar() http.CookieJar {
	return cli.client.Jar
}
