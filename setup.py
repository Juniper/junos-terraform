from setuptools import setup, find_packages
from setuptools.command.install import install
import subprocess
import os
import shutil
import platform
import sys
from pathlib import Path

class CustomInstall(install):
    def run(self):
        # Run each step
        self.run_pyang()
        self.run_generate_plugin()
        self.go_build_and_install()
        # self.create_symlink(device, junos_version, os_type)
        # self.run_unit_tests()

        # Call default install
        install.run(self)

    def run_pyang(self):
        yang_dir = "yang_models"
        output_file = Path("junos_tf") / "combined.json"
        os.makedirs(output_file.parent, exist_ok=True)
        print(f"▶️ Running Pyang plugin on all YANG files in {yang_dir}...")

        # Build the command with all YANG files as arguments
        yang_files = list(Path(yang_dir).glob("*.yang"))
        if not yang_files:
            print("No YANG files found to process.")
            return
        
        cmd = [
            "pyang",
            "--plugindir", "junos_tf/pyang_plugin",
            "-f", "jtaf",
            *map(str, yang_files),
        ]

        # Run the command, redirecting stdout to the combined.json file
        with open(output_file, "w") as out_f:
            subprocess.run(cmd, stdout=out_f, check=True)

    def run_generate_plugin(self):
        print("🧹 Running generate_plugin.py script...")

        combined_json = Path("junos_tf") / "combined.json"
        test_xml = Path("junos_tf") / "test.xml"
        terraform_dir = Path("terraform_providers")
        os.makedirs(terraform_dir, exist_ok=True)

        cmd = [
            sys.executable,  # uses the current Python interpreter
            "junos_tf/generate_plugin.py",
            "-j", str(combined_json),
            "-x", str(test_xml)
        ]

        subprocess.run(cmd, check=True)

    def go_build_and_install(self):
        print("🏗️ Building Go provider...")
        cwd = os.getcwd()
        build_dir = os.path.join(cwd, "terraform_providers")
        subprocess.run(["go", "build"], cwd=build_dir, check=True)
        # subprocess.run(["go", "install"], cwd=build_dir, check=True)

    # def create_symlink(self, device, version, os_type):
    #     print("🔗 Creating symlink for Terraform...")
    #     home = str(Path.home())
    #     plugin_path = os.path.join(home, f".terraform.d/plugins/juniper/providers/junos-{device}/{version}/{os_type}")
    #     os.makedirs(plugin_path, exist_ok=True)

    #     built_provider = Path("terraform_providers") / "terraform_providers"  # replace with actual binary name if needed
    #     symlink_path = Path(plugin_path) / "terraform-provider-junos"
    #     if symlink_path.exists() or symlink_path.is_symlink():
    #         symlink_path.unlink()
    #     os.symlink(built_provider.resolve(), symlink_path)

    # def run_unit_tests(self):
    #     print("🧪 Running unit tests...")
    #     subprocess.run(["python3", "-m", "pytest", "unit_tests"], check=True)

setup(
    name='junos-terraform-generator',
    version='0.1',
    packages=find_packages(),
    include_package_data=True,
    # cmdclass={'install': CustomInstall},
    install_requires=[
        'pytest',
        'pyang',
        'lxml',
    ],
    entry_points={
        'console_scripts': [
            "generate_plugin=junos_tf.generate_plugin:main",
            "populate_tf=junos_tf.populate_tf:main",
        ]
    },
    package_data={
        "junos_tf": ["pyang_plugin/*.py"],
    },
    python_requires='>=3.7',
    zip_safe=False,
    classifiers=[
        'Programming Language :: Python :: 3',
        'Operating System :: OS Independent',
    ],
)
