// Copyright (c) 2017-2022, Juniper Networks Inc. All rights reserved.
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

package processProviders

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	s "strings"
	"sync/atomic"
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
	Val     string `xml:"value,attr"`
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

// ElementName is a helper type for traversing the data tree.
type ElementName struct {
	name  string
	count int
}

// Create variable to append groups if the flag is set.
var isGrpFlag bool = false

// Create variable to check if xpath found or not.
var isXpathFound bool = false

// Create variable to store groups in the yin file.
var grpNode []Node

// Enums contains the enumerated choices from the YANG files under a given 'uses' grouping
// The first map key is the group name and the second map key is the enum name
// It's possible to store data under the enum according to the YANG spec, so a string var is required
var enums map[string]map[string]string

var elementNameList []ElementName

// String variable for module name.
var strModuleName string

// String variable for input Xpath.
var inputXpath string = ""

// String variable for yin file.
var inputYinFile string = ""

// String variable for xpath for structure.
var structXpath string

// String variable for import variables.
var strImport string

// String variable to hold path to TFtemplate folder
var targetDir string

// String variable for structure.
var strStruct string
var strStructEnd string

// string variable for schema.
var strSchema string
var strClientInit string
var strSendTrans string
var strSendTransId string
var strSetIdValue string

// String variable for create and update function.
var strGetFunc string
var strSetFunc string
var strVarAssign string

// String variable for Create Fn.
var strCreate string

// String variable for Read Fn.
var strRead string

// String variable for Update Fn.
var strUpdate string

// String variable for Delete Fn.
var strDelete string

// String variable for list of yang files.
var yang_file_list []string

// String variable for provider.go file.
var providerFileData string

// String variable for tabs in schema.
var strSchemaTab string

// Variable to hold current XPath being processed
var currentXPath string

// Issue Counter
var issueCounter int

// XPath Counter
var xpathCounter int

var issue_xpaths []string

