package strex

import (
	"fmt"
	"strings"
)

func Join(sep string, args ...interface{}) string {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%v", arg))
	}
	return strings.Join(strArgs, sep)
}

func Concat(args ...interface{}) string {
	return Join("", args...)
}

func RemoveSpace(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func TrimAndTitle(str string) string {
	return RemoveSpace(strings.Title(str))
}
