from setuptools import setup, find_packages
from setuptools.command.install import install

setup(
    name='junosterraform',
    version='1.0.1',
    url='https://github.com/aburston/junos-terraform',
    tests_require=["pytest"],
    test_suite="junosterraform.unit_tests",
    author="Juniper Networks",
    description="Junos Terraform Framework",
    packages=["junosterraform", "jtaf_pyang_plugin", "terraform_provider"],
    scripts=["junosterraform/jtaf-provider","junosterraform/jtaf-ansible","junosterraform/jtaf-xml2yaml", "jtaf_pyang_plugin/jtaf-pyang-plugindir", "junosterraform/jtaf-xml2tf", "junosterraform/jtaf-yang2go",
             "junosterraform/jtaf-yang2ansible"],
    package_data = {
        "terraform_provider": [
            "*",
            "netconf/*",
            "go-netconf/*",
            "go-netconf/drivers/driver/*",
            "go-netconf/drivers/junos/*",
            "go-netconf/drivers/junos/lowlevel/*",
            "go-netconf/drivers/ssh/*",
            "go-netconf/drivers/ssh/lowlevel/*",
            "go-netconf/helpers/junos_helpers/*",
            "go-netconf/rpc/*",
            "go-netconf/session/*",
            "go-netconf/transport/*"
        ],
        "junosterraform": [ 
            "templates/*",
            "jtaf_common.py"
            ]
    },
    install_requires=[
        'pytest',
        'pyang',
        'lxml',
        'jinja2',
        'setuptools',
        'pyyaml'
    ]
)