// Syntactic helper to reduce repetition.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CreateProviders consumes a YANG file, Xpath file and module name to create a provider
func CreateProviders(jcfg cfg.Config) error {

	PrintHeader("Autogenerating Terraform Provider code from _xpath files")

	yangFilePath := jcfg.YangDir
	xpathFilePath := jcfg.XpathPath
	moduleFilePath := jcfg.ProviderDir

	// Create list of yang files present
	listFiles(yangFilePath, jcfg)

	// Read data from xpath file
	datIn, err := ioutil.ReadFile(xpathFilePath)
	if err != nil {
		return err
	}

	// XML decoding
	bufIn := bytes.NewBuffer(datIn)
	decIn := xml.NewDecoder(bufIn)

	// Create Node based structure, Node is defined above
	var inNode Node
	err = decIn.Decode(&inNode)
	if err != nil {
		return err
	}

	// parse the xpaths provided to generate terraform based modules
	counter := 0
	numofJobs := len(inNode.Nodes)
	// fmt.Printf(numofJobs)

	// Create terraform main and test file
	createTerraformMain(jcfg)
	createTerraformTest(jcfg)

	for _, n5 := range inNode.Nodes {

		// fmt.Println("DEBUG: Working with XPath Expression(n5.Key): -> ", n5.Key)
		currentXPath = n5.Key
		// ADD Check here to see if n5.KEY is not empty
		// test by adding empty xpath to xpath_example
		if currentXPath == "" {
			fmt.Println("EMPTY XPATH found, remove this from the xpath file. Issue Marked")
			issue_xpaths = append(issue_xpaths, currentXPath)
			issueCounter += 1
		}
		if currentXPath != "" {
			if n5.XMLName.Local == "xpath" {
				inputXpath = n5.Key
				strParts := s.Split(inputXpath, "/")
				yangCheck := "conf-" + strParts[1] + "@"
				for _, file := range yang_file_list {
					if s.Contains(file, yangCheck) {
						inputYinFile = file + ".yin"
						break
					}
				}
				isGrpFlag = true
				for _, n2 := range n5.Nodes {
					if n2.XMLName.Local == "group-flag" {
						if n2.Key == "false" {
							isGrpFlag = false
						}
					}
				}
				// Set global variables to default values
				initialize_global_variables()
				// Parse data from Yin file
				dat, err := ioutil.ReadFile(inputYinFile)
				if err != nil {
					return err
				}
				// XML decoding
				buf := bytes.NewBuffer(dat)
				dec := xml.NewDecoder(buf)
				// Create Node based structure, Node is defined above
				var n Node
				err = dec.Decode(&n)
				if err != nil {
					return err
				}
				// Process all the groups in yin file and store them
				create_group_nodes([]Node{n})
				// Process all of the 'uses enums' as choice-ident and choice-value
				for _, v := range grpNode {
					// fmt.Println("DEBUG: Looking at nodes in V1: ", v.Key)
					// fmt.Println("In file: ", inputYinFile)
					// if v.Key == "control_route_filter_type" {
					// fmt.Println("looking for choice-ident")
					for _, v2 := range v.Nodes {
						// fmt.Println("DEBUG: Looking at nodes in v2: ", v2.Key)
						if v2.Key == "choice-ident" {
							// fmt.Println("Found choice-ident...")
							for _, v3 := range v2.Nodes {
								// fmt.Println("looking for enumeration...")
								if v3.Key == "enumeration" {
									// fmt.Println("Found enumeration for choice-ident...")
									// fmt.Println("Length of enums: ", len(v3.Nodes))
									for _, v4 := range v3.Nodes {
										// fmt.Println(v4.Key)
										if enums[v.Key] == nil {
											enums[v.Key] = make(map[string]string)
										}
										enums[v.Key][v4.Key] = v4.Key
									}
								}
							}
						}
					}
				}
				// End of 'uses' enum processing for lists.

				isXpathFound = false

				// Start processing of the file data
				// Notes : "-" and "." is not allowed in go as variable name. need to replace it with "_"
				start([]Node{n})

				if isXpathFound {
					// After all the data processing is done, create the file.
					err = createFile(moduleFilePath, jcfg)
					if err != nil {
						fmt.Println("Issue creating file. Check presence of directory and permissions")
						os.Exit(0)
					}
				} else {
					issueCounter += 1
					fmt.Printf("[ISSUE]: Terraform API resource for %s could NOT be created\nXPATH doesn't have valid match in Yang Files therefore has been removed from xapth_inputs file.\n", currentXPath)
					issue_xpaths = append(issue_xpaths, currentXPath)
				}
			}
		}
		counter++
		printProgressBar(counter, numofJobs, "Progress", "Complete", 25, "=")
	}

	providerFileData += `			"junos-` + jcfg.ProviderName + `_commit": junosCommit(),
	        "junos-` + jcfg.ProviderName + `_destroycommit": junosDestroyCommit(),
			},
		    ConfigureContextFunc: providerConfigure,
	    } 
    }`

	// Create provider.go file
	var fileName string = "provider.go"
	fileName = moduleFilePath + "/" + fileName
	fPtr, err := os.Create(fileName)
	defer fPtr.Close()
	check(err)

	fixXPath_Inputs()

	// Write to the file
	_, err = fPtr.WriteString(providerFileData)

	// List summary data
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Number of Xpaths processed: ", xpathCounter)
	fmt.Println("Number of potential issues: ", issueCounter)
	fmt.Println("	Search for [ISSUE]")

	// Change path to the terraform_provider dir in this project
	// we don't know the exact location, so find it through some path building
	path, _ := os.Executable()
	pathBits := strings.Split(path, "/")
	pathBitsLen := len(pathBits)
	pathBits = pathBits[:pathBitsLen-3]
	pathBits = append(pathBits, "terraform_providers")

	tpPath := ""
	pathBitsLen = len(pathBits)
	for k, v := range pathBits {
		tpPath += v
		if k != pathBitsLen-1 {
			tpPath += "/"
		}
	}

	// Copy the files from ../terraform_providers to the `providerDir` from the config file
	files, err := ioutil.ReadDir(tpPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	PrintHeader("Copying the rest of the required Go files")
	for _, f := range files {
		if strings.Contains(f.Name(), ".go") {
			input, err := ioutil.ReadFile(tpPath + "/" + f.Name())
			if err != nil {
				fmt.Println(err)
			}

			err = ioutil.WriteFile(jcfg.ProviderDir+"/"+f.Name(), input, 0644)
			if err != nil {
				fmt.Println("Error creating", jcfg.ProviderDir+"/"+f.Name())
			}

			fmt.Printf("Copied file: %+v to %+v\n", f.Name(), jcfg.ProviderDir)
		}
	}

	// fmt.Println("------------------------------------------------------------")
	// fmt.Println(fmt.Sprintf("- Copied a total of %d .go files from %s to %s -", fileCopyCount, filepath.Base(tpPath), filepath.Base(jcfg.ProviderDir)))
	// fmt.Println("------------------------------------------------------------")

	PrintHeader("Creating Go Mod")
	err = ioutil.WriteFile(jcfg.ProviderDir+"/go.mod", []byte(fmt.Sprintf(gomodcontent, jcfg.ProviderName)), 0644)
	if err != nil {
		fmt.Println("Error creating", jcfg.ProviderDir+"/go.mod")
	}

	PrintHeader("Creating provider config")
	// fmt.Println(providerConfigContent)
	err = ioutil.WriteFile(jcfg.ProviderDir+"/config.go", []byte(fmt.Sprintf(providerConfigContent, jcfg.ProviderName)), 0644)
	if err != nil {
		fmt.Println("Error creating", jcfg.ProviderDir+"/config.go")
	}
	// No errors, so return nil.
	return nil
}

// Initialize global variables to default values
func initialize_global_variables() {

	grpNode = []Node{}

	enums = make(map[string]map[string]string)

	elementNameList = []ElementName{}

	strModuleName = ""

	structXpath = ""

	strImport = `
package main

import (
    "context"
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

)

`

	strStruct = `
// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex - interface is keyword in golang
`

	strStructEnd = ""

	strSchema = `

        Schema: map[string]*schema.Schema{
            "resource_name": &schema.Schema{
                Type:    schema.TypeString,
                Required: true,
            },`

	strClientInit = `(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var err error
	client := m.(*ProviderConfig)

    id := d.Get("resource_name").(string)
    `

	strSendTrans = `
    err = client.SendTransaction("", config, false)
    check(ctx, err)
    `

	strSendTransId = `
    err = client.SendTransaction(id, config, false)
    check(ctx, err)
    `

	strSetIdValue = `
    d.SetId(fmt.Sprintf("%s_%s", client.Host, id))
    `

	strGetFunc = " "
	strSetFunc = " "
	strVarAssign = ""

	strCreate = `

// v_ is appended before every variable so it doesn't give any conflict
// with any keyword in golang. ex- interface is keyword in golang
`

	strRead = ""

	strUpdate = ""

	strDelete = `
	_, err = client.DeleteConfig(id,false)
    check(ctx, err)

    d.SetId("")
    `

	strSchemaTab = "\t"
}

// The function parses all the elements in the yin/yang file and creates a list of
// the container/list/leaf/leaf-list elements and generate module files based on them.
func start(nodes []Node) {
	for _, n := range nodes {

		// All the modules in Juniper yang starts with augment of configuration
		// first augment is always configuration so we will search it and break
		// the loop
		if n.XMLName.Local == "augment" {
			var n1 Node
			// The augment for configuration has the very top group and has a single element
			// so take the 1st node in the list.
			// grp node is extracted before this start function.
			// the group may not be at the very top so iterate all the groups and match
			// the group name. If the group name matches, we need to iterate to all
			// sub nodes.
			for _, n2 := range grpNode {
				if n2.Key == n.Nodes[0].Key {
					n1 = n2
					break
				}
			}
			// n1 is group, the parent container will be the next element.
			// the parent container is expected to be a single element.
			for _, n2 := range n1.Nodes {
				if n2.XMLName.Local == "container" || (n2.XMLName.Local == "list" && n2.Key == "logical-systems") {
					// Append information at start of the code blocks in the file
					// It needs the parent container name and xpath so it is being done here than
					// in the main function itself.
					strXpath := inputXpath
					strXpath = s.Replace(strXpath, "/", "", 1)
					strXpath = s.ReplaceAll(strXpath, "/", " ")
					strXpath = s.Title(strXpath)
					strXpath = s.ReplaceAll(strXpath, " ", "")
					initializeFunctionString(strXpath)
					// initializeFunctionString(n2.Key)

					// Check if the xpath matches. then start the schema/struct creation.
					matchXpath(n2)
					break
				}
			}
			// Only first augment needs to be traversed, so breaking the loop.
			break
		}
		start(n.Nodes)
	}
}

// check if the xpath provided is a valid xpath in yang file
func matchXpath(nodes Node) {
	// Handle variable for xpath hierarchy for structure
	// remove the 1st / from the xpath as it is not required
	// for the rest xpath replace / with > as it will be used in struct.
	strXpath := inputXpath
	strXpath = s.Replace(strXpath, "/", "", 1)
	// structXpath = s.ReplaceAll(structXpath, "/", ">")

	strParts := s.Split(strXpath, "/")
	structXpath = strParts[0]

	nodeCheck := nodes
	schemaTab := "\t"
	var strStructHierarchy string = "config"

	if isGrpFlag {
		schemaTab += "\t"
		strStructHierarchy += ".Groups"
	}

	// If any element in xpath is list, we need to add key of it in structure to
	// be able to pass it to the device as valid configuration.
	// If any element before the last element is list in xpath, we will process and
	// add key as a parameter in the generated module.
	// We accept leaf as last element in xpath. In such scenario, leaf can't be part of
	// structure. it need to be an element. so if last element is leaf, call handleLeaf
	// after processing the previous container. list would have been processed.
	var structXpath_last_elem string
	// TODO: TIDY THIS UP node_last_elemt Node   // no need to store last element, it will be nodeCheck
	var node_last_elemt_2 Node

	// If the topmost container is chosen, don't process xpath. --> IGNORE
	if len(strParts) > 1 {
		var matchFound bool
		for itr := 1; itr < len(strParts); itr++ {
			matchFound = false
			for _, n := range nodeCheck.Nodes {
				// If the next element is a container , list , leaf-list or leaf
				// it can be a possible chance for xpath match.
				// If the element is uses, then find corresponding group and handle it.
				if check_node_tag(n.XMLName.Local) {
					if n.Key == strParts[itr] {
						// If element matches, break the inner loop
						// run next loop on sublements of this node.
						nodeCheck = n
						matchFound = true
					}
				} else if n.XMLName.Local == "uses" {
					nodeGrp, flag := matchGroupingXpath(n.Key, strParts[itr])
					if flag {
						// If element matches, break the inner loop
						// run next loop on sublements of this node.
						nodeCheck = nodeGrp
						matchFound = true
					}
				} else if n.XMLName.Local == "choice" {
					nodeGrp, flag := matchChoiceXpath(n, strParts[itr])
					if flag {
						// If element matches, break the inner loop
						// run next loop on sublements of this node.
						nodeCheck = nodeGrp
						matchFound = true
					}
				}

				if matchFound {
					if itr == (len(strParts) - 1) {
						isXpathFound = true
						structXpath_last_elem = strParts[itr]
					} else {
						if structXpath == "" {
							structXpath = strParts[itr]
						} else {
							structXpath += ">" + strParts[itr]
						}
						// Handle for list and create parameter for key.
						if nodeCheck.XMLName.Local == "list" {
							if itr == (len(strParts) - 2) {
								// this is done to avoid duplication of key element when passed in xpath as end-element
								structXpath_last_elem = strParts[itr+1]
							}

							uses := ""

							for _, tmp := range nodeCheck.Nodes {
								if tmp.XMLName.Local == "uses" {
									if tmp.Key != "apply-advanced" {
										// fmt.Println("DEBUG: Found: uses-> ", tmp.Key)
										uses = tmp.Key
									}
								}
							}

							strStructHierarchy, structXpath, schemaTab = setListXpathMatch(uses, nodeCheck, schemaTab, structXpath, strStructHierarchy, structXpath_last_elem)
						}

						// For 2nd last element store the node.
						if itr == (len(strParts) - 2) {
							node_last_elemt_2 = nodeCheck
						}
					}
					// If match found, break the for loop.
					break
					// End of if check for match found.
				}

				// End of looping of nodes.
			}
			if !matchFound {
				// issueCounter += 1
				// fmt.Printf("[ISSUE]: Xpath not found in file, check it. : %s \n", strXpath)
				return
			}

			// End of for loop for xpath.
		}

		// If last element is leaf, further processing will be handled here and return from here.
		if nodeCheck.XMLName.Local == "leaf" || nodeCheck.XMLName.Local == "leaf-list" {

			if node_last_elemt_2.XMLName.Local == "list" {
				handleLeaf(nodeCheck, strStructHierarchy, schemaTab)
			} else {
				// It is a container. add structure for it.
				val_ := s.ReplaceAll(node_last_elemt_2.Key, "-", "__")
				val_ = s.ReplaceAll(val_, ".", "__")
				id := check_element_name(node_last_elemt_2.Key)
				if id != 0 {
					val_ += "__" + strconv.Itoa(int(id)) //string(id)
				}
				// fmt.Println("DEBUG: Check val_name 0: ", val_)
				strStruct += "\n" + schemaTab + "V_" + val_ + "  struct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + node_last_elemt_2.Key + "\"`"
				strStructHierarchy += ".V_" + val_

				handleLeaf(nodeCheck, strStructHierarchy, schemaTab+"\t")

				strStruct += "\n" + schemaTab + "} `xml:\"" + structXpath + "\"`"
			}
			strStruct += strStructEnd

			return
		}

		// Append structXpath for last element.
		if structXpath == "" {
			structXpath = structXpath_last_elem
		} else {
			structXpath += ">" + structXpath_last_elem
		}

		// End of if check.
	}

	handleParentNodeXpath(nodeCheck, strStructHierarchy, schemaTab)
	strStruct += strStructEnd
}

// Match xpath for a group represented as 'uses', resolve for the corresponding structure
func matchGroupingXpath(nodeName string, xpathElem string) (Node, bool) {
	// We have created a list of all groups as grpNode.
	// If the group name matches, we need to iterate to all sub nodes.
	// First match the grpNode with this node.
	var n1 Node
	for _, n2 := range grpNode {
		if n2.Key == nodeName {
			n1 = n2
			break
		}
	}

	nodeCheck := n1
	var flag bool = false
	for _, n := range nodeCheck.Nodes {
		// If the next element is a container , list , leaf-list or leaf
		// it can be a possible chance for xpath match.
		// If the element is uses, then find corresponding group and handle it.
		if check_node_tag(n.XMLName.Local) {
			if n.Key == xpathElem {
				nodeCheck = n
				flag = true
				break
			}
		} else if n.XMLName.Local == "uses" {
			nodeGrp, flag := matchGroupingXpath(n.Key, xpathElem)
			if flag {
				// If element matches, break the inner loop
				// run next loop on sublements of this node.
				return nodeGrp, flag
			}
		} else if n.XMLName.Local == "choice" {
			nodeGrp, flag := matchChoiceXpath(n, xpathElem)
			if flag {
				// If element matches, break the inner loop
				// run next loop on sublements of this node.
				return nodeGrp, flag
			}
		}
	}
	return nodeCheck, flag
}

// Match xpath for a choice, resolve for the corresponding structure
func matchChoiceXpath(nodeChoice Node, xpathElem string) (Node, bool) {
	// Iterate through the list of cases in the option choice.
	// If it matches, handle it.
	nodeCheck := nodeChoice
	var flag bool = false
	for _, nodeCase := range nodeChoice.Nodes {
		if nodeCase.XMLName.Local == "case" {
			for _, n := range nodeCase.Nodes {
				// If the next element is a container , list , leaf-list or leaf
				// it can be a possible chance for xpath match.
				// if check_node_tag(n.XMLName.Local) {
				if n.Key == xpathElem {
					nodeCheck = n
					flag = true
					break
				} else if n.XMLName.Local == "uses" {
					nodeGrp, flag := matchGroupingXpath(n.Key, xpathElem)
					if flag {
						// If element matches, break the inner loop
						// run next loop on sublements of this node.
						return nodeGrp, flag
					}
				} else if n.XMLName.Local == "choice" {
					nodeGrp, flag := matchChoiceXpath(n, xpathElem)
					if flag {
						// If element matches, break the inner loop
						// run next loop on sublements of this node.
						return nodeGrp, flag
					}
				}
			}
		}
		if flag {
			break
		}
	}
	return nodeCheck, flag
}

// Handle xpath matching for list. Need to add key also in case of list during xpath matching
func setListXpathMatch(uses string, nodeCheck Node, schemaTab string, structXpath string, strStructHierarchy string, structXpath_last_elem string) (string, string, string) {
	var keyValue string

	for _, n1 := range nodeCheck.Nodes {

		if n1.XMLName.Local == "key" {
			keyValue = n1.Val
			break
		}

	}

	// Assign values for list and its key values.
	val_ := s.ReplaceAll(nodeCheck.Key, "-", "__")
	val_ = s.ReplaceAll(val_, ".", "__")
	// Duplicate name check for list.
	id := check_element_name(nodeCheck.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}
	// fmt.Println("DEBUG: Check val_name 1: ", val_)
	strStruct += "\n" + schemaTab + "V_" + val_ + "  struct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodeCheck.Key + "\"`"
	strStructEnd = "\n" + schemaTab + "} `xml:\"" + structXpath + "\"`" + strStructEnd
	strStructHierarchy += ".V_" + val_
	schemaTab += "\t"

	elements := s.Split(keyValue, " ")

	for _, keyVar := range elements {

		// required to remove choice-ident and choice-value from 'uses' within a YANG list construct
		if keyVar != "choice-ident" && keyVar != "choice-value" {
			// this is done to avoid duplication of key element when passed in xpath as end-element
			if structXpath_last_elem != keyVar {
				val_ = s.ReplaceAll(keyVar, "-", "__")
				val_ = s.ReplaceAll(val_, ".", "__")
				// Duplicate name check for key.
				id = check_element_name(keyVar)
				if id != 0 {
					val_ += "__" + strconv.Itoa(int(id)) //string(id)
				}

				strSchema += "\n\t\t\t\"" + val_ + "\": &schema.Schema{\n\t\t\t\tType:    schema.TypeString,"
				strSchema += "\n\t\t\t\tOptional: true,"
				strSchema += "\n\t\t\t\tDescription:    \"xpath is: " + strStructHierarchy + "\",\n\t\t\t},"

				strStruct += "\n" + schemaTab + "V_" + val_ + "  *string  `xml:\"" + keyVar + ",omitempty\"`"
				strGetFunc += "\tV_" + val_ + " := d.Get(\"" + val_ + "\").(string)\n"
				strSetFunc += "\tif err :=d.Set(\"" + val_ + "\", " + strStructHierarchy + ".V_" + val_ + ");err != nil{\n"
				strSetFunc += "\t\treturn diag.FromErr(err)\n"
				strSetFunc += "\t}\n"
				strVarAssign += "\t" + strStructHierarchy + ".V_" + val_ + " = &V_" + val_ + "\n"
			}
		}
	}

	// TODO: Add the enums for the 'uses' to the element slice
	for keyVar, _ := range enums[uses] {

		if structXpath_last_elem != keyVar {
			val_ = s.ReplaceAll(keyVar, "-", "__")
			val_ = s.ReplaceAll(val_, ".", "__")
			// Duplicate name check for key.
			id = check_element_name(keyVar)
			if id != 0 {
				val_ += "__" + strconv.Itoa(int(id)) //string(id)
			}

			strSchema += "\n\t\t\t\"" + val_ + "\": &schema.Schema{\n\t\t\t\tType:    schema.TypeString,"
			strSchema += "\n\t\t\t\tOptional: true,"
			strSchema += "\n\t\t\t\tDescription:    \"xpath is: " + strStructHierarchy + "\",\n\t\t\t},"

			// This can be tricky to understand. Look at the generated code for more insight.
			// We use pointers for choice-ident enums, because Junos presents us empty elements.
			// Therefore, if the element isn't presented, the choice isn't configured.
			// So, we use pointers, if the element is presented, then the pointer is no longer nil
			// on unmarshal.
			strStruct += "\n" + schemaTab + "V_" + val_ + "  *string  `xml:\"" + keyVar + ",omitempty\"`"
			strGetFunc += "\tV_" + val_ + " := d.Get(\"" + val_ + "\").(string)\n"

			strSetFunc += "\tif " + strStructHierarchy + ".V_" + val_ + " != nil { \n\t\td.Set(\"" + val_ + "\", \" \") } else { d.Set(\"" + val_ + "\", \"\")}\n\n"

			strVarAssign += "\tif V_" + val_ + " != \"\" { " + strStructHierarchy + ".V_" + val_ + " = &V_" + val_ + " }\n"
		}

	}

	structXpath = ""

	return strStructHierarchy, structXpath, schemaTab
}

// initialize values of global variables used for generating the terraform module
func initializeFunctionString(name string) {
	name = s.ReplaceAll(name, "-", "__")
	name = s.ReplaceAll(name, ".", "__")
	strModuleName = name

	strVarAssign = "\n\tconfig := xml" + name + "{}\n"

	// fmt.Println("DEBUG: Check name 3: ", name)

	strStruct += "type xml" + name + " struct {\n\tXMLName xml.Name `xml:\"configuration\"`"

	if isGrpFlag {
		// strVarAssign += "\tconfig.ApplyGroup = id\n\tconfig.Groups.Name = id\n"
		strVarAssign += "\tconfig.Groups.Name = id\n"
		strStruct += "\n\tGroups  struct {\n\t\tXMLName\txml.Name\t`xml:\"groups\"`\n\t\tName\tstring\t`xml:\"name\"`"
		strStructEnd = "\n\t} `xml:\"groups\"`\n"
		// strStructEnd = "\n\t} `xml:\"groups\"`\n\tApplyGroup string `xml:\"apply-groups\"`"
	}
	strRead = "\n\tconfig := &xml" + name + "{}\n\n\terr = client.MarshalGroup(id, config)\n\tcheck(ctx, err)\n"

	// Append text for Create Function.
	strCreate += "func junos" + name + "Create" + strClientInit

	// Append text for Read Function.
	strRead = "\n\nfunc junos" + name + "Read" + strClientInit + strRead

	// Append text for Update Function.
	strUpdate = "\n\nfunc junos" + name + "Update" + strClientInit

	// Append text for Delete Function.
	strDelete = "\n\nfunc junos" + name + "Delete" + strClientInit + strDelete

	// Append text with function name to schema string.
	var schemaTemp string = "\n\nfunc junos" + name + "() *schema.Resource {\n\treturn &schema.Resource{\n"
	schemaTemp = schemaTemp + "\t\tCreateContext: junos" + name + "Create,\n"
	schemaTemp = schemaTemp + "\t\tReadContext: junos" + name + "Read,\n"
	schemaTemp = schemaTemp + "\t\tUpdateContext: junos" + name + "Update,\n"
	schemaTemp = schemaTemp + "\t\tDeleteContext: junos" + name + "Delete,"
	strSchema = schemaTemp + strSchema
}

// Parent Node xpath needs to be handled with closing braces appended after parsing sub-elements
func handleParentNodeXpath(nodes Node, strStructHierarchy string, schemaTab string) {
	val_ := s.ReplaceAll(nodes.Key, "-", "__")
	val_ = s.ReplaceAll(val_, ".", "__")
	id := check_element_name(nodes.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}

	// Initialization for the structure present at top of file.
	// fmt.Println("DEBUG: Check val_name 2: ", val_)
	// fmt.Println("DEBUG: Check nodes.Key: ", nodes.Key)
	strStruct += "\n" + schemaTab + "V_" + val_ + "  struct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodes.Key + "\"`"
	tab := schemaTab + "\t"
	strStructHierarchy += ".V_" + val_
	for _, n := range nodes.Nodes {
		// Simple gate whether to include children
		processInner := false
		// Here we need to see if the inner data is part of our explicit XPath expression.
		// Remember, these XPaths are the longest match and should not include siblings, if they're not included in the full path
		// fmt.Println("DEBUG is: ", n.Key)
		xPathParts := s.Split(currentXPath, "/")
		for _, v := range xPathParts {
			// fmt.Println("DEBUG: searching children...")
			if n.Key == v {
				processInner = true
				// fmt.Println("DEBUG: found inner child: ", n.Key, n.XMLName.Local, n.XMLName, n.XMLName.Space)
			}
		}

		if processInner {
			// If there is list or container, hierarchy needs to be added.
			// If leaf or leaf-list then hierarchy doesn't need to ne there.
			// If uses, then grouping needs to be resolved.
			if n.XMLName.Local == "container" || n.XMLName.Local == "list" {
				handleContainer(n, strStructHierarchy, tab)
			} else if n.XMLName.Local == "leaf" || n.XMLName.Local == "leaf-list" {
				handleLeaf(n, strStructHierarchy, tab)
			} else if n.XMLName.Local == "uses" {
				handleGrouping(n, strStructHierarchy, tab)
			}
		}
	}
	// Close the structure.
	strStruct += "\n" + schemaTab + "} `xml:\"" + structXpath + "\"`"
}

// handle the structure creation for the 'container'/'list' defined in yang files
func handleContainer(nodes Node, strStructHierarchy string, schemaTab string) {
	val_ := s.ReplaceAll(nodes.Key, "-", "__")
	val_ = s.ReplaceAll(val_, ".", "__")
	id := check_element_name(nodes.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}

	// this logic removes embedded apply-macro and apply-groups that are superfluious and potentially breaking
	if nodes.Key != "apply-macro" && nodes.Key != "apply-groups" {
		strStruct += "\n" + schemaTab + "V_" + val_ + "\tstruct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodes.Key + "\"`"
		tab := schemaTab + "\t"
		strStructHierarchy += ".V_" + val_
		for _, n := range nodes.Nodes {
			// If there is list or container, hierarchy needs to be added.
			// If leaf or leaf-list then hierarchy doesn't need to be there.
			// If uses, then grouping needs to be resolved.
			if n.XMLName.Local == "container" || n.XMLName.Local == "list" {
				handleContainer(n, strStructHierarchy, tab)
			}
			if n.XMLName.Local == "leaf" || n.XMLName.Local == "leaf-list" {
				handleLeaf(n, strStructHierarchy, tab)
			} else if n.XMLName.Local == "uses" {
				handleGrouping(n, strStructHierarchy, tab)
			}
		}
		strStruct += "\n" + schemaTab + "} `xml:\"" + nodes.Key + "\"`"
	}
}

// handle the structure creation for the 'leaf'/'leaf-list' defined in yang files
func handleLeaf(nodes Node, strStructHierarchy string, schemaTab string) {

	// This fixes the apply-groups and apply-groups-except nodes appearing in the providers.
	if nodes.Key != "apply-groups" && nodes.Key != "apply-groups-except" {
		var desc string
		// Extract description for the node.
		for _, n := range nodes.Nodes {
			if n.XMLName.Local == "description" {
				for _, n1 := range n.Nodes {
					if n1.XMLName.Local == "text" {
						desc = string([]byte(n1.Content))
					}
				}
			}
		}

		val_ := s.ReplaceAll(nodes.Key, "-", "__")
		val_ = s.ReplaceAll(val_, ".", "__")
		id := check_element_name(nodes.Key)
		if id != 0 {
			val_ += "__" + strconv.Itoa(int(id)) //string(id)
		}

		strStruct += "\n" + schemaTab + "V_" + val_ + "  *string  `xml:\"" + nodes.Key + ",omitempty\"`"
		strSchema += "\n\t\t\t\"" + val_ + "\": &schema.Schema{\n\t\t\t\tType:    schema.TypeString,"
		strSchema += "\n\t\t\t\tOptional: true,"
		strSchema += "\n\t\t\t\tDescription:    \"xpath is: " + strStructHierarchy + ". " + desc + "\",\n\t\t\t},"
		strGetFunc += "\tV_" + val_ + " := d.Get(\"" + val_ + "\").(string)\n"
		strSetFunc += "\tif err :=d.Set(\"" + val_ + "\", " + strStructHierarchy + ".V_" + val_ + ");err != nil{\n"
		strSetFunc += "\t\treturn diag.FromErr(err)\n"
		strSetFunc += "\t}\n"
		strVarAssign += "\t" + strStructHierarchy + ".V_" + val_ + " = &V_" + val_ + "\n"
	}
}

// For groups defined in yang file, it needs to be resolved to corresponding container/list
func handleGrouping(nodes Node, strStructHierarchy string, schemaTab string) {
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
		if n2.XMLName.Local == "container" || n2.XMLName.Local == "list" {
			handleContainer(n2, strStructHierarchy, schemaTab)
		}
		if n2.XMLName.Local == "leaf" || n2.XMLName.Local == "leaf-list" {
			handleLeaf(n2, strStructHierarchy, schemaTab)
		} else if n2.XMLName.Local == "uses" {
			handleGrouping(n2, strStructHierarchy, schemaTab)
		}
	}
}

