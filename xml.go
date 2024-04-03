package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Função para extrair URLs de um sitemap
func extractUrlsFromSitemap(sitemapContent string) []string {
	var urls []string
	re := regexp.MustCompile(`<loc>(.*?)<\/loc>`)
	matches := re.FindAllStringSubmatch(sitemapContent, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}
	return urls
}

// Função para remover URLs duplicadas
func removeDuplicates(urls []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, url := range urls {
		if _, ok := seen[url]; !ok {
			seen[url] = struct{}{}
			result = append(result, url)
		}
	}
	return result
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run Robot.go <arquivo_de_urls_xml>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var allUrls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urlString := scanner.Text()
		if urlString == "" {
			continue
		}

		fmt.Printf("Lendo arquivo XML: %s\n", urlString)

		// Tentar acessar o arquivo XML
		client := http.Client{
			Timeout: time.Second * 10, // Timeout de 10 segundos
		}
		resp, err := client.Get(urlString)
		if err != nil {
			fmt.Printf("Erro ao acessar %s: %s. Pulando para o próximo URL.\n", urlString, err)
			continue
		}
		defer resp.Body.Close()

		// Ler o conteúdo do arquivo XML
		xmlContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Erro ao ler o conteúdo do arquivo XML %s: %s. Pulando para o próximo URL.\n", urlString, err)
			continue
		}

		// Extrair URLs do sitemap XML
		urls := extractUrlsFromSitemap(string(xmlContent))
		allUrls = append(allUrls, urls...)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler o arquivo: %s\n", err)
	}

	// Remover URLs duplicadas
	allUrls = removeDuplicates(allUrls)

	outputFilePath := filepath.Join(".", "ou.txt")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo urls.txt: %s\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, u := range allUrls {
		outputFile.WriteString(u + "\n")
	}

	fmt.Println("URLs extraídas dos sitemaps e salvas em urls.txt.")
}

