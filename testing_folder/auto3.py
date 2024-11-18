from jinja2 import FileSystemLoader, environment

parent = "interfaces"
data = [
	{
		"interfaces" : {
			"attributes" : [
				{
					"attr_name" : "interface",
					"attr_type" : "list"
				}
			]
		}
	},
    {
    	"interface" : {
    		"type": "object",
    		"attributes" : [
    			{
    				"attr_name": "name",
    				"attr_type": "string"
    			},
    			{
    				"attr_name": "description",
    				"attr_type": "string"
    			},
    			{
    				"attr_name": "mtu",
    				"attr_type": "int64"
    			},
                {
                    "attr_name" : "vlan_tagging",
                    "attr_type" : "bool"
                },
				{
					"attr_name" : "unit",
					"attr_type" : "list"
				}
    		],
    	},
    },
	{
		"unit" : {
			"type": "list",
			"attributes" : [
				{
					"attr_name": "name",
					"attr_type" : "string"
				},
				{
					"attr_name" : "description",
					"attr_type" : "string"
				},
				{
					"attr_name" : "vlan_id",
					"attr_type" : "int32"
				}
			]
		}
	}
]

file_loader = FileSystemLoader('./')
env = environment.Environment(loader=file_loader)

tmp1 = env.get_template('XML_struct.j2')
xml_struct = tmp1.render(data=data)
tmp2 = env.get_template('Terraform_file.j2')
terraform_file = tmp2.render(data=data)
tmp3 = env.get_template('schema.j2')
schema = tmp3.render(data=data)
tmp4 = env.get_template('create.j2')
create = tmp4.render(data=data)
tmp5 = env.get_template('read.j2')
read = tmp5.render(data=data)
tmp6 = env.get_template('update_delete.j2')
update_delete = tmp6.render(data=data)
with open('resource_Interfaces.go', 'w') as file:
	file.write(xml_struct)
	file.write(terraform_file)
	file.write(schema)
	file.write(create)
	file.write(read)
	file.write(update_delete)



# file_loader = FileSystemLoader('./')
# env = environment.Environment(loader=file_loader)
# # tmp = env.get_template('XML_struct.j2')
# # tmp = env.get_template('Terraform_file.j2')
# # tmp = env.get_template('schema.j2')
# # tmp = env.get_template('create.j2')
# tmp = env.get_template('read.j2')
# # tmp = env.get_template('update_delete.j2')
# print(tmp.render(parent = parent, data=data))
