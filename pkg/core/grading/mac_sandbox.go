package grading

import (
	"io"
	"os"
	"path/filepath"
)

const rule = `(version 1)
(debug deny)
(allow default)
(deny network*)
(deny file-write*)`

const RuleName = "mac-sandbox-rule.sb"

func RuleFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Join(wd, RuleName)
}

func CreateMacSandboxRules(wr io.Writer) error {
	_, err := wr.Write([]byte(rule))
	return err
}
