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
