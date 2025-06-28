from setuptools import setup, find_packages
from setuptools.command.install import install

setup(
    name='junosterraform',
    version='1.0.0',
    url='https://github.com/aburston/junos-terraform',
    author="Juniper Networks",
    description="Junos Terraform Framework",
    packages=find_packages(),
    scripts=["junosterraform/generate_provider.py", "junosterraform/populate_tf.py"],
    install_requires=[
        'pytest',
        'pyang',
        'lxml',
        'jinja2'
    ]
)