// Function to extract uses and store them as groups.
// This is a pre-parser done on the file at the starting of the functionality.
func create_group_nodes(nodes []Node) {
	for _, n := range nodes {
		if n.XMLName.Local == "grouping" {
			grpNode = append(grpNode, n)
		}
		create_group_nodes(n.Nodes)
	}
}

// check what is the type of a particular node in the yang file
// the check is to verify if the particular node is required to be added as part of structure.
func check_node_tag(text string) bool {
	if text == "container" {
		return true
	} else if text == "list" {
		return true
	} else if text == "leaf" {
		return true
	} else if text == "leaf-list" {
		return true
	}
	return false
}

// In case of node-element with same name in yang file
// keep appending digits at the end of the variables
func check_element_name(text string) int {
	var cnt int = 0
	var iter int = 0
	for itr, n := range elementNameList {
		if n.name == text {
			cnt = n.count
			iter = itr
			break
		}
	}

	if cnt == 0 {
		var temp ElementName
		temp.name = text
		temp.count = 1
		elementNameList = append(elementNameList, temp)
	} else {
		elementNameList[iter].count += 1
	}
	return cnt
}

// Generate terraform Modules
func createFile(moduleFilePath string, jcfg cfg.Config) error {

	providerFileData += "\t\t\t\"junos-" + jcfg.ProviderName + "_" + strModuleName + "\": junos" + strModuleName + "(),\n"

	// Create go file with top container/module name.
	var fileName string = s.Join([]string{"resource", strModuleName}, "_")
	fileName = s.Join([]string{fileName, "go"}, ".")
	fileName = moduleFilePath + "/" + fileName
	fPtr, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fPtr.Close()

	// Append at the end of the schema which is at bottom of created file.
	strSchema += "\n\t\t},\n\t}\n}"

	// Append at end of structure which is near the top of the created file.
	strStruct += "\n}"
	// Append for the create function.
	strCreate += strGetFunc + "\n" + strVarAssign + strSendTrans + strSetIdValue +
		"\n\treturn junos" + strModuleName + "Read(ctx,d,m)" + "\n}"
	// Append for the update function.
	strUpdate += strGetFunc + "\n" + strVarAssign + strSendTransId +
		"\n\treturn junos" + strModuleName + "Read(ctx,d,m)" + "\n}"
	// Append for the read function.
	strRead += strSetFunc + "\n\treturn nil\n}"
	// Append for the delete function.
	strDelete += "\n\treturn nil\n}"

	// Write to the file.
	_, err = fPtr.WriteString(strImport)
	_, err = fPtr.WriteString(strStruct)
	_, err = fPtr.WriteString(strCreate)
	_, err = fPtr.WriteString(strRead)
	_, err = fPtr.WriteString(strUpdate)
	_, err = fPtr.WriteString(strDelete)
	_, err = fPtr.WriteString(strSchema)

	// Create terraform file
	createTerraform(strSchema, jcfg)

	fmt.Printf("Terraform API resource_%s created \n", strModuleName)
	xpathCounter += 1
	return nil
}

