from setuptools import setup, find_packages
from setuptools.command.install import install

setup(
    name='junosterraform',
    version='1.0.0',
    url='https://github.com/aburston/junos-terraform',
    author="Juniper Networks",
    description="Junos Terraform Framework",
    #packages=find_packages(),
    packages=["junosterraform", "jtaf_pyang_plugin"],
    scripts=["junosterraform/jtaf-provider", "jtaf_pyang_plugin/jtaf-pyang-plugindir", "junosterraform/jtaf-junos-xml2tf"],
    package_data = {
        "junosterraform": [ "terraform_provider/*", "terraform_provider/netconf/*", "templates/*" ]
    },
    install_requires=[
        'pytest',
        'pyang',
        'lxml',
        'jinja2'
    ]
)
