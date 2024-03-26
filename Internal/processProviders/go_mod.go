package processProviders

const gomodcontent = `
module terraform-provider-junos-%+v

go 1.22

require (
        github.com/chrismarget/lambda-tf-registry v0.0.1
        github.com/goreleaser/goreleaser v1.24.0
        github.com/hashicorp/terraform-plugin-log v0.4.0
        github.com/hashicorp/terraform-plugin-sdk/v2 v2.16.0
        github.com/vinpatel24/go-netconf v0.1.3
        golang.org/x/crypto v0.18.0
)`
