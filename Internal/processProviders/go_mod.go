package processProviders

const gomodcontent = `
module terraform-provider-junos-%+v

go 1.14

require (
        github.com/davedotdev/go-netconf v0.1.5
        github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
)`
