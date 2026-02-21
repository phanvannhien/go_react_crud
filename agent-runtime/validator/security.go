package validator

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type SecurityIssue struct {
	File  string
	Line  int
	Rule  string
	Value string
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)api[_-]?key\s*=\s*".+"`),
	regexp.MustCompile(`(?i)secret\s*=\s*".+"`),
	regexp.MustCompile(`(?i)password\s*=\s*".+"`),
	regexp.MustCompile(`(?i)token\s*=\s*".+"`),
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`), // AWS key
}

var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`exec\.Command\(.+\+.+\)`), // string concat in exec
	regexp.MustCompile(`os\.Exec\(`),
	regexp.MustCompile(`syscall\.`),
	regexp.MustCompile(`eval\(`),
	regexp.MustCompile(`rm\s+-rf`),
}

func ScanSecurity(root string) []SecurityIssue {

	var issues []SecurityIssue

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") &&
			!strings.HasSuffix(path, ".js") &&
			!strings.HasSuffix(path, ".ts") {
			return nil
		}

		fileIssues := scanFile(path)
		issues = append(issues, fileIssues...)

		return nil
	})

	return issues
}

func scanFile(path string) []SecurityIssue {

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var issues []SecurityIssue
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()

		for _, pattern := range secretPatterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					File:  path,
					Line:  lineNumber,
					Rule:  "HardcodedSecret",
					Value: strings.TrimSpace(line),
				})
			}
		}

		for _, pattern := range dangerousPatterns {
			if pattern.MatchString(line) {
				issues = append(issues, SecurityIssue{
					File:  path,
					Line:  lineNumber,
					Rule:  "DangerousPattern",
					Value: strings.TrimSpace(line),
				})
			}
		}
	}

	return issues
}

func RunGoSec() []string {

	cmd := exec.Command("gosec", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return []string{"gosec failed: " + string(output)}
	}

	return nil
}