// copy file from source location to destination
func CopyFile(source string, dest string) {
	// Read the file.
	temp, _ := ioutil.ReadFile(source)
	ioutil.WriteFile(dest, temp, 0777)
}

// List files and get filenames.
func listFiles(yangFilePath string, jcfg cfg.Config) {

	os.Chdir(yangFilePath)

	out, err := exec.Command("ls").Output()
	if err != nil {
		log.Fatalf("%s", err)
	}

	output := string(out[:])
	lines := s.Split(output, "\n")
	for _, line := range lines {
		matched, _ := regexp.Match(`.yang`, []byte(line))
		if matched {
			temp := s.Split(line, ".yang")
			yang_file_list = append(yang_file_list, temp[0])
		}
	}

	providerFileData = `
// Copyright (c) 2017-2022, Juniper Networks Inc. All rights reserved.
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

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"terraform-provider-junos-vqfx/netconf"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
`

	providerFileData += "\nconst groupStrXML = `<load-configuration action=\"merge\" format=\"xml\">\n%s\n</load-configuration>`\n"

	providerFileData += "\nconst deleteStr = `<edit-config>\n\t<target>\n\t\t<candidate/>\n\t</target>\n\t<default-operation>none</default-operation>\n\t<config>\n\t\t<configuration>\n\t\t\t<groups operation=\"delete\">\n\t\t\t\t<name>%s</name>\n\t\t\t</groups>\n\t\t\t<apply-groups operation=\"delete\">%s</apply-groups>\n\t\t</configuration>\n\t</config>\n</edit-config>`\n"

	providerFileData += "\nconst commitStr = `<commit/>`\n"

	providerFileData += "\nconst getGroupXMLStr = `<get-configuration>\n\t<configuration>\n\t<groups><name>%s</name></groups>\n\t</configuration>\n</get-configuration>`\n"

	providerFileData += "\nconst ApplyGroupXML = `<load-configuration action=\"merge\" format=\"xml\">\n%s\n</load-configuration>`\n"

	providerFileData += "\ntype configuration struct {\n\tApplyGroup []string `xml:\"apply-groups\"`\n}\n"

	//providerFileData += "\ntype applyGroupID struct {\n\tApplyGroup string `xml:\"apply-groups\"`\n}"

	providerFileData += `
// var mockMap map[string]incomingConfig
var mockMapMutex sync.Mutex

// ProviderConfig is to hold client information
type ProviderConfig struct {
	netconf.Client
	Host string
}

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func check(ctx context.Context, err error) {
	if err != nil {
		// Some of these errors will be "normal".
		tflog.Debug(ctx, err.Error())
		f, _ := os.OpenFile("jtaf_logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(err.Error() + "\n")
		f.Close()
		return
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Host:     d.Get("host").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		SSHKey:   d.Get("sshkey").(string),
	}

	configFilePath, ok := os.LookupEnv("MOCK_FILE")
	var client netconf.Client
	var err error

	if ok {
		filePtr, err := os.OpenFile(configFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		client = FileClient{filePtr: filePtr}

	} else {
		client, err = config.Client()
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return &ProviderConfig{client, config.Host}, nil
}

var _ netconf.Client = &FileClient{}

// FileClient represents a fake client for testing purposes.
type FileClient struct {
	// You can add fields for testing purposes here.
	filePtr *os.File
}

// Close is a functional thing to close the FileClient (no-op in this case).
func (bc FileClient) Close() error {
	return nil
}

// updateRawConfig simulates updating the configuration on a network device.
func (bc FileClient) updateRawConfig(applyGroup string, netconfCall string, commit bool) (string, error) {
	// Simulate the update operation (you can customize this part).
	// Extract the string between <name> tags
	nameStart := strings.Index(netconfCall, "<name>")
	nameEnd := strings.Index(netconfCall, "</name>")
	if nameStart == -1 || nameEnd == -1 {
		return "", fmt.Errorf("Failed to extract the group name from the netconfCall")
	}
	groupName := netconfCall[nameStart+6 : nameEnd]

	// Add the groupName to the applyGroupsList
	addToApplyGroupsList(groupName)

	var groupString string
	groupString = fmt.Sprintf(groupStrXML, netconfCall)
	groupString = fmt.Sprintln("Group String: ", groupString)
	_, err := bc.filePtr.WriteString(groupString)
	if err != nil {
		return "", err
	}

	if commit {
		bc.filePtr.WriteString("\nCommiting from Update\n")
		_, err := bc.filePtr.WriteString(commitStr)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Updated config for group: %s", applyGroup), nil
}

// DeleteConfig simulates deleting a configuration on a network device.
func (bc FileClient) DeleteConfig(applyGroup string, commit bool) (string, error) {
	// Simulate the delete operation (you can customize this part).
	return fmt.Sprintf("Deleted config for group: %s", applyGroup), nil
}

// SendCommit simulates sending a commit to a network device.
func (bc FileClient) SendCommit() error {
	// Simulate the commit operation (you can customize this part).
	bc.sortApplyGroupsList()
	if err := bc.SendApplyGroups(); err != nil {
		return err
	}
	bc.filePtr.WriteString("\nCommiting from the SendCommit function\n")
	return nil
}

// MarshalGroup simulates retrieving and marshaling configuration data for a group.
func (bc FileClient) MarshalGroup(id string, obj interface{}) error {
	// Simulate the retrieval and marshaling of configuration data (you can customize this part).
	// For testing purposes, let's just marshal an example object and save it to a file.
	return nil
}

// SendTransaction simulates sending a transaction to a network device.
func (bc FileClient) SendTransaction(id string, obj interface{}, commit bool) error {
	// Simulate sending a transaction (you can customize this part).
	// For testing purposes, let's just write the transaction data to a file.
	cfg, err := xml.Marshal(obj) // Indent with four spaces
	if err != nil {
		return err
	}
	mockMapMutex.Lock()
	defer mockMapMutex.Unlock()

	// _, err = bc.filePtr.Write(append(cfg, byte('\n')))
	// if err != nil {
	// 	fmt.Printf("Error writing to XML file: %v\n", err)
	// 	return err
	// }

	// updateRawConfig deletes old group by, re-creates it then commits.
	// As far as Junos cares, it's an edit.
	if id != "" {
		// fmt.Println("Updating raw")
		if _, err = bc.updateRawConfig(id, string(cfg), commit); err != nil {
			return err
		}
		return nil
	}
	if _, err = bc.sendRawConfig(string(cfg), commit); err != nil {
		return err
	}
	return nil
}

// sendRawConfig is a wrapper for driver.SendRaw()
func (bc FileClient) sendRawConfig(netconfCall string, commit bool) (string, error) {

	// Extract the string between <name> tags
	nameStart := strings.Index(netconfCall, "<name>")
	nameEnd := strings.Index(netconfCall, "</name>")
	if nameStart == -1 || nameEnd == -1 {
		return "", fmt.Errorf("Failed to extract the group name from the netconfCall")
	}
	groupName := netconfCall[nameStart+6 : nameEnd]

	// Add the groupName to the applyGroupsList
	addToApplyGroupsList(groupName)

	groupString := fmt.Sprintf(groupStrXML, netconfCall)

	_, err := bc.filePtr.WriteString(groupString)
	if err != nil {
		return "", err
	}
	bc.filePtr.WriteString("\n\n")
	if commit {
		bc.filePtr.WriteString("\nCommiting from Sending\n")
		_, err := bc.filePtr.WriteString(commitStr)
		if err != nil {
			return "", err
		}
	}

	return "d", nil
}

// Helper function to add an id to the global list.
func addToApplyGroupsList(id string) {
	applyGroupsMutex.Lock()
	defer applyGroupsMutex.Unlock()
	applyGroupsList = append(applyGroupsList, id)
}

// Helper function to sort the global list.
func (bc FileClient) sortApplyGroupsList() {
	applyGroupsMutex.Lock()
	defer applyGroupsMutex.Unlock()

	// Create a map to track unique items
	uniqueGroups := make(map[string]bool)

	// Filter out empty strings and remove duplicates
	filteredGroups := make([]string, 0)
	for _, group := range applyGroupsList {
		if group != "" && !uniqueGroups[group] {
			uniqueGroups[group] = true
			filteredGroups = append(filteredGroups, group)
		}
	}

	// Sort the filtered list
	sort.Strings(filteredGroups)

	// Update the global applyGroupsList with the sorted and filtered list
	applyGroupsList = filteredGroups
}

var applyGroupsList []string
var applyGroupsMutex sync.Mutex

func (bc FileClient) SendApplyGroups() error {

	// Concatenate the strings in applyGroupsList.
	applyGroupsMutex.Lock()
	defer applyGroupsMutex.Unlock()

	var applyG configuration
	applyG.ApplyGroup = make([]string, len(applyGroupsList))
	for i, item := range applyGroupsList {
		applyG.ApplyGroup[i] = item
	}

	cfg, err := xml.Marshal(applyG)
	if err != nil {
		return err
	}

	_, err = bc.filePtr.WriteString("\n")
	if err != nil {
		return err
	}

	_, err = bc.filePtr.WriteString("\n\nApplying Groups\n")
	if err != nil {
		return err
	}

	applyGroupString := fmt.Sprintf(ApplyGroupXML, string(cfg))

	_, err = bc.filePtr.WriteString(applyGroupString)
	if err != nil {
		fmt.Printf("Error writing to XML file: %v\n", err)
		return err
	}

	return nil
}

// Provider returns a Terraform Provider.
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
	`
	providerFileData = strings.Replace(providerFileData, "[providerName]", jcfg.ProviderName, -1)
}

