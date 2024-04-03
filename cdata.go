package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Função para limpar as URLs do formato CDATA
func cleanCDATA(url string) string {
	// Remove o prefixo "<![CDATA[" e o sufixo "]]>"
	cleanedURL := strings.TrimPrefix(url, "<![CDATA[")
	cleanedURL = strings.TrimSuffix(cleanedURL, "]]>")
	return cleanedURL
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run CleanURLs.go <arquivo_de_urls>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var cleanedURLs []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urlString := scanner.Text()
		if urlString == "" {
			continue
		}

		// Limpar a URL do formato CDATA
		cleanedURL := cleanCDATA(urlString)
		cleanedURLs = append(cleanedURLs, cleanedURL)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler o arquivo: %s\n", err)
	}

	outputFilePath := filepath.Join(".", "cleaned_urls.txt")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo cleaned_urls.txt: %s\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, u := range cleanedURLs {
		outputFile.WriteString(u + "\n")
	}

	fmt.Println("URLs limpas e salvas em cleaned_urls.txt.")
}

