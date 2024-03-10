package grading

import "io"

const rule = `(version 1)
(debug deny)
(allow default)
(deny network*)
(deny file-write*)
(deny file-read*)`

const RuleName = "mac-sandbox-rule.sb"

func CreateMacSandboxRules(wr io.Writer) error {
	_, err := wr.Write([]byte(rule))
	return err
}