func copyDir(src, dst, extension string, fileCopyCount *uint32) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if !strings.Contains(info.Name(), extension) {
				return nil
			}
			atomic.AddUint32(fileCopyCount, 1)
		}
		fmt.Println(fmt.Sprintf("- Copying %s to %s", info.Name(), filepath.Base(dst)))
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))
		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}
		in, _ := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		fh.Chmod(info.Mode())
		_, err = io.Copy(fh, in)
		return err
	})
}

func printProgressBar(iteration, total int, prefix, suffix string, length int, fill string) {
	percent := float64(iteration) / float64(total)
	filledLength := int(length * iteration / total)
	end := ">"

	if iteration == total {
		end = "="
	}

	bar := strings.Repeat(fill, filledLength) + end + strings.Repeat("-", (length-filledLength))
	fmt.Printf("\r     %s [%s] %.2f%% %s", prefix, bar, percent*100, suffix)
	fmt.Println()
	fmt.Println()

	if iteration == total {
		fmt.Println()
	}
}

func createTerraformMain(jcfg cfg.Config) {

	// SETTING UP FILE PATH DIRECTORY AND NAMING
	// Get the absolute path of the executable binary
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Navigate up three directories
	targetDir = filepath.Join(executablePath, "..", "..", "..", "/TFtemplates")

	mainFileData :=
		`
# replace {text} with your own test setup

terraform {
	required_providers {
		junos-{device-type} = {
			source = "{input source path}"
			version = "{input version here}"
		}
	}
}

provider "junos-{device-type}" {
	host = "localhost"
	port = 8300
	username = "root"
	password = "juniper123"
	sshkey = ""
}

module "{test-folder-name}" {
	source = "./{test folder name}"

	providers = {junos-{device-type} = junos-{device-type}}

	depends_on = [junos-{device-type}_destroycommit.commit-main]
}


resource "junos-{device-type}_commit" "commit-main" {
	resource_name = "commit"
	depends_on = [module.{test-folder-name}]
}

resource "junos-{device-type}_destroycommit" "commit-main" {
	resource_name = "destroycommit"
}
	`

	mainFileData = strings.Replace(mainFileData, "{device-type}", jcfg.ProviderName, -1)

	//	DEFINE THE FILE CONTENT AND WRITE TO FILE
	fileContent := []byte(mainFileData)

	resultFile := "main.tf"

	// Define the file path where you want to save the .tf file
	filePath := filepath.Join(targetDir, resultFile)

	// Create and write the Terraform configuration file
	if err := ioutil.WriteFile(filePath, fileContent, 0644); err != nil {
		log.Fatalf("Error writing .tf file: %v", err)
	}
}

