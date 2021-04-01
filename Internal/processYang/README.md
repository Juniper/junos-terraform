# processYang

This module is to be used to generate yin files from yang files for which terraform api's are to be generated. 

## Yang file Details to be checked

Check the correct version and device type for the yang files. Yang files for junos devices can be downloaded from the 
following github repo -

``https://github.com/Juniper/yang.git`` 

## Generate yin files and xpaths from the yang files

Yin files are an xml representation of the yang files. It makes it easier for the golang code to understand the data. 
All the xpaths are generated for a corresponding yang file. User can refer it while generating terraform api for the 
yang model. 
Refer following links for more details regarding yin files and xpaths 

* [RFC 6020 : Yin](https://tools.ietf.org/html/rfc6020#section-11)
* [XPath from W3C](https://www.w3.org/TR/1999/REC-xpath-19991116/)

To generate yin file and xpath follow the following steps. 
1) Copy yang files to a particular repository on the device. 
2) Execute the following command -
 
`` cd cmd/processyang `` 

`` go build ``

`` ./processYang -config /var/tmp/config.toml ``

The yin-files and xpath-files will be generated in the same repository.

The Sample config file are provided in jtaf/Samples/config.toml . 
The same arguments can also be passed from command line during binary execution . 

