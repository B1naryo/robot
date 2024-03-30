package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "strings"
    "sync"
)

// Function to extract URLs from the robots.txt file
func extractUrlsFromRobots(robotsUrl string, wg *sync.WaitGroup, urlsChan chan<- []string) {
    defer wg.Done()

    response, err := http.Get(robotsUrl)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error accessing robots.txt at %s: %v\n", robotsUrl, err)
        return
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        fmt.Fprintf(os.Stderr, "Error accessing robots.txt at %s. Status code: %d\n", robotsUrl, response.StatusCode)
        return
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading robots.txt content at %s: %v\n", robotsUrl, err)
        return
    }

    pattern := regexp.MustCompile(`Disallow: (.*)`)
    matches := pattern.FindAllStringSubmatch(string(body), -1)

    var urls []string
    for _, match := range matches {
        // Remove "/robots.txt" from the URL
        url := robotsUrl[:len(robotsUrl)-11] // Remove 10 caracteres da URL
        // Concatenando com o diretório encontrado
        url += match[1]
        urls = append(urls, url)
    }

    urlsChan <- urls
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    var wg sync.WaitGroup
    urlsChan := make(chan []string)

    go func() {
        for scanner.Scan() {
            url := scanner.Text()
            if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
                url = "http://" + url
            }
            robotsUrl := strings.TrimSuffix(url, "/") + "/robots.txt"
            wg.Add(1)
            go extractUrlsFromRobots(robotsUrl, &wg, urlsChan)
        }
        wg.Wait()
        close(urlsChan)
    }()

    var output strings.Builder
    for urls := range urlsChan {
        if len(urls) > 0 {
            for _, url := range urls {
                output.WriteString(url)
                output.WriteString("\n")
            }
            output.WriteString("")
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
        os.Exit(1)
    }

    outputFile := "urlsrobots"
    err := ioutil.WriteFile(outputFile, []byte(output.String()), 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error saving extracted URLs from robots.txt to %s: %v\n", outputFile, err)
        os.Exit(1)
    }

    fmt.Printf("Extracted URLs from robots.txt files have been saved to %s\n", outputFile)
}