func createTerraformTest(jcfg cfg.Config) {

	testFileData :=
		`
terraform {
	required_providers {
		junos-{device-type} = {
			source = "{input source path}"
			version = "{input version here}"
		}
	}
}
`

	testFileData = strings.Replace(testFileData, "{device-type}", jcfg.ProviderName, -1)

	testFile := "test.tf"

	// Define the file path where you want to save the .tf file
	filePath := filepath.Join(targetDir, testFile)
	// var filePtr *os.File

	// // Open the file for writing (create it if it doesn't exist, truncate it if it does)
	// filePtr, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	// if err != nil {
	// 	return
	// }

	// Convert content to a byte slice and append a newline
	fileContent := []byte(testFileData + "\n")

	// Create and write the Terraform configuration file
	if err := ioutil.WriteFile(filePath, fileContent, 0644); err != nil {
		log.Fatalf("Error writing .tf file: %v", err)
	}
}

func createTerraform(strSchema string, jcfg cfg.Config) {

	// Split the first line by spaces
	elements := strings.Fields(strSchema)
	resource := ""

	// Find the name of the resource from the PROVIDER RESOURCE
	if len(elements) >= 2 {
		secondElement := elements[1]
		// Trim "junos" and parentheses from the string
		resource = strings.Trim(secondElement, "junos")
		resource = strings.Trim(resource, "()")
	} else {
		fmt.Println("Not enough elements on the first line.")
	}

	// 	FIND THE ARGUMENTS FROM THE SCEHEMA --> PARSE SCHEMA
	// Define a regular expression pattern to match the keys within double quotes --> This is to find INPUTS
	pattern := `"\w+":`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find all matches in the input string
	matches := re.FindAllString(strSchema, -1)

	// Create a slice to store the extracted keys
	var keys []string

	// Extract and print the keys without the quotes and colon
	for _, match := range matches {
		key := match[1 : len(match)-2] // Remove the double quotes and colon
		keys = append(keys, key)
	}

	// CREATING THE FILE CONTENTS
	providerType := "junos-" + jcfg.ProviderName
	resourceType := providerType + "_" + resource
	resourceName := "*replace with name for resource*"
	blockHead := "resource" + " \"" + resourceType + "\"" + " \"" + resourceName + "\" {"

	// headBlock := "provider " + "\"" + providerType + "\" {\n" + "	host = \"10.x.x.x\"\n	port = 22\n	username = \"\"\n	password = \"\"\n	sshkey = \"\"\n}\n"

	// Initialize an empty result string
	var argumentBlock string

	// Concatenate the strings with newline characters
	for _, str := range keys {
		argumentBlock += "	" + str + " = \"\"" + "\n"
	}
	argumentBlock += "}"

	finalTemplate := blockHead + "\n" + argumentBlock
	addToTF(finalTemplate)
	// finalTemplate = headBlock + "\n" + finalTemplate

	// //	DEFINE THE FILE CONTENT AND WRITE TO FILE
	// fileContent := []byte(finalTemplate)

	// resultFile := resource + ".tf"
	// // Define the file path where you want to save the .tf file
	// filePath := filepath.Join(targetDir, resultFile)

	// // Create and write the Terraform configuration file
	// if err := ioutil.WriteFile(filePath, fileContent, 0644); err != nil {
	// 	log.Fatalf("Error writing .tf file: %v", err)
	// }
}

