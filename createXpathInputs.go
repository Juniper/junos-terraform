package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// XMLData represents XML data with XPath-like keys
type XMLData map[string][]string

func main() {

	// Directory where XML files are located
	dir := "user_config_files"

	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the files in the directory
	for _, file := range files {
		if file.IsDir() {
			// Skip directories
			continue
		}

		// Check if the file has a .xml extension
		if filepath.Ext(file.Name()) == ".xml" {
			// Open the XML file
			xmlFile, err := os.Open(filepath.Join(dir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer xmlFile.Close()

			// Initialize a map to store XPaths
			xpathMap := make(map[string]bool)

			// Initialize an empty XMLData map
			data := make(XMLData)

			// Convert XML to XPaths and populate the XMLData map
			getXPathsAndDataFromXML(xpathMap, data, xmlFile)

			// Create Xpath Input file
			createXpathInputsFile(xpathMap)

		}
	}
}

// Function traverses XML file, adds XPaths to the map, and stores XML data
func getXPathsAndDataFromXML(xpathMap map[string]bool, data XMLData, xmlFile *os.File) {
	doc := xml.NewDecoder(xmlFile)
	currentPath := []string{}
	currentDataKey := []string{} // Keeps track of the current data key

	// Start decoding the XML file
	for {
		t, err := doc.Token()
		if err != nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			currentPath = append(currentPath, se.Name.Local)
			xpath := strings.Join(currentPath, "/")
			xpathMap[xpath] = true
			// Append the XPath to the data map with an empty value
			data[xpath] = append(data[xpath], "")
			// Append the current element to the currentDataKey
			currentDataKey = append(currentDataKey, se.Name.Local)
		case xml.CharData:
			// Get the text content and build the XPath-like key
			content := strings.TrimSpace(string(se))
			// fmt.Println(content)
			if len(content) > 0 {
				key := strings.Join(currentDataKey, "/")
				content = content + ","
				data[key] = append(data[key], content)
			}
		case xml.EndElement:
			currentPath = currentPath[:len(currentPath)-1]
			// Remove the last element from the currentDataKey
			currentDataKey = currentDataKey[:len(currentDataKey)-1]
		}
	}
}

func createXpathInputsFile(xpathMap map[string]bool) {
	// Part 2: Format and print XPaths directly
	// Create a new file for formatted output

	formattedXpathFile, err := os.Create("xpath_inputs.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer formattedXpathFile.Close()

	// Create an opening <file-list> tag
	fmt.Fprintf(formattedXpathFile, "<file-list>\n")

	// Iterate over the extracted XPaths and format them
	for xpath := range xpathMap {
		if strings.Contains(xpath, "/") {
			// Format the line into XPath
			formattedLine := fmt.Sprintf("    <xpath name=\"/%s\"/>\n", xpath)
			// Write the formatted line to the output file
			fmt.Fprintf(formattedXpathFile, formattedLine)
		}
	}

	// Create a closing </file-list> tag
	fmt.Fprintf(formattedXpathFile, "</file-list>\n")

	fmt.Println()
	fmt.Println("Xpath file written successfully in xpath_inputs.xml")
	fmt.Println()
}
