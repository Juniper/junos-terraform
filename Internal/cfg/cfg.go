// This file contains the config file information for storing data about JTAF to reduce input repetition.
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

// This file contains the config file information for storing data about JTAF to reduce input repetition
package cfg

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config used to decode configuration file.
type Config struct {
	YangDir      string `toml:"yangDir"`
	ProviderDir  string `toml:"providerDir"`
	XpathPath    string `toml:"xpathPath"`
	FileType     string `toml:"fileType"`
	ProviderName string `toml:"providerName"`
}

// GetConfig loads the config from a .TOML file and removes the problem of repetition.
func GetConfig(name string) (Cfg Config, err error) {
	c := Config{}
	_, err = toml.DecodeFile(name, &c)

	if err != nil {
		return c, err
	}

	// Let's check for empty fields.
	emptyField := ""

	if c.YangDir == "" {
		emptyField = "yangDir"
	} else if c.ProviderDir == "" {
		emptyField = "providerDir"
	} else if c.XpathPath == "" {
		emptyField = "xpathPath"
	}

	// This is done as an Errorf type call. It's an error of our workflow but we're returning an error to main().
	if emptyField != "" {
		err := fmt.Errorf("empty required config field missing: %s", emptyField)
		return c, err
	}

	if c.FileType == "" {
		c.FileType = "text"
	}

	return c, nil
}
