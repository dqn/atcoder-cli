package atcoder

import (
	"io"
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

func execCommand(command Command, replacements []*replacement) *exec.Cmd {
	name := applyReplacement(command.Name, replacements)
	args := make([]string, len(command.Args))
	for i, v := range command.Args {
		args[i] = applyReplacement(v, replacements)
	}
	return exec.Command(name, args...)
}

func execCommandWithStdin(command Command, replacements []*replacement, input string) (*exec.Cmd, error) {
	cmd := execCommand(command, replacements)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(stdin, input)
	if err != nil {
		return nil, err
	}
	err = stdin.Close()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
