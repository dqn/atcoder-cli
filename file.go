package atcoder

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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

func createSourceFile(dir, problem, templetePath string) error {
	ext := templetePath[strings.LastIndex(templetePath, "."):]
	src, err := os.OpenFile(templetePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer src.Close()

	dstPath := fmt.Sprintf("%s/%s%s", dir, problem, ext)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
