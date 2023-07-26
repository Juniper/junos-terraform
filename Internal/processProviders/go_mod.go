package processProviders

const gomodcontent = `
module terraform-provider-junos-%+v

go 1.17

require (
        github.com/vinpatel24/go-netconf v0.1.5
	github.com/hashicorp/terraform-plugin-log v0.4.0
        github.com/hashicorp/terraform-plugin-sdk/v2 v2.16
)`
