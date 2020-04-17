package atcoder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

const baseURL = "https://atcoder.jp"

type AtCoderClient struct {
	config    *Config
	client    *http.Client
	csrfToken string
}

type Command struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

type Config struct {
	Username     string    `yaml:"username"`
	Password     string    `yaml:"password"`
	TemplatePath string    `yaml:"template_path"`
	Test         Command   `yaml:"test"`
	Pretest      []Command `yaml:"pretest"`
	Posttest     []Command `yaml:"posttest"`
}

type test struct {
	input  string
	output string
}

func New(config *Config) *AtCoderClient {
	return &AtCoderClient{
		client: &http.Client{},
		config: config,
	}
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

func (a *AtCoderClient) Login() error {
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
			"username":   {a.config.Username},
			"password":   {a.config.Password},
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

func (a *AtCoderClient) Init(url string) error {
	resp, err := a.client.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	s := string(body)
	contest, problem := parseContest(s)

	dir := fmt.Sprintf("atcoder/%s", contest)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	tests := parseTests(s)

	ch := make(chan error, 2)
	go func() {
		var wg sync.WaitGroup
		wg.Add(cap(ch))
		go func() {
			ch <- createTestFile(fmt.Sprintf("%s/%s.txt", dir, problem), tests)
			wg.Done()
		}()
		go func() {
			templatePath := a.config.TemplatePath
			ch <- createSourceFile(dir, problem, templatePath)
			wg.Done()
		}()
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AtCoderClient) Test(contest, problem string) (bool, error) {
	path := fmt.Sprintf("./atcoder/%s/%s", contest, problem)
	ext := ".cpp" // TODO
	replacements := []*replacement{
		{
			regexp.MustCompile(`{{\s*file_path\s*}}`),
			[]byte(path + ext),
		},
	}

	tests, err := readTests(path + ".txt")
	if err != nil {
		return false, err
	}

	for _, v := range a.config.Pretest {
		if _, err := execCommand(v, replacements).Output(); err != nil {
			return false, err
		}
	}

	for _, v := range tests {
		cmd, err := execCommandWithStdin(a.config.Test, replacements, v.input)
		if err != nil {
			return false, err
		}

		b, err := cmd.Output()
		if err != nil {
			return false, err
		}

		ans := strings.TrimSpace(string(b))
		if ans != v.output {
			println(ans, v.output)
			return false, nil
		}
	}

	for _, v := range a.config.Posttest {
		if _, err := execCommand(v, replacements).Output(); err != nil {
			return false, err
		}
	}
	return true, nil
}
