robot.go

This script is designed to create URLs containing directories found in the robots.txt file. It scans a list of URLs provided by the user (via go run robot.go urls.txt) to check if the URLs have the robots.txt and sitemap.xml files. If it finds these files, it saves these directories concatenated with the URL.

xml.go

This script complements the previous script. If the previous script finds a URL with a file named sitemap{}.xml, it will not be able to save the directories because it differs from sitemap.xml. Therefore, it collects these URLs with the .xml extension, adds them to a text file, and executes: go run xml.go urls_xml.txt. It saves all the directories of these XML files concatenated with the URLs.

input.go

After collecting all the directories in robots and XML, you can use this script to extract all URLs with input data entries. This script is straightforward. It scans a list of URLs to check if they have input data entries, and if found, it saves all the URLs with input data entries.

