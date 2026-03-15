package random

import (
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/m-oons/lexi/internal/utils"
)

func LoadWords(input string, separator string) ([]string, error) {
	isFile := utils.IsFile(input)
	sep := normalizeSeparator(separator)
	if sep == "" {
		if isFile {
			sep = "\n"
		} else {
			sep = ","
		}
	}

	var lines string
	if isFile {
		content, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}
		lines = string(content)
	} else {
		lines = input
	}

	return parseLines(lines, sep), nil
}

func LoadDefaultWords(separator string) []string {
	return parseLines(string(defaultWords), "\n")
}

func FilterWords(words []string, minLength int, maxLength int, regex string) ([]string, error) {
	if minLength <= 0 {
		minLength = 1
	}
	if maxLength <= 0 {
		maxLength = math.MaxInt32
	}

	var matcher *regexp.Regexp
	if regex != "" {
		compiled, err := regexp.Compile(regex)
		if err != nil {
			return nil, err
		}
		matcher = compiled
	}

	filtered := make([]string, 0)
	for _, word := range words {
		if len(word) < minLength || len(word) > maxLength {
			continue
		}

		if matcher != nil && !matcher.MatchString(word) {
			continue
		}

		filtered = append(filtered, word)
	}

	return filtered, nil
}

func parseLines(lines string, separator string) []string {
	parts := strings.Split(lines, separator)
	words := make([]string, 0, len(parts))

	for _, part := range parts {
		word := strings.TrimSpace(part)
		if word == "" {
			continue
		}

		words = append(words, word)
	}

	return words
}

func normalizeSeparator(separator string) string {
	if separator == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		`\\n`, "\n",
		`\\r`, "\r",
		`\\t`, "\t",
	)

	return replacer.Replace(separator)
}
