package atcoder

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

const atcoderrc string = `
username: dqn
password: ''
template_path: ./test/template.cpp
extension: .cpp
test:
  name: ./a.out
  args:
pretest:
  - name: g++
    args:
      - -Wall
      - -std=c++14
      - '{{ file_path }}'
posttest:
  - name: rm
    args:
      - ./a.out
`

func loadConfig() *Config {
	var c Config
	if err := yaml.Unmarshal([]byte(atcoderrc), &c); err != nil {
		panic(err)
	}
	c.Password = os.Getenv("PASSWORD")
	return &c
}

func TestLogin(t *testing.T) {
	a := New(loadConfig())
	if err := a.Login(); err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	a := New(loadConfig())
	a.Login()
	url := "https://atcoder.jp/contests/abc126/tasks/abc126_a"
	if err := a.Init(url); err != nil {
		t.Fatal(err)
	}
}

func TestTest(t *testing.T) {
	a := New(loadConfig())
	a.Login()
	url := "https://atcoder.jp/contests/abc126/tasks/abc126_a"
	a.Init(url)
	copyFile("test/a.cpp", "atcoder/abc126/a.cpp")
	_, err := a.Test("abc126", "a")
	if err != nil {
		t.Fatal(err)
	}
}
