package generator

import (
	"ascii-art-web/pkg/fileMgr"
	"fmt"
	"log"
	"strings"
)

var styleNames = []string{"standard", "shadow", "thinkertoy"}
var styles = make(map[string]string)
var style = make(map[rune][]string)

// GenArt() calls the functions to generate ascii art string
func GenArt(txt, styleNm string) (string, error) {
	txtLns, err := checkInput(txt)
	if err != nil {
		return "", err
	}
	var art strings.Builder

	if _, exists := styles[styleNm]; !exists {
		art.WriteString("can't find banner style " + styleNm + "\n")
		art.WriteString("default to standard banner.\n")
		styleNm = "standard"
	}
	getStyle(styleNm)
	for _, txtLn := range txtLns {
		if txtLn == "" {
			art.WriteString("\n")
			continue
		}
		for j := range 8 {
			for _, rn := range txtLn {
				art.WriteString(style[rn][j])
			}
			art.WriteString("\n")
		}
	}
	return art.String(), nil
}

// checkInput() checks txt string and split it by newline
// and cleans up the text of carriage returns, leadind and trailing newlines.
// Returns []txtLns by spliting txt with newlines
func checkInput(txt string) ([]string, error) {
	txtLns := strings.Split(txt, "\n")
	isEmpty := true

	for _, txtLn := range txtLns {
		if txtLn != "" {
			isEmpty = false
		}
		for _, rn := range txtLn {
			if rn < 32 || rn > 127 {
				return nil, fmt.Errorf("character not a printable ASCII char: %s", string(rn))
			}
		}
	}
	if isEmpty {
		return nil, fmt.Errorf("no character to convert")
	}
	return txtLns, nil
}

// getStyle read <styleName>.txt and store the ascii art runes in a map[rune][]string.
func getStyle(styleNm string) {
	rawStyle := strings.Split(styles[styleNm], "\n")
	for i := 1; i < len(rawStyle); i = i + 9 {
		style[rune(32+i/9)] = rawStyle[i : i+8]
	}
}

// GetStyles() load banner style files into a map of strings [styleNm]styles
func GetStyles() {
	for _, styleNm := range styleNames {
		rawStyle := fileMgr.ReadFile("./assets/styles/" + styleNm + ".txt")
		styles[styleNm] = rawStyle
		checkBannerFormat(rawStyle)
	}
}

func checkBannerFormat(rawStyle string) {
	rawStyleSplit := strings.Split(rawStyle, "\n")
	if len(rawStyleSplit) != 856 {
		log.Fatal("Incorrect banner style format")
	}
	for i := 1; i < 856; i = i + 9 {
		if len(rawStyleSplit[i]) != len(rawStyleSplit[i+1]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+2]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+3]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+4]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+5]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+6]) ||
			len(rawStyleSplit[i]) != len(rawStyleSplit[i+7]) {
			log.Fatal("Incorrect banner style format")
		}
	}
}
