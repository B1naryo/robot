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
	// Open the file with URLs
	file, err := os.Open("x.txt")
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Open a file for writing
	outputFile, err := os.Create("urls_xml_inputs.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Create a writer to write to the output file
	writer := bufio.NewWriter(outputFile)

	// Create a slice to store unique URLs
	var uniqueURLs []string

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Iterate over each line of the file
	for scanner.Scan() {
		url := scanner.Text()

		// Check if the URL is already in the list of unique URLs
		if contains(uniqueURLs, url) {
			continue
		}

		// Add the URL to the list of unique URLs
		uniqueURLs = append(uniqueURLs, url)

		// Increment the WaitGroup counter for each goroutine
		wg.Add(1)

		// Call the go function to perform URL checking in a separate goroutine
		go func(url string) {
			defer wg.Done()
			if hasInputs(url) || strings.Contains(url, "?") {
				// If the URL contains inputs or parameters, save it to the output file
				_, err := writer.WriteString(url + "\n")
				if err != nil {
					fmt.Println("Error writing to output file:", err)
					return
				}
			}
		}(url)
	}

	// Check if there was any error during file reading
	if err := scanner.Err(); err != nil {
		fmt.Println("Error during file reading:", err)
		return
	}

	// Wait for all goroutines to finish before proceeding
	wg.Wait()

	// Flush to ensure all data is written to the file
	writer.Flush()

	// Remove duplicate URLs from the output file
	removeDuplicates("urls_xml_inputs.txt")

	fmt.Println("URLs with inputs found have been saved in the urls_with_inputs.txt file.")
}

// Function to check if a URL contains inputs
func hasInputs(url string) bool {
	// Perform a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error accessing URL:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Check if the response body contains the word "input"
	if strings.Contains(strings.ToLower(string(body)), "input") {
		return true
	}

	return false
}

// Function to check if a slice contains a certain element
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Function to remove duplicate URLs from a file
func removeDuplicates(filename string) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a map to store unique URLs
	uniqueURLs := make(map[string]struct{})

	// Read each line of the file
	for scanner.Scan() {
		url := scanner.Text()
		// Add the URL to the map of unique URLs
		uniqueURLs[url] = struct{}{}
	}

	// Check if there was any error during file reading
	if err := scanner.Err(); err != nil {
		fmt.Println("Error during file reading:", err)
		return
	}

	// Open the file for writing
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Write the unique URLs to the output file
	writer := bufio.NewWriter(outputFile)
	for url := range uniqueURLs {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
	writer.Flush()
}
