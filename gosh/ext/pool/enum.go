package pool

import "strings"

const Shell = "bash"
const Python = "python2"
const Python3 = "python3"

func OperatorEnum(v string) string {
	switch strings.ToLower(v) {
	case "py", "py2":
		return Python
	case "py3":
		return Python3
	case "sh":
		return Shell
	default:
		return Shell
	}
}
