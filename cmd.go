package atcoder

import (
	"fmt"
	"os/exec"
	"regexp"
)

type replacement struct {
	re    *regexp.Regexp
	value []byte
}

func applyReplacement(s string, replacements []*replacement) string {
	b := []byte(s)
	for _, v := range replacements {
		b = v.re.ReplaceAll(b, v.value)
	}
	return string(b)
}

func executeCommand(command Command, replacements []*replacement) error {
	name := applyReplacement(command.Name, replacements)
	args := make([]string, len(command.Args))
	for i, v := range command.Args {
		args[i] = applyReplacement(v, replacements)
	}
	out, err := exec.Command(name, args...).Output()
	fmt.Print(string(out))
	return err
}
