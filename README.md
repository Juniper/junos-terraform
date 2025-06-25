# JUNOS Terraform Automation Framework (JTAF)

Terraform is traditionally used for managing virtual infrastructure, but there are organisations out there that use Terraform end-to-end and also want to manage configuration state using the same methods for managing infrastructure. Sure, we can run a provisioner with Terraform, but that wasn't asked for!

Much the same as you can use Terraform to create an AWS EC2 instance, you can manage the configurational state of Junos. In essence, we treat Junos configuration as declarative resources.

So what is JTAF? It's a framework, meaning, it's an opinionated set of tools and steps that allow you to go from YANG models to a custom Junos Terraform provider. With all frameworks, there are some dependencies.

To use JTAF, you'll need machine that can run **Go, Python, Git and Terraform.** This can be Linux, OSX or Windows. Some easy to consume videos are below.

## Quick start

```bash
git clone https://github.com/aburston/junos-terraform
git clone https://github.com/juniper/yang
python -m venv venv
. venv/bin/activate
pip install ./junos-terraform
cd junos-terraform
pyang --plugindir ./junos_tf/pyang_plugin -f jtaf -p ../yang/18.2/18.2R3/common ../yang/18.2/18.2R3/junos-qfx/conf/*.yang > junos.json
cp examples/evpn-vxlan-dc/dc2/dc2-spine1.xml config.xml```
```
Edit `config.xml` to remove `<rpc>` and `<configuration>` tags.
Now run the following commands:
```bash
generate_plugin.py -j junos.json -x config.xml
populate_tf.py config.xml
cd terraform_providers
go build
go install
```
