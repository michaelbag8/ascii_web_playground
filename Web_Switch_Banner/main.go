package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.TrimLeft(content, "\n")
	rawFile := strings.Split(content, "\n\n")

	bannerMap := make(map[rune][]string)
	charCode := ' '

	for _, raw := range rawFile {
		lines := strings.Split(raw, "\n")
		if len(lines) < 8 {
			continue
		}
		bannerMap[charCode] = lines[:8]
		charCode++
	}
	return bannerMap, nil
}

func RenderLines(input string, banner map[rune][]string) []string {
	var output []string
	for row := 0; row < 8; row++ {
		var result strings.Builder
		for _, ch := range input {
			if lines, ok := banner[ch]; ok && len(lines) > row {
				result.WriteString(lines[row])
			}
		}
		output = append(output, result.String())
	}
	return output
}

func GenerateArt(input string, banner map[rune][]string) string {
	var result strings.Builder

	lines := strings.Split(input, "\\n")
	for _, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}
		content := RenderLines(line, banner)
		for _, row := range content {
			result.WriteString(row)
			result.WriteString("\n")
		}
		
	}
	return result.String()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /ascii-art", handleAsciiArt)
	mux.HandleFunc("GET /ascii-art-switch", handleSwitch)

	fmt.Println("server is running at http://localhost:8080/")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
