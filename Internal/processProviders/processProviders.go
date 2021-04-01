// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
// * Neither the name of the Juniper Networks nor the
//   names of its contributors may be used to endorse or promote products
//   derived from this software without specific prior written permission.
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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	s "strings"

	"github.com/Juniper/junos-terraform/Internal/cfg"
)

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

// Create variable to store groups in the yin file.
var grpNode []Node

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

// String variable for structure.
var strStruct string
var strStructEnd string

//string variable for schema.
var strSchema string
var strClientInit string
var strSendTrans string
var strSendTransId string
var strSetIdValue string
var strClientClose string

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

// Syntactic helper to reduce repetition.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CreateProviders consumes a YANG file, Xpath file and module name to create a provider
func CreateProviders(jcfg cfg.Config) error {

	yangFilePath := jcfg.YangDir
	xpathFilePath := jcfg.XpathPath
	moduleFilePath := jcfg.ProviderDir

	// Create list of yang files present
	listFiles(yangFilePath)

	// Read data from file
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
	for _, n5 := range inNode.Nodes {
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

			// Start processing of the file data
			// Notes : "-" and "." is not allowed in go as variable name. need to replace it with "_"
			start([]Node{n})

			// After all the data processing is done, create the file.
			createFile(moduleFilePath)
		}
	}

	providerFileData += `			"junos-qfx_commit": junosQFXCommit(),
		},
		ConfigureFunc: returnProvider,
	}
}
`
	// Create provider.go file
	var fileName string = "provider.go"
	fileName = moduleFilePath + "/" + fileName
	fPtr, err := os.Create(fileName)
	check(err)

	// Write to the file
	_, err = fPtr.WriteString(providerFileData)

	// No errors, so return nil.
	return nil
}

