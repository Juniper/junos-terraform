resource "junos-vsrx" "config_vsrx" {
  resource_name = "example_resource"
  engineer = [
    {
      name = "Test Engineer"
      age = 30
      commute = "car"
      food = [
        {
            office = { 
                coffee = [
                    {
                    size = "tall"
                    milk = "true"
                    }
                ]
                donut = true
            }
        },
        {
            home = {
                coffee = ""
                salad = ""
            }
        }
      ]
    }
  ]
}
