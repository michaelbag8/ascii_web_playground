package ascii

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadBanner(name string) (map[rune][]string, error){
	data, err := os.ReadFile(name)
	if err!=nil{
		return nil, err
	}

	if len(data) == 0{
		return nil, errors.New("File is empty")
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	rawFile := strings.Split(content, "\n\n")

	bannerMap := make(map[rune][]string)
	charCode := ' '


	for _, raw := range rawFile{
		lines := strings.Split(raw, "\n")
		if len(lines)< 8{
			return nil, errors.New("Lines are more than 8")
		}
		bannerMap[charCode]= lines[:8]
		charCode++
	}
	return  bannerMap, nil
}



func RenderLines(input string, banner map[rune][]string) []string{
	var output []string
	for row:=0; row < 8; row++{
		var result strings.Builder
		for _, ch := range input{
			result.WriteString(banner[ch][row])
		}
		output = append(output, result.String())
	}
	return output
}

func GenerateArt(input string, banner map[rune][]string) string{
	var result strings.Builder
	userInput := strings.Split(input, "\\n")
	for i, lines := range userInput{
		if lines == ""{
			if i < len(userInput)-1{
				result.WriteString("\n")
			}
			continue
		}

		for _, ch := range lines{
			if ch < 32 || ch > 126{
				fmt.Fprintln(os.Stderr, "Non ascii character", ch)
				os.Exit(1)
			}
		}
		row := RenderLines(input, banner)
		for _, rows := range row{
			result.WriteString(rows)
			result.WriteString("\n")
		}
		if i < len(userInput)-1{
				result.WriteString("\n")
			}
	}
	return result.String()
}