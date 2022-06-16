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

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/Juniper/junos-terraform/Internal/cfg"
	"github.com/Juniper/junos-terraform/Internal/processYang"
)

const _ver = "0.1.5"

// Syntactic helper to reduce repetition.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Check if running from venv
func check_venv_exists() {
	app	:= "python3"
	script	:= "checkVenv.py"
	cmd := exec.Command(app, script)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	if string(stdout) == "false\n" {
		fmt.Println("ERROR: Please run this in a python3 virtual environment.\n")
		os.Exit(1)
	}
}

func check_pyang_installed() {
	app := "pyang"
	ver := "-v"
	cmd := exec.Command(app, ver)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("ERROR: Please install pyang in the virtual environment.\n")
		os.Exit(1)
	}
}

func PrintLogo() {
	const jtafLogo = `
                    ___ _____ ___  ______ 
                   |_  |_   _/ _ \ |  ___|
                     | | | |/ /_\ \| |_   
                     | | | ||  _  ||  _|  
                 /\__/ / | || | | || |    
                 \____/  \_/\_| |_/\_| `

	fmt.Println(jtafLogo)
	fmt.Println("                           " + _ver)
}

func main() {

	// The user can pass the configuration as part of config file or command line arguments.
	jcfg := cfg.Config{}

	// Version flag
	verFlag := flag.Bool("v", false, "Check the version")

	// Get the config location.
	flagConfig := flag.String("config", "", "Path to config file")

	// Get flags from the user for JTAF.
	flagYang := flag.String("yangDir", "", "Absolute path to Yang files directory")
	flagFileType := flag.String("fileType", "", "fileType for the xpath to be generated for Providers")

	flag.Parse()

	// Check the version.
	if *verFlag {
		fmt.Println("v", _ver)
		return
	}

	// Get the config.
	if *flagConfig != "" {
		jcfg, err := cfg.GetConfig(*flagConfig)
		if err != nil {
			fmt.Println("Error retrieving configuration: ", err)
		}

		PrintLogo()

		check_venv_exists()
		check_pyang_installed()
		check(processYang.CreateYinFileAndXpath(jcfg))
	} else if *flagYang != "" || *flagFileType != "" {
		// If config file path is not present then check for individual elements.
		jcfg.FileType = *flagFileType
		jcfg.YangDir = *flagYang

		PrintLogo()

		check_venv_exists()
		check_pyang_installed()
		check(processYang.CreateYinFileAndXpath(jcfg))
	} else {
		fmt.Println("One or more mandatory inputs are missing, exiting...")
		os.Exit(0)
	}
}
