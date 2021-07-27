This Folder is a sample implementation of configuring the below configuration on device using jtaf and terraform. 
This is a much complex implementation. Refer ``security_address-set`` for simple and detailed explanation as starting point for this.

```
set security nat source pool ut-pool address 107.127.96.40/32 to 107.127.96.43/32
set security nat source pool t-pool address 172.16.0.8/32 to 172.16.0.10/32
set security nat source rule-set snat-untrust-to-trust from zone untrust
set security nat source rule-set snat-untrust-to-trust to zone trust
set security nat source rule-set snat-untrust-to-trust rule 1 match source-address 0.0.0.0/0
set security nat source rule-set snat-untrust-to-trust rule 1 then source-nat pool t-pool
set security nat source rule-set snat-trust-to-untrust from zone trust
set security nat source rule-set snat-trust-to-untrust to zone untrust
set security nat source rule-set snat-trust-to-untrust rule 2 match source-address 10.0.0.0/17
set security nat source rule-set snat-trust-to-untrust rule 2 then source-nat pool ut-pool
```

### Details - 
* module_file/xpath_sample.xml - The file has xpath used to generate the terraform provider
* module_file/resource_*.go - the generated modules.
* module_file/junos-es-conf-security@2019-01-01.yang - the yang file used to generate the module.
* module_file/temp.toml -> the config file created to pass as argument to processYang and processProviders 
* test.tf - The file has api written to set the configuration. 
* terraform-provider-junos-device -> the generated binary 

### Understanding the implementation
1) processYang generates yin file and xpath using temp.toml. 
Related xpath for this configuration is mentioned in module_file/xpath_sample.xml which can be found with help of generated xpath file.
Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
2) processProviders generates module using temp.toml. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
3) terraform provider is generated. Refer [link](https://github.com/Juniper/junos-terraform/blob/master/README.md)
4) Copy the generated provider and initialize it via terraform init api. 
5) create a testfile for execution. test.tf is the sample file used.


### How test.tf is written 
* We are generating five different modules for this example. All of them belong to the hierarchy ``/security/nat/source/``.
* If we check the generated xpath file, we will find that there are around 300+ sub-elements in this hierarchy but we need to set only 5 of them. 
* The 1st concern with this which makes it mandatory to be generated at leaf-level is that it has choice in its yang. 
JTAF can't resolve it itself, which needs to be set so the user needs to provide the exact element(/choice/case/element) in the hierarchy. 
Any parent hierarchy having choice as a child element will create issues as jtaf can't resolve it while generation. 
* The 2nd concern would be that with 300 elements there may be syntactic check which will fail as by default they will be passed as an xml with empty tags for those elements.
* The 3rd concern is that some of them may be lists and there key can't be empty so the user will need to fill information for these elements as well which is not required. 
check ``bgp_traceoptions_file`` for a simpler example and detailed explanation regarding this.  
