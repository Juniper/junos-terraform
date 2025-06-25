from setuptools import setup, find_packages
from setuptools.command.install import install

setup(
    name='junos-terraform',
    version='1.0.0',
    url='https://github.com/aburston/junos-terraform',
    author="Juniper Networks",
    description="Junos Terraform Framework",
    packages=find_packages(),
    install_requires=[
        'pytest',
        'pyang',
        'lxml',
        'jinja2'
    ]
)
