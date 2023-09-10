package highlighting

import (
	"regexp"
	"strings"

	"github.com/TwiN/go-color"
)

var pattern = "\\b(" + strings.Join(getKeywords(), "|") + ")\\b"
var keywordColors = map[string]string{
	"var":  color.Purple,
	"if":   color.Purple,
	"else": color.Purple,
	"func": color.Blue,
}

type Highlighter struct {
}

func NewHighlighter() *Highlighter {
	return &Highlighter{}
}

func (h *Highlighter) Highlight(input string) string {

	regex := regexp.MustCompile(pattern)

	coloredInput := regex.ReplaceAllStringFunc(input, func(match string) string {
		return keywordColors[match] + match + color.Reset
	})

	return coloredInput
}

func getKeywords() []string {
	keywords := make([]string, 0, len(keywordColors))
	for keyword := range keywordColors {
		keywords = append(keywords, keyword)
	}
	return keywords
}
