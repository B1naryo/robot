cat <<EOF > README.md
# robot

This Go script facilitates the extraction of URLs disallowed by robots.txt files from a list of websites. By providing a list of website URLs via standard input, the script accesses each site's robots.txt file, extracts the disallowed URLs, and saves them to an output file for later analysis.

## Usage

1. Ensure you have Go installed on your system.
2. Clone the repository or download the script.
3. Navigate to the directory containing the script.
4. Run the script with the command:

   \`\`\`
   go run main.go > domains-robot
   \`\`\`

OR:
cat subdomains-alive | robot > domains-robot

6. Input a list of website URLs, each on a separate line.
7. The script will process each URL, extract the disallowed URLs from the corresponding robots.txt file, and save them to an output file named \`urlsrobots\`.

## Features

- **Concurrent processing:** The script utilizes Goroutines to process multiple URLs concurrently, enhancing performance.
- **Error handling:** It provides error messages for any issues encountered during URL extraction.
- **Output file:** The extracted URLs are saved to an output file for further analysis.

## Example

Suppose you have a list of website URLs stored in a file named \`websites.txt\`:

\`\`\`
http://example.com
http://example.org
http://example.net
\`\`\`

You can run the script and provide this file as input:

\`\`\`bash
cat websites.txt | robot > urls-robots
\`\`\`

The script will extract the disallowed URLs from each site's robots.txt file and save them to \`urlsrobots\`.

## Author

- **Name:** Sandro H... Cerqueira
- **Codename:** bynar1o
EOF



