package atcoder

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func createTestFile(path string, tests []*test) error {
	cases := make([]string, 0, 6)
	for _, v := range tests {
		cases = append(cases, v.input+"\n\n"+v.output+"\n")
	}
	s := strings.Join(cases, "---\n")

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(s))
	if err != nil {
		return err
	}

	return nil
}

func createSourceFile(dir, problem, templatePath string) error {
	ext := templatePath[strings.LastIndex(templatePath, "."):]
	dstPath := fmt.Sprintf("%s/%s%s", dir, problem, ext)
	return copyFile(templatePath, dstPath)
}

func readTests(path string) ([]*test, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cases := strings.Split(string(b), "\n---\n")
	tests := make([]*test, len(cases))

	for i, v := range cases {
		samples := strings.Split(v, "\n\n")
		tests[i] = &test{
			input:  strings.TrimSpace(samples[0]),
			output: strings.TrimSpace(samples[1]),
		}
	}
	return tests, nil
}
