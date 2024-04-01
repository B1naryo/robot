package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"
)


func extractUrlsFromSitemap(sitemapContent string, baseUrl string) []string {
	var urls []string
	re := regexp.MustCompile(`<loc>(.*?)<\/loc>`)
	matches := re.FindAllStringSubmatch(sitemapContent, -1)
	for _, match := range matches {
		urls = append(urls, match[1])
	}
	fmt.Printf("URLs extraídas do sitemap de %s.\n", baseUrl)
	return urls
}

func extractUrlsFromRobots(robotsContent string, baseUrl string) []string {
	var urls []string
	re := regexp.MustCompile(`Disallow: (.*)`)
	matches := re.FindAllStringSubmatch(robotsContent, -1)
	for _, match := range matches {
		urls = append(urls, baseUrl+match[1])
	}
	fmt.Printf("URLs extraídas do robots.txt de %s.\n", baseUrl)
	return urls
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run Robot.go <arquivo_de_urls>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var extractedUrls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urlString := scanner.Text()
		if urlString == "" {
			continue
		}

		baseUrl := "http://" + urlString
		fmt.Printf("URL base: %s\n", baseUrl)

		// Verificar se a URL é válida
		_, err := url.Parse(baseUrl)
		if err != nil {
			fmt.Printf("URL inválida: %s\n", baseUrl)
			continue
		}

		// Tentar acessar o site
		client := http.Client{
			Timeout: time.Second * 10, // Timeout de 10 segundos
		}
		resp, err := client.Head(baseUrl)
		if err != nil {
			fmt.Printf("Erro ao acessar %s: %s. Pulando para o próximo URL.\n", baseUrl, err)
			continue
		}
		resp.Body.Close()

		// Extrair URLs do sitemap
		sitemapUrl := baseUrl + "/sitemap.xml"
		fmt.Printf("URL do sitemap: %s\n", sitemapUrl)
		resp, err = client.Get(sitemapUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			sitemapContent, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				extractedUrls = append(extractedUrls, extractUrlsFromSitemap(string(sitemapContent), baseUrl)...)
			} else {
				fmt.Printf("Erro ao ler o conteúdo do sitemap de %s: %s\n", baseUrl, err)
			}
		} else {
			fmt.Printf("Erro ao acessar o sitemap de %s: %s\n", baseUrl, err)
		}

		// Extrair URLs do robots.txt
		robotsUrl := baseUrl + "/robots.txt"
		fmt.Printf("URL do robots.txt: %s\n", robotsUrl)
		resp, err = client.Get(robotsUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			robotsContent, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				extractedUrls = append(extractedUrls, extractUrlsFromRobots(string(robotsContent), baseUrl)...)
			} else {
				fmt.Printf("Erro ao ler o conteúdo do robots.txt de %s: %s\n", baseUrl, err)
			}
		} else {
			fmt.Printf("Erro ao acessar o robots.txt de %s: %s\n", baseUrl, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler o arquivo: %s\n", err)
	}

	outputFilePath := filepath.Join(".", "urls.txt")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo urls.txt: %s\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	for _, u := range extractedUrls {
		outputFile.WriteString(u + "\n")
	}

	fmt.Println("Diretórios extraídos de robots.txt e sitemap.xml e salvos em urls.txt.")
}



