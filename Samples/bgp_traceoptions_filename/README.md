This Folder is a sample implementation of configuring the below configuration on device using jtaf and terraform. 
Filename(bgp.log) is a leaf element and we are generating module at the same hierarchy. 

```set protocols bgp group <name> traceoptions file bgp.log```


### Details of Files present in the directory 
* module_file/junos-qfx-conf-protocols@2019-01-01.yang -> the yang file used to generate the module.
* module_file/xpath_sample.xml -> The file has xpath used to generate the terraform provider for xpath /protocols/bgp/group/traceoptions/file/filename
* module_file/temp.toml -> the config file created to pass as argument to processYang and processProviders 
* module_file/resource_ProtocolsBgpGroupTraceoptionsFileFilename.go -> the generated module.
* test.tf -> The file has api written to set the configuration. 
* terraform-provider-junos-device -> the generated binary 


### Understanding the implementation
1) processYang generates yin file and xpath using temp.toml. 
Related xpath for this configuration will be ``/protocols/bgp/group/traceoptions/file/filename`` which can be found with help of generated xpath file.
Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
2) processProviders generates module using temp.toml. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
3) terraform provider is generated. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
4) Copy the generated provider and initialize it via terraform init api. 
5) create a testfile for execution. test.tf is the sample file used.


### How test.tf is written 
* We refer the generated modules for its creation. resource_ProtocolsBgpGroupTraceoptionsFileFilename.go is to be referred in this example. 
* The hierarchy used for the module generation is /protocols/bgp/group/traceoptions/file/filename as mentioned in xpath_sample.xml file. 
* xmlProtocolsBgpGroupTraceoptionsFileFilename in the generated module gives the xml hierarchy as seen in the device. 
* junosProtocolsBgpGroupTraceoptionsFileFilename in the generated module provides the list of variables to be used to set the configuration. 
The variables may have duplicate names with incremental numbers attached at the end. 
The description of each variable provides their exact hierarchy which helps user to relate the variable to the xml schema.   
* The container or list names are not required to be mentioned in the test.tf file. 
* The key for each list also needs to be mentioned while setting a sub-element as it is required to translate to correct configuration.
In this example we have a list ``group`` in the hierarchy ``/protocols/bgp/group/traceoptions/file/filename``. 
``name`` is the variable assigned to its key element and it needs to be mentioned in the test.tf
* The variable ``filename`` is assigned to leaf-element ``filename`` and needs to be mentioned in the test.tf .
* The variable resource_name helps to identify if it is a group based configuration and is required. 
In this example in xpath_sample.xml we have mentioned ``group-flag`` as false so it will not be set in the database.  