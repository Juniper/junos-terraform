This Folder is a sample implementation of configuring the below configuration on device using jtaf and terraform. 

```
set policy-options policy-statement DEF-IMPORT-FWTRUST-TABLE term t1 from route-filter 10.0.0.0/17 exact
```

### Details - 
* module_file/xpath_sample.xml - The file has xpath used to generate the terraform provider
* module_file/resource_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress.go - the generated module.
* module_file/junos-es-conf-policy-options@2019-01-01.yang - the yang file used to generate the module.
* module_file/temp.toml -> the config file created to pass as argument to processYang and processProviders 
* test.tf - The file has api written to set the configuration. 
* terraform-provider-junos-device -> the generated binary 

### Understanding the implementation
1) processYang generates yin file and xpath using temp.toml. 
Related xpath for this configuration will be ``/policy-options/policy-statement/term/from/route-filter/address`` which can be found with help of generated xpath file.
Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
2) processProviders generates module using temp.toml. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
3) terraform provider is generated. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
4) Copy the generated provider and initialize it via terraform init api. 
5) create a testfile for execution. test.tf is the sample file used.


### How test.tf is written 
* We refer the generated modules for its creation. resource_Policy__OptionsPolicy__StatementTermFromRoute__FilterAddress.go is to be referred in this example. 
* The hierarchy used for the module generation is ``/policy-options/policy-statement/term/from/route-filter/address`` as mentioned in xpath_sample.xml file. 
* xmlPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress in the generated module gives the xml hierarchy as seen in the device. 
* junosPolicy__OptionsPolicy__StatementTermFromRoute__FilterAddress in the generated module provides the list of variables to be used to set the configuration. 
The variables may have duplicate names with incremental numbers attached at the end. 
The description of each variable provides their exact hierarchy which helps user to relate the variable to the xml schema.   
* The container or list names are not required to be mentioned in the test.tf file. 
* The key for each list also needs to be mentioned while setting a sub-element as it is required to translate to correct configuration.
In this example we have three lists ``policy-statement`` , ``term`` and ``route-filter`` in the hierarchy ``/policy-options/policy-statement/term/from/route-filter/address``. 
``name`` is the variable assigned to both of its key element ``policy-statement`` and ``term`` and they are differentiated by integers at the end and it needs to be mentioned in the test.tf
``route-filter`` has three keys as mentioned in yang ``choice__ident``, ``choice__value`` and ``address`` . 
Even though the hierarchy chosen here only sets ``address`` and ``choice-indent``, ``choice-value`` also needs to be set. We can set it as empty or leave it in test.tf so it will be sent as empty by terraform as default value. 
* The variable resource_name helps to identify if it is a group based configuration and is required. 
In this example in xpath_sample.xml we have mentioned ``group-flag`` as false so it will not be set in the database.


