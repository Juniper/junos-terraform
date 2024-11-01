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

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        if filepath.Ext(file.Name()) == ".xml" {
            xmlFile, err := os.Open(filepath.Join(dir, file.Name()))
            if err != nil {
                log.Fatal(err)
            }
            defer xmlFile.Close()

            xpathMap := make(map[string]bool)
            data := make(XMLData)

            // Extract XPaths and populate XMLData map
            getXPathsAndDataFromXML(xpathMap, data, xmlFile)

            // Create filtered unique Xpath inputs file
            createUniqueXpathInputsFile(xpathMap)
        }
    }
}

// createUniqueXpathInputsFile filters unique top-level elements and writes them to an output file
func createUniqueXpathInputsFile(xpathMap map[string]bool) {
    uniqueTopLevel := make(map[string]bool)
    for xpath := range xpathMap {
        elements := strings.Split(xpath, "/")
        if len(elements) > 1 {
            uniqueTopLevel[elements[1]] = true
        }
    }

    // Create the output file with <file-list> wrapper
    outputFile, err := os.Create("filtered_xpath_inputs.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer outputFile.Close()

    // Write opening tag for file-list
    _, err = outputFile.WriteString("<file-list>\n")
    if err != nil {
        log.Fatal(err)
    }

    // Write each unique top-level element as an <xpath> entry
    for topElement := range uniqueTopLevel {
        _, err := outputFile.WriteString(fmt.Sprintf("   <xpath name=\"/%s\"/>\n", topElement))
        if err != nil {
            log.Fatal(err)
        }
    }

    // Write closing tag for file-list
    _, err = outputFile.WriteString("</file-list>\n")
    if err != nil {
        log.Fatal(err)
    }
}

// Function to extract XPaths from XML file
func getXPathsAndDataFromXML(xpathMap map[string]bool, data XMLData, xmlFile *os.File) {
    doc := xml.NewDecoder(xmlFile)
    currentPath := []string{}

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
        case xml.EndElement:
            if len(currentPath) > 0 {
                currentPath = currentPath[:len(currentPath)-1]
            }
        }
    }
}