// Initialize global variables to default values
func initialize_global_variables() {

	grpNode = []Node{}

	elementNameList = []ElementName{}

	strModuleName = ""

	structXpath = ""

	strImport = `
package main

import (
    "encoding/xml"
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	strClientInit = `(d *schema.ResourceData, m interface{}) error {

    pcfg := m.(*ProviderConfig)
    client, err := pcfg.Client()
    check(err)

    id := d.Get("resource_name").(string)
    `

	strSendTrans = `
    err = client.SendTransaction("", config, commit)
    check(err)
    `

	strSendTransId = `
    err = client.SendTransaction(id, config, commit)
    check(err)
    `

	strSetIdValue = `
    d.SetId(fmt.Sprintf("%s_%s", pcfg.Cfg.Host, id))
    `

	strClientClose = `
    err = client.Close()
    check(err)
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
    _, err = client.DeleteConfig(id)
    check(err)

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
				if n2.XMLName.Local == "container" {
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
	// remove the 1st / from the xapth as it is not required
	// for the rest xpath replace / with > as it will be used in struct.
	strXpath := inputXpath
	strXpath = s.Replace(strXpath, "/", "", 1)
	// structXpath = s.ReplaceAll(structXpath, "/", ">")

	strParts := s.Split(strXpath, "/")
	structXpath = strParts[0]

	nodeCheck := nodes
	schemaTab := "\t"
	var strStructHierarchy string = "config"

	if isGrpFlag == true {
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

	// If the topmost container is chosen, don't process xpath.
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
				}

				if matchFound {
					if itr == (len(strParts) - 1) {
						structXpath_last_elem = strParts[itr]
					} else {
						if structXpath == "" {
							structXpath = strParts[itr]
						} else {
							structXpath += ">" + strParts[itr]
						}
						// Handle for list and create parameter for key.
						if nodeCheck.XMLName.Local == "list" {
							strStructHierarchy, structXpath, schemaTab = setListXpathMatch(nodeCheck, schemaTab, structXpath, strStructHierarchy)
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
			if matchFound == false {
				fmt.Printf("Xpath not found in file, check it. : %s \n", strXpath)
				// os.Exit(0)
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
				val_ = s.ReplaceAll(node_last_elemt_2.Key, ".", "__")
				id := check_element_name(node_last_elemt_2.Key)
				if id != 0 {
					val_ += "__" + strconv.Itoa(int(id)) //string(id)
				}
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
		}
	}
	return nodeCheck, flag
}

// Handle xpath matching for list. Need to add key also in cae of list during xpath matching
func setListXpathMatch(nodeCheck Node, schemaTab string, structXpath string, strStructHierarchy string) (string, string, string) {
	var keyValue string
	for _, n1 := range nodeCheck.Nodes {
		if n1.XMLName.Local == "key" {
			keyValue = n1.Val
			break
		}
	}

	// Assign values for list and its key values.
	val_ := s.ReplaceAll(nodeCheck.Key, "-", "__")
	val_ = s.ReplaceAll(nodeCheck.Key, ".", "__")
	// Duplicate name check for list.
	id := check_element_name(nodeCheck.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}
	strStruct += "\n" + schemaTab + "V_" + val_ + "  struct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodeCheck.Key + "\"`"
	strStructEnd = "\n" + schemaTab + "} `xml:\"" + structXpath + "\"`" + strStructEnd
	strStructHierarchy += ".V_" + val_
	schemaTab += "\t"

	val_ = s.ReplaceAll(keyValue, "-", "__")
	val_ = s.ReplaceAll(keyValue, ".", "__")
	// Duplicate name check for key.
	id = check_element_name(keyValue)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}

	strSchema += "\n\t\t\t\"" + val_ + "\": &schema.Schema{\n\t\t\t\tType:    schema.TypeString,"
	strSchema += "\n\t\t\t\tOptional: true,"
	strSchema += "\n\t\t\t\tDescription:    \"xpath is: " + strStructHierarchy + "\",\n\t\t\t},"
	strStruct += "\n" + schemaTab + "V_" + val_ + "  string  `xml:\"" + keyValue + "\"`"
	strGetFunc += "\tV_" + val_ + " := d.Get(\"" + val_ + "\").(string)\n"
	strSetFunc += "\td.Set(\"" + val_ + "\", " + strStructHierarchy + ".V_" + val_ + ")\n"
	strVarAssign += "\t" + strStructHierarchy + ".V_" + val_ + " = V_" + val_ + "\n"

	structXpath = ""

	return strStructHierarchy, structXpath, schemaTab
}

// initialize values of global variables used for generating the terraform module
func initializeFunctionString(name string) {
	name = s.ReplaceAll(name, "-", "__")
	name = s.ReplaceAll(name, ".", "__")
	strModuleName = name

	strVarAssign = "\n\tconfig := xml" + name + "{}\n"

	strStruct += "type xml" + name + " struct {\n\tXMLName xml.Name `xml:\"configuration\"`"

	if isGrpFlag == true {
		strVarAssign += "\tconfig.ApplyGroup = id\n\tconfig.Groups.Name = id\n"
		strStruct += "\n\tGroups  struct {\n\t\tXMLName\txml.Name\t`xml:\"groups\"`\n\t\tName\tstring\t`xml:\"name\"`"
		strStructEnd = "\n\t} `xml:\"groups\"`\n\tApplyGroup string `xml:\"apply-groups\"`"
	}
	strRead = "\n\tconfig := &xml" + name + "{}\n\n\terr = client.MarshalGroup(id, config)\n\tcheck(err)\n"

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
	schemaTemp = schemaTemp + "\t\tCreate: junos" + name + "Create,\n"
	schemaTemp = schemaTemp + "\t\tRead: junos" + name + "Read,\n"
	schemaTemp = schemaTemp + "\t\tUpdate: junos" + name + "Update,\n"
	schemaTemp = schemaTemp + "\t\tDelete: junos" + name + "Delete,"
	strSchema = schemaTemp + strSchema
}

// Parent Node xpath needs to be handled with closing braces appended after parsing sub-elements
func handleParentNodeXpath(nodes Node, strStructHierarchy string, schemaTab string) {
	val_ := s.ReplaceAll(nodes.Key, "-", "__")
	val_ = s.ReplaceAll(nodes.Key, ".", "__")
	id := check_element_name(nodes.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}
	// Initialization for the structure present at top of file.
	strStruct += "\n" + schemaTab + "V_" + val_ + "  struct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodes.Key + "\"`"
	tab := schemaTab + "\t"
	strStructHierarchy += ".V_" + val_
	for _, n := range nodes.Nodes {
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
	// Close the structure.
	strStruct += "\n" + schemaTab + "} `xml:\"" + structXpath + "\"`"
}

// handle the structure creation for the 'container'/'list' defined in yang files
func handleContainer(nodes Node, strStructHierarchy string, schemaTab string) {
	val_ := s.ReplaceAll(nodes.Key, "-", "__")
	val_ = s.ReplaceAll(nodes.Key, ".", "__")
	id := check_element_name(nodes.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}
	strStruct += "\n" + schemaTab + "V_" + val_ + "\tstruct {\n" + schemaTab + "\tXMLName xml.Name `xml:\"" + nodes.Key + "\"`"
	tab := schemaTab + "\t"
	strStructHierarchy += ".V_" + val_
	for _, n := range nodes.Nodes {
		// If there is list or container, hierarchy needs to be added.
		// If leaf or leaf-list then hierarchy doesn't need to ne there.
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

// handle the structure creation for the 'leaf'/'leaf-list' defined in yang files
func handleLeaf(nodes Node, strStructHierarchy string, schemaTab string) {
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
	val_ = s.ReplaceAll(nodes.Key, ".", "__")
	id := check_element_name(nodes.Key)
	if id != 0 {
		val_ += "__" + strconv.Itoa(int(id)) //string(id)
	}

	strStruct += "\n" + schemaTab + "V_" + val_ + "  string  `xml:\"" + nodes.Key + "\"`"
	strSchema += "\n\t\t\t\"" + val_ + "\": &schema.Schema{\n\t\t\t\tType:    schema.TypeString,"
	strSchema += "\n\t\t\t\tOptional: true,"
	strSchema += "\n\t\t\t\tDescription:    \"xpath is: " + strStructHierarchy + ". " + desc + "\",\n\t\t\t},"
	strGetFunc += "\tV_" + val_ + " := d.Get(\"" + val_ + "\").(string)\n"
	strSetFunc += "\td.Set(\"" + val_ + "\", " + strStructHierarchy + ".V_" + val_ + ")\n"
	strVarAssign += "\t" + strStructHierarchy + ".V_" + val_ + " = V_" + val_ + "\n"
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
func createFile(moduleFilePath string) {

	providerFileData += "\t\t\t\"junos-qfx_" + strModuleName + "\": junos" + strModuleName + "(),\n"

	// Create go file with top container/module name.
	var fileName string = s.Join([]string{"resource", strModuleName}, "_")
	fileName = s.Join([]string{fileName, "go"}, ".")
	fileName = moduleFilePath + "/" + fileName
	fPtr, err := os.Create(fileName)
	check(err)

	// Append at the end of the schema which is at bottom of created file.
	strSchema += "\n\t\t},\n\t}\n}"

	// Append at end of structure which is near the top of the created file.
	strStruct += "\n}"
	// Append for the create function.
	strCreate += strGetFunc + "\tcommit := true\n" + strVarAssign + strSendTrans + strSetIdValue +
		strClientClose + "\n\treturn junos" + strModuleName + "Read(d,m)" + "\n}"
	// Append for the update function.
	strUpdate += strGetFunc + "\tcommit := true\n" + strVarAssign + strSendTransId + strClientClose +
		"\n\treturn junos" + strModuleName + "Read(d,m)" + "\n}"
	// Append for the read function.
	strRead += strSetFunc + strClientClose + "\n\treturn nil\n}"
	// Append for the delete function.
	strDelete += strClientClose + "\n\treturn nil\n}"

	// Write to the file.
	_, err = fPtr.WriteString(strImport)
	_, err = fPtr.WriteString(strStruct)
	_, err = fPtr.WriteString(strCreate)
	_, err = fPtr.WriteString(strRead)
	_, err = fPtr.WriteString(strUpdate)
	_, err = fPtr.WriteString(strDelete)
	_, err = fPtr.WriteString(strSchema)
	fmt.Printf("Terraform API resource_%s created \n", strModuleName)
}

// copy file from source location to destination
func CopyFile(source string, dest string) {
	// Read the file.
	temp, _ := ioutil.ReadFile(source)
	ioutil.WriteFile(dest, temp, 0777)
}

// List files and get filenames.
func listFiles(yangFilePath string) {

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
// Copyright (c) 2017-2021, Juniper Networks Inc. All rights reserved.
//
// License: Apache 2.0
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
// * Neither the name of the Juniper Networks nor the
//   names of its contributors may be used to endorse or promote products
//   derived from this software without specific prior written permission.
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
	"log"

	gonetconf "github.com/davedotdev/go-netconf/helpers/junos_helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProviderConfig is to hold client information
type ProviderConfig struct {
	Cfg Config
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Client Refreshes the client
func (PC *ProviderConfig) Client() (*gonetconf.GoNCClient, error) {

	client, err := PC.Cfg.Client()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}

func returnProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:     d.Get("host").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		SSHKey:   d.Get("sshkey").(string),
	}

	return &ProviderConfig{config}, nil
}

// Provider returns a Terraform ResourceProvider.
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
}
