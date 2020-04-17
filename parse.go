package atcoder

import "strings"

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
