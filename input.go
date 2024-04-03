package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	// Abrir o arquivo com as URLs
	file, err := os.Open("x.txt")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar um scanner para ler o arquivo linha por linha
	scanner := bufio.NewScanner(file)

	// Abrir arquivo para escrita
	outputFile, err := os.Create("urls_xml_inputs.txt")
	if err != nil {
		fmt.Println("Erro ao criar arquivo de saída:", err)
		return
	}
	defer outputFile.Close()

	// Criar um writer para escrever no arquivo de saída
	writer := bufio.NewWriter(outputFile)

	// Criar um slice para armazenar as URLs únicas
	var uniqueURLs []string

	// Criar um WaitGroup para esperar que todas as goroutines terminem
	var wg sync.WaitGroup

	// Iterar sobre cada linha do arquivo
	for scanner.Scan() {
		url := scanner.Text()

		// Verificar se a URL já está na lista de URLs únicas
		if contains(uniqueURLs, url) {
			continue
		}

		// Adicionar a URL à lista de URLs únicas
		uniqueURLs = append(uniqueURLs, url)

		// Incrementar o contador do WaitGroup para cada goroutine
		wg.Add(1)

		// Chamar a função go para executar a verificação de URL em uma goroutine separada
		go func(url string) {
			defer wg.Done()
			if hasInputs(url) {
				// Se a URL contiver inputs, salvá-la no arquivo de saída
				_, err := writer.WriteString(url + "\n")
				if err != nil {
					fmt.Println("Erro ao escrever no arquivo de saída:", err)
					return
				}
			}
		}(url)
	}

	// Verificar se houve algum erro durante a leitura do arquivo
	if err := scanner.Err(); err != nil {
		fmt.Println("Erro durante a leitura do arquivo:", err)
		return
	}

	// Aguardar a conclusão de todas as goroutines antes de prosseguir
	wg.Wait()

	// Flush para garantir que todos os dados sejam escritos no arquivo
	writer.Flush()

	// Remover URLs duplicadas do arquivo de saída
	removeDuplicates("urls_xml_inputs.txt")

	fmt.Println("URLs com inputs encontradas foram salvas no arquivo urls_com_inputs.txt.")
}

// Função para verificar se uma URL contém inputs
func hasInputs(url string) bool {
	// Realizar uma requisição GET para a URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao acessar a URL:", err)
		return false
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return false
	}

	// Verificar se o corpo da resposta contém a palavra "input"
	if strings.Contains(strings.ToLower(string(body)), "input") {
		return true
	}

	return false
}

// Função para verificar se um slice contém um determinado elemento
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Função para remover URLs duplicadas de um arquivo
func removeDuplicates(filename string) {
	// Abrir o arquivo para leitura
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar um scanner para ler o arquivo linha por linha
	scanner := bufio.NewScanner(file)

	// Criar um mapa para armazenar as URLs únicas
	uniqueURLs := make(map[string]struct{})

	// Ler cada linha do arquivo
	for scanner.Scan() {
		url := scanner.Text()
		// Adicionar a URL ao mapa de URLs únicas
		uniqueURLs[url] = struct{}{}
	}

	// Verificar se houve algum erro durante a leitura do arquivo
	if err := scanner.Err(); err != nil {
		fmt.Println("Erro durante a leitura do arquivo:", err)
		return
	}

	// Abrir o arquivo para escrita
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erro ao criar arquivo de saída:", err)
		return
	}
	defer outputFile.Close()

	// Escrever as URLs únicas no arquivo de saída
	writer := bufio.NewWriter(outputFile)
	for url := range uniqueURLs {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo de saída:", err)
			return
		}
	}
	writer.Flush()
}

