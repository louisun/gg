package gg

import (
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	breakLineChar = "\n"
	breakLineMark = "---"
)

func wordWrap(f font.Face, text string, width float64, indent bool) []string {
	var (
		lines  []string
		result []string
	)

	rawLines := strings.Split(text, breakLineChar)
	for i, line := range rawLines {
		line = strings.TrimSpace(line)
		if indent && line != "" {
			line = string('\u0009') + string('\u0009') + line
		}
		lines = append(lines, lineSplit(f, line, width)...)
		// break line mark
		if i != len(rawLines) {
			lines = append(lines, breakLineMark)
		}
	}

	for i, line := range lines {
		if i != len(lines) {
			if line != breakLineMark {
				result = append(result, line, "")
			} else {
				result = append(result, "")
			}
		}
	}
	return result
}

func lineSplit(f font.Face, text string, width float64) []string {
	var (
		lines        = make([]string, 0)
		runes        = make([]rune, 0)
		advance      = fixed.Int26_6(0)
		previousRune = rune(-1)
	)

	for _, currentRune := range text {
		if previousRune >= 0 {
			advance += f.Kern(previousRune, currentRune)
		}

		a, ok := f.GlyphAdvance(currentRune)
		if !ok {
			continue
		}

		advance += a

		// new line
		if float64(advance>>6) > width {
			lines = append(lines, string(runes))
			runes = make([]rune, 0)
			runes = append(runes, currentRune)

			a, ok := f.GlyphAdvance(currentRune)
			if !ok {
				a = fixed.Int26_6(0)
			}

			advance = a
			previousRune = rune(-1)
		} else {
			runes = append(runes, currentRune)
			previousRune = currentRune
		}
	}

	if len(runes) != 0 {
		lines = append(lines, string(runes))
	}

	return lines
}
