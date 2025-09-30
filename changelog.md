
# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
 
## [Unreleased]

### Added

### Changed
 
### Fixed

## [1.1.0] - yyyy-mm-dd
 
### Added
- jtaf-yang2go command to combine yang files to JSON conversion and provider creation ([#6](https://github.com/Juniper/junos-terraform/issues/6))

### Changed
- jtaf-provider, jtaf-xml2tf, and jtaf-yang2go to accept multiple xml configurations of the same device type ([#72](https://github.com/Juniper/junos-terraform/issues/72))
- jtaf-xml2tf to support a base configuration, groups, and apply groups ([#65](https://github.com/Juniper/junos-terraform/issues/65))
 
### Fixed
- Dependency on private go-netconf repository ([#61](https://github.com/Juniper/junos-terraform/issues/61))
- Unexpected output rpc-reply messages ([#71](https://github.com/Juniper/junos-terraform/issues/71))
- Leaf-list error in jtaf-xml2tf ([#65](https://github.com/Juniper/junos-terraform/issues/65))
 
## [1.0.0] - 2025-07-01

### Added
- Many updates to make JTAF production ready ([Release 1.0.0](https://github.com/Juniper/junos-terraform/releases/tag/1.0.0))

## [0.1.1] - 2025-06-26

### Added
- Many updates and examples ([Release 0.1.1](https://github.com/Juniper/junos-terraform/releases/tag/0.1.1))

## [0.1] - 2021-04-14

### Added
- First release of API to generate Junos modules for Terraform ([Release 0.1](https://github.com/Juniper/junos-terraform/releases/tag/0.1))
 