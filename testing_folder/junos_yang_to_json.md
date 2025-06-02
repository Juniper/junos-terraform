## Coverting Junos YANG Models to JSON

Below are the steps to convert Junos YANG models into respective JSON data structure, these steps have been verified on Ubuntu OS 22.04

1. Clone Junos YANG Models repository
      `https://github.com/Juniper/yang`
2. Install pyang python library
      `pip3 install pyang`
3. Clone below repository to get JTAF plugin
      `https://github.com/vinpatel24/junos-terraform.git`

With above repoistories cloned and pyang installed, follow below steps to convert Junos YANG models into JSON.

1. Checkout to `upgradeFramework` branch of the above Junos Terraform repository
    ```
    cd junos-terraform
    git checkout upgradeFramework
    ```
    Git repository branch can be verified like below
     ```
     root@ubuntu:~/junos-terraform# git branch
      master
      * upgradeFramework
     ```
3. Create a folder named pyang-plugin within the junos-terraform repository
    ```
    root@ubuntu:~/junos-terraform# mkdir pyang-plugin
    ```
4. Move jtaf_json.py file into the newly created folder
    ```
    root@ubuntu:~/junos-terraform# mv jtaf_json.py pyang-plugin
    ```
5. Go to the respective Junos version directory whose YANG models need to be converted, below is an example for version 22.3
    ```
    root@ubuntu:~/junos-terraform# cd../yang/22.3/22.3R1/conf/
    ```
6. Execute the pyang command, by providing the `plugindir` option and all the Junos YANG files, using `-f`, which are to be convereted. For below case interfaces and root YANG models are considered.
   In addition, common yang models are to be provided to the plugin for successful conversion, they can provided with `-p` option. Since the data is a lot to be printed on screen hence its
   redirected to a file named `interfaces.json`.
   ```
   pyang --plugindir=/root/junos-terraform/pyang-plugin -f jtaf junos-conf-root\@2022-01-01.yang junos-conf-interfaces\@2022-01-01.yang -p ../../common/ >interface.json
   ```
   
