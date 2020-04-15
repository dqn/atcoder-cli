package atcoder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const baseURL = "https://atcoder.jp"

type AtCoderClient struct {
	client *http.Client
}

func New() *AtCoderClient {
	return &AtCoderClient{client: &http.Client{}}
}

func (a *AtCoderClient) newRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest("POST", baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.URL.Path = path
	return req, nil
}

func (a *AtCoderClient) getCSRFToken() (string, error) {
	resp, err := a.client.Get(fmt.Sprintf("%s/login", baseURL))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := string(body)
	q := `var csrfToken = "`
	index := strings.LastIndex(s, q)
	if index == -1 {
		return "", fmt.Errorf("cannot find csrf token")
	}

	b := make([]byte, 0, 48)
	for i := index + len(q); s[i] != byte('"'); i++ {
		b = append(b, s[i])
	}
	return string(b), nil
}

func (a *AtCoderClient) Login(username, password string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	a.client.Jar = jar

	csrfToken, err := a.getCSRFToken()
	if err != nil {
		return err
	}

	resp, err := a.client.PostForm(
		fmt.Sprintf("%s/login", baseURL),
		url.Values{
			"username":   {username},
			"password":   {password},
			"csrf_token": {csrfToken},
		},
	)
	if err != nil {
		return err
	}

	if resp.Request.URL.Path != "/home" {
		return fmt.Errorf("failed to login")
	}

	return nil
}
