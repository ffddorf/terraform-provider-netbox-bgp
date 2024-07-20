package provider

import (
	"regexp"
	"strings"
)

var (
	// \w = word characters (== [0-9A-Za-z_])
	// \s = whitespace (== [\t\n\f\r ])
	matchSpecial                   = regexp.MustCompile(`[^\w\s-]`)
	matchMultiWhitespacesAndDashes = regexp.MustCompile(`[\s-]+`)
)

func slugify(name string) string {
	var result string
	// Special chars are stripped
	result = matchSpecial.ReplaceAllString(name, "")
	// Blocks of multiple whitespaces and dashes will be replaced by a single dash
	result = matchMultiWhitespacesAndDashes.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	return strings.ToLower(result)
}
