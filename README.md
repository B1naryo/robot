Robot

The robot automates the extraction of URLs from sitemap.xml and robots.txt files associated with provided domain names. This tool streamlines the process of gathering URLs for analysis by penetration testers, web developers, and security analysts.

Features:

Extracts URLs from sitemap.xml and robots.txt files.
Handles HTTP requests efficiently and manages timeouts.
Saves extracted URLs to a specified output file for analysis.
Usage:

Clone the repository.
Run the script with the following command: 
cat domains | robot urls
Provide the input file containing domain names.
The extracted URLs will be saved to the urls.txt file in the current directory.
Requirements:

Go programming language
Internet connection
Note: Ensure that the provided domain names are accessible and have valid sitemap.xml and robots.txt files.


