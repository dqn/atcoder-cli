package atcoder

import (
	"strconv"
	"strings"
)

type submissionStatus struct {
	Status   string
	Interval uint
	Score    string
}

func dropUntil(s, until string) string {
	begin := strings.Index(s, until)
	if begin == -1 {
		return ""
	}
	return s[begin+len(until):]
}

func extractBetween(s, a, b string) string {
	s = dropUntil(s, a)
	i := strings.Index(s, b)
	if i == -1 {
		return ""
	}
	return s[:i]
}

func parseContest(b []byte) (contest, task string) {
	s := string(b)
	contest = dropUntil(s, "contests/")[:6]
	task = strings.ToLower(extractBetween(s, "<title>", " "))
	return
}

func parseSample(s string) string {
	return strings.TrimSpace(extractBetween(s, "<pre>", "</pre>"))
}

func parseTests(b []byte) []*test {
	s := string(b)
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

func parseSubmissionId(b []byte) string {
	return extractBetween(string(b), `data-id="`, `"`)
}

func parseSubmissionStatus(b []byte) *submissionStatus {
	s := string(b)
	status := extractBetween(s, `"\u003e`, `\u003c/span`)
	score := extractBetween(s, `"Score":"`, `"`)

	interval, err := strconv.Atoi(extractBetween(s, `"Interval":`, ","))
	if err != nil {
		interval = 0
	}

	return &submissionStatus{status, uint(interval), score}
}
