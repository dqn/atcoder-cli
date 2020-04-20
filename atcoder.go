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
	Extension    string    `yaml:"extension"`
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
	a.csrfToken = csrfToken

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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	contest, task := parseContest(b)

	dir := fmt.Sprintf("atcoder/%s", contest)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	tests := parseTests(b)

	ch := make(chan error, 2)
	go func() {
		var wg sync.WaitGroup
		wg.Add(cap(ch))
		go func() {
			ch <- createTestFile(fmt.Sprintf("%s/%s.txt", dir, task), tests)
			wg.Done()
		}()
		go func() {
			dstPath := fmt.Sprintf("%s/%s%s", dir, task, a.config.Extension)
			ch <- copyFile(a.config.TemplatePath, dstPath)
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

func (a *AtCoderClient) Test(contest, task string) (bool, error) {
	path := fmt.Sprintf("./atcoder/%s/%s", contest, task)
	replacements := []*replacement{
		{
			regexp.MustCompile(`{{\s*file_path\s*}}`),
			[]byte(path + a.config.Extension),
		},
	}

	tests, err := readTests(path + ".txt")
	if err != nil {
		return false, err
	}

	for _, v := range a.config.Pretest {
		if b, err := execCommand(v, replacements).CombinedOutput(); err != nil {
			fmt.Print(string(b))
			return false, err
		}
	}

	ok := true
	for i, v := range tests {
		cmd, err := execCommandWithStdin(a.config.Test, replacements, v.input)
		if err != nil {
			return false, err
		}

		b, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Print(string(b))
			return false, err
		}

		ans := strings.TrimSpace(string(b))
		if ans == v.output {
			printAC(i)
		} else {
			printWA(i, v.output, ans)
			ok = false
		}
	}
	printDivider()

	for _, v := range a.config.Posttest {
		if b, err := execCommand(v, replacements).CombinedOutput(); err != nil {
			fmt.Print(string(b))
			return ok, err
		}
	}
	return ok, nil
}

func (a *AtCoderClient) Submit(contest, task string) (bool, error) {
	path := fmt.Sprintf("./atcoder/%s/%s%s", contest, task, a.config.Extension)
	sourceCode, err := readFileContent(path)
	if err != nil {
		return false, err
	}
	resp, err := a.client.PostForm(
		fmt.Sprintf("%s/contests/%s/submit", baseURL, contest),
		url.Values{
			"data.TaskScreenName": {fmt.Sprintf("%s_%s", contest, task)},
			"data.LanguageId":     {"3003"}, // TODO
			"sourceCode":          {sourceCode},
			"csrf_token":          {a.csrfToken},
		},
	)
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	sid := parseSubmissionId(b)
	a.getSubmissionStatus(contest, sid)

	return true, nil
}

func (a *AtCoderClient) getSubmissionStatus(contest, id string) (*submissionStatus, error) {
	url := fmt.Sprintf("%s/contests/%s/submissions/me/status/json", baseURL, contest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("reload", "true")
	query.Add("sids[]", id)
	req.URL.RawQuery = query.Encode()

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	status := parseSubmissionStatus(b)
	// fmt.Println(string(b))
	// fmt.Println(status)

	return status, nil
}
