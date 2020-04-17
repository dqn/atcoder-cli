package atcoder

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func dropUntil(s, until string) string {
	begin := strings.Index(s, until)
	if begin == -1 {
		return ""
	}
	return s[begin+len(until):]
}

func parseContest(s string) (contest, problem string) {
	contest = dropUntil(s, "contests/")[:6]
	problem = strings.ToLower(dropUntil(s, "<title>")[0:1])
	return
}

func parseSample(s string) string {
	s = dropUntil(s, "<pre>")
	return strings.TrimSpace(s[:strings.Index(s, "</pre>")])
}

func parseTests(s string) []*test {
	tests := make([]*test, 0, 6)
	for {
		s = dropUntil(s, "Sample")
		if s == "" {
			return tests
		}
		input := parseSample(s)
		s = dropUntil(s, "Sample")
		output := parseSample(s)
		tests = append(tests, &test{input, output})
	}
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
