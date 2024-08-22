package util

import "strings"

var escapeLikeSearchChar = []string{
	"%",
	"_",
}

func SanitizeSearchTerm(term string, escapeChar string) string {
	term = strings.Replace(term, escapeChar, "", -1)

	for _, c := range escapeLikeSearchChar {
		term = strings.Replace(term, c, escapeChar+c, -1)
	}

	return term
}
