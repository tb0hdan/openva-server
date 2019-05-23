package stringutil

import (
	"regexp"
	"strings"
)

func SplitWords(path string) (newString string) {
	return doWordMagic(path, " ")
}

func MergeWords(path string) (newString string) {
	return doWordMagic(path, "")
}

func doWordMagic(path, separator string) (newString string) {
	re := regexp.MustCompile(`[/|_|-|-|(|)|\.]`)
	for _, str := range strings.Split(re.ReplaceAllString(path, separator), separator) {
		if strings.TrimSpace(str) == "" {
			continue
		}
		newString += separator + str
	}
	newString = strings.TrimSpace(newString)
	return
}
