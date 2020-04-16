package atcoder

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	a := New()
	if err := a.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD")); err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	a := New()
	a.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	a.Init("https://atcoder.jp/contests/abc126/tasks/abc126_a")
}
