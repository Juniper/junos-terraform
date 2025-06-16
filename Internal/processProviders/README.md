# processProviders

The framework generates junos compatible terraform providers based on the input provided by the user.
The provider remains independent of junos version and device type. The user needs to provide the device compatible yang
files and the xpaths for which the provider is to be generated. 

## Steps to generate the providers 

### Pre-Requisite

Generate the yin file and xpath for the yang files with processYang module.

### Steps 

`` cd cmd/processyang `` 

`` go build ``

`` ./processProviders -config /var/tmp/config.toml ``

Update the xpath_list_to_be_generated.xml with the yin files, xpaths and group details.

The terraform api will be generated in the repository as provided inn config file .

The Sample config file are provided in jtaf/Samples/config.toml . 
The same arguments can also be passed from command line during binary execution . 