package atcoder

import (
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

func createTestFile(path string, tests []*test, ch chan error) {
	cases := make([]string, 0, 6)
	for _, v := range tests {
		cases = append(cases, v.input+"\n\n"+v.output+"\n")
	}
	s := strings.Join(cases, "---\n")

	file, err := os.Create(path)
	if err != nil {
		ch <- err
		return
	}
	defer file.Close()

	file.Write([]byte(s))
	ch <- nil
}

func createSourceFile(dir string, ch chan error) {
	ch <- nil
}
