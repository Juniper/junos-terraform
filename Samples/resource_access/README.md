## Sample Example

The config file used for this example -
config.toml

In this example we use yang file -  
* junos-qfx-conf-access@2019-01-01.yang

Three terraform api are to be created from these yang as mentioned in xpath_test.xml

The below yin and xpath files are generated from yang files.
* junos-qfx-conf-access@2019-01-01.yin
* junos-qfx-conf-access@2019-01-01_xpath.txt
* junos-qfx-conf-access@2019-01-01_xpath.xml

Corresponding provider and api are generated from this sample set. 
* resource_AccessDomainMapApply__MacroData.go
* resource_AccessGx__PlusPartition.go
* resource_AccessRadsecDestinationDynamic__Requests.go
* provider.go