func addToTF(content string) {

	testFile := "test.tf"

	// Define the file path where you want to save the .tf file
	filePath := filepath.Join(targetDir, testFile)

	//fmt.Println(filePath)
	fileContent := []byte(content + "\n")

	var filePtr *os.File

	filePtr, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer filePtr.Close()

	if _, err = filePtr.Write(append(fileContent, byte('\n'))); err != nil {
		log.Fatalf("Error writing .tf file: %v", err)
	}
}

// Struct to represent the XML structure
type FileList struct {
	XMLName xml.Name `xml:"file-list"`
	XPaths  []xpath  `xml:"xpath"`
}

type xpath struct {
	Name string `xml:"name,attr"`
}

func fixXPath_Inputs() {
	// Open and read the XML file

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}

	// Navigate back one folder
	for i := 0; i < 1; i++ {
		cwd = filepath.Dir(cwd)
	}

	// Append the file name to the path
	fileName := "xpath_inputs.xml"
	filePath := filepath.Join(cwd, fileName)

	xmlFile, err := os.OpenFile(filePath, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	// Decode the XML data
	decoder := xml.NewDecoder(xmlFile)
	var fileList FileList
	err = decoder.Decode(&fileList)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	// Create a map for fast lookup of issue XPaths
	issueXPathsMap := make(map[string]bool)
	for _, xpath := range issue_xpaths {
		issueXPathsMap[xpath] = true
	}

	// Truncate the file to remove its content
	xmlFile.Truncate(0)
	xmlFile.Seek(0, 0)

	// Write the file-list start tag
	xmlFile.WriteString("<file-list>\n")

	// Iterate through the XPaths and filter out the unwanted ones
	for _, xpath := range fileList.XPaths {
		if !issueXPathsMap[xpath.Name] {
			// If the XPath is not in the issue_xpaths list, write it to the same file
			xmlFile.WriteString("    ")
			xpathString := fmt.Sprintf("<xpath name=\"%s\"/>\n", xpath.Name)
			fmt.Fprintf(xmlFile, xpathString)
		}
	}

	// Write the file-list end tag
	xmlFile.WriteString("</file-list>")
}
