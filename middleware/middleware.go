package middleware

import (
	echo "github.com/IkezawaYuki/lucky-strike"
	"regexp"
	"strconv"
	"strings"
)

type (
	Skipper    func(echo.Context) bool
	BeforeFunc func(ctx echo.Context)
)

func captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}
	values := groups[0][1:]
	replace := make([]string, 2*len(values))
	for i, v := range values {
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}

func DefaultSkipper(ctx echo.Context) bool {
	return false
}
