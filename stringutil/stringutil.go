package stringutil

import (
	"regexp"
	"strings"
)

func SplitWords(path string) (newString string) {
	re := regexp.MustCompile(`[/|_|-|-|(|)|\.]`)
	for _, str := range strings.Split(re.ReplaceAllString(path, " "), " ") {
		if strings.TrimSpace(str) == "" {
			continue
		}
		newString += " " + str
	}
	return
}
