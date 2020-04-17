package atcoder

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

const atcoderrc string = `username: dqn
password: ''
templete_path: './test/templete.cpp'
`

func loadConfig() *Config {
	var c Config
	yaml.Unmarshal([]byte(atcoderrc), &c)
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
