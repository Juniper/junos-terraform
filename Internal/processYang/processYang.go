// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// THIS SOFTWARE IS PROVIDED BY Juniper Networks, Inc. ''AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL Juniper Networks, Inc. BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package processYang

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	s "strings"
	"unicode/utf8"

	"github.com/Juniper/junos-terraform/Internal/cfg"
)

// PrintHeader accepts a message of any length (ideally no more than 80 chars) and pretty prints it in a box
func PrintHeader(message string) {
	header := strings.Repeat("-", utf8.RuneCountInString(message)) + "----" + "\n"
	footer := strings.Repeat("-", utf8.RuneCountInString(message)) + "----" + "\n"
	fmt.Print(header, "- "+message+" -\n", footer)
}

// Node is a helper type for traversing the data tree.
type Node struct {
	XMLName xml.Name
	Key     string `xml:"name,attr"`
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

// Create variable to store groups in the yin file.
var grpNode []Node

// String variable for module name.
var strModuleName string = ""

// String variable for Create Fn.
var strCreate string = ""

// String variable for XML representation.
var startXML string = ""

// String variable for list of yang files.
var yangFileList []string

// Syntactic helper to reduce repetition.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Create Yin files from Yang files and also generate the xpath for the elements
func CreateYinFileAndXpath(jcfg cfg.Config) error {

	filePath := jcfg.YangDir
	fileType := jcfg.FileType

	// Create list of yang files present.
	err := listFiles(filePath)
	if err != nil {
		return err
	}

	// Generate yin file for all yang files.
	generateYinFile(filePath)

	PrintHeader("Creating _xpath files from the Yin files")

	for _, inputYinFile := range yangFileList {

		strCreate = ""
		strModuleName = ""
		grpNode = nil

		// Read data from file.
		dat, err := ioutil.ReadFile(inputYinFile + ".yin")
		if err != nil {
			fmt.Println("DEBUG: Error reading .yin file")
			return err
		}

		// XML decoding.
		buf := bytes.NewBuffer(dat)
		dec := xml.NewDecoder(buf)

		// Create Node based structure, Node is defined above.
		var n Node
		err = dec.Decode(&n)
		if err != nil {
			fmt.Println("DEBUG: Error decododing XML")
			return err
		}

		// Process all the groups in yin file and store them.
		create_group_nodes([]Node{n})

		// Start processing of the data in the file.
		start([]Node{n})

		// create the xpath for the yin/yang file.
		err = createFile(inputYinFile+"_xpath", fileType)
		if err != nil {
			return err
		}
	}
	// No error, return nil.
	return nil
}

// List files and get filenames.
func listFiles(filePath string) error {
	os.Chdir(filePath)

	out, err := exec.Command("ls").Output()
	if err != nil {
		return err
	}

	// Retained for debugging purposes.
	// fmt.Println(string(out))

	output := string(out[:])
	lines := s.Split(output, "\n")
	for _, line := range lines {
		// TODO: Consider using regexp.Compile() for this.
		matched, _ := regexp.Match(`.yang`, []byte(line))
		if matched {
			temp := s.Split(line, ".yang")
			yangFileList = append(yangFileList, temp[0])
		}
	}
	// No error, return nil.
	return nil
}

// Generate yin file for all yang files.
func generateYinFile(filePath string) {

	if !isCommandAvailable("pyang") {
		panic("pyang is not installed")
	}

	// Retained for debugging purposes.
	// fmt.Println(yangFileList)

	PrintHeader("Creating Yin files from Yang file directory: " + filePath)

	for _, file := range yangFileList {
		// The search path is required for included models
		// pyang doesn't provide any output for creating Yin files
		_, err := exec.Command("pyang", "-f", "yin", file+".yang", "-o", file+".yin", "-p", filePath).Output()
		if err != nil {
			panic("pyang error: " + err.Error())
		}

		fmt.Printf("Yin file for %s is generated\n", file)
	}
}

// Create group nodes from []Node.
func create_group_nodes(nodes []Node) {
	for _, n := range nodes {
		if n.XMLName.Local == "grouping" {
			grpNode = append(grpNode, n)
		}
		create_group_nodes(n.Nodes)
	}
}

// The function parses all the elements in the yin/yang file and creates a list of
// the container/list/leaf/leaf-list elements and generate xpath files based on them.
func start(nodes []Node) {
	for _, n := range nodes {

		// All the modules in Juniper yang starts with augment of configuration.
		// First augment is always configuration so we will search it and break
		// the loop.
		if n.XMLName.Local == "augment" {
			var n1 Node
			// The augment for configuration has the very top group and single element
			// so take the 1st node in the list.
			// The group may not be at the very top so iterate all the groups and match
			// the group name. If the group name matches, we need to iterate to all
			// sub nodes.
			for _, n2 := range grpNode {
				if n2.Key == n.Nodes[0].Key {
					n1 = n2
					break
				}
			}
			// n1 is group, the parent container will be the next element.
			for _, n2 := range n1.Nodes {
				if n2.XMLName.Local == "container" {
					strModuleName = n1.Nodes[0].Key
					handleContainer(n2, "", "")
					break
				}
			}
			// Only 1st augment needs to be traversed, so breaking the loop.
			break
		}
		start(n.Nodes)
	}
}

// Parses and appends the path for container and leafs in xpath
func handleContainer(nodes Node, strXpath string, strTab string) {
	strXpath = strXpath + "/" + nodes.Key
	strCreate += strXpath + "\n"
	startXML += strTab + "<" + nodes.Key + "> \n\t" + strTab + "<xpath>" + strXpath + "</xpath>\n"
	strTabTmp := strTab
	strTab += "\t"
	for _, n := range nodes.Nodes {
		// if container append it in the path.
		if n.XMLName.Local == "container" {
			handleContainer(n, strXpath, strTab)
		}
		// if leaf append it in the path.
		if n.XMLName.Local == "leaf" {
			handleContainer(n, strXpath, strTab)
		}
		// if leaf-list append it in the path.
		if n.XMLName.Local == "leaf-list" {
			handleContainer(n, strXpath, strTab)
		}
		// if list append it in the path.
		if n.XMLName.Local == "list" {
			handleContainer(n, strXpath, strTab)
		}
		// if uses , then find corresponding grouping and handle it.
		if n.XMLName.Local == "uses" {
			handleGrouping(n, strXpath, strTab)
		}
		//if choice handle cases
		if n.XMLName.Local == "choice" {
			handleChoices(n, strXpath, strTab)
		}
	}
	startXML += strTabTmp + "</" + nodes.Key + ">\n"
}

// For groups defined in yang file, it needs to be resolved to corresponding container/list
func handleGrouping(nodes Node, strXpath string, strTab string) {

	// We have created a list of all groups as grpNode.
	// If the group name matches, we need to iterate to all sub nodes.
	// First match the grpNode with this node.
	var n1 Node
	for _, n2 := range grpNode {
		if n2.Key == nodes.Key {
			n1 = n2
			break
		}
	}
	// n1 is group, the container and leaf will be the sub elements.
	// It is not mandatory to have one individual container so iterate to
	// each sub node individually.
	for _, n2 := range n1.Nodes {
		if n2.XMLName.Local == "uses" {
			handleGrouping(n2, strXpath, strTab)
		}
		if n2.XMLName.Local == "container" {
			handleContainer(n2, strXpath, strTab)
		}
		if n2.XMLName.Local == "leaf" {
			handleContainer(n2, strXpath, strTab)
		}
		// if leaf-list append it in the path.
		if n2.XMLName.Local == "leaf-list" {
			handleContainer(n2, strXpath, strTab)
		}
		// if list append it in the path.
		if n2.XMLName.Local == "list" {
			handleContainer(n2, strXpath, strTab)
		}
		//if choice handle cases
		if n2.XMLName.Local == "choice" {
			handleChoices(n2, strXpath, strTab)
		}
	}
}

// For choice defined in yang file, it needs to be resolved to corresponding elements
func handleChoices(nodeChoice Node, strXpath string, strTab string) {

	for _, nodeCase := range nodeChoice.Nodes {
		if nodeCase.XMLName.Local == "case" {
			for _, n2 := range nodeCase.Nodes {
				if n2.XMLName.Local == "uses" {
					handleGrouping(n2, strXpath, strTab)
				}
				if n2.XMLName.Local == "container" {
					handleContainer(n2, strXpath, strTab)
				}
				if n2.XMLName.Local == "leaf" {
					handleContainer(n2, strXpath, strTab)
				}
				// if leaf-list append it in the path.
				if n2.XMLName.Local == "leaf-list" {
					handleContainer(n2, strXpath, strTab)
				}
				// if list append it in the path.
				if n2.XMLName.Local == "list" {
					handleContainer(n2, strXpath, strTab)
				}
				//if choice handle cases
				if n2.XMLName.Local == "choice" {
					handleChoices(n2, strXpath, strTab)
				}
			}
		}
	}
}

// Create file with the top container / module name.
func createFile(file string, fileType string) error {

	if fileType == "text" || fileType == "both" {
		var fileName string = s.Join([]string{file, "txt"}, ".")
		fPtr, err := os.Create(fileName)
		if err != nil {
			return err
		}

		// Write to the file.
		_, err = fPtr.WriteString(strCreate)
		if err != nil {
			return err
		}

		fmt.Printf("Creating Xpath file: %s\n", fileName)
	}
	if fileType == "xml" || fileType == "both" {
		var fileNameXML string = s.Join([]string{file, "xml"}, ".")
		fPtr, err := os.Create(fileNameXML)
		if err != nil {
			return err
		}

		// Write to the file.
		_, err = fPtr.WriteString(startXML)
		if err != nil {
			return err
		}

		fmt.Printf("Creating Xpath file: %s\n", fileNameXML)
	}
	// No error, return nil.
	return nil
}

// Print nodes helper. This is to be used for debugging
func print_nodes(nodes []Node, itr int) {
	for _, n := range nodes {
		fmt.Printf("nodes() %s : %s  - %d \n", n.XMLName.Local, n.Key, itr)
		print_nodes(n.Nodes, itr+1)
	}
}

// Is the command available helper.
func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "-v")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
