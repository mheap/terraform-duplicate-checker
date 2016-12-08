package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
)

func main() {
	searches := map[string][]string{}

	searches["azurerm_network_interface"] = []string{"name"}
	searches["azurerm_network_security_group"] = []string{"name"}
	searches["azurerm_public_ip"] = []string{"name"}
	searches["azurerm_sql_database"] = []string{"name"}
	searches["azurerm_sql_server"] = []string{"name"}
	searches["azurerm_storage_account"] = []string{"name"}
	searches["azurerm_storage_container"] = []string{"name"}
	searches["azurerm_subnet"] = []string{"name", "address_prefix"}
	searches["azurerm_virtual_machine"] = []string{"name"}
	searches["azurerm_virtual_network"] = []string{"name"}

	run(searches)
}

func run(searches map[string][]string) {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("unable to read from stdin: %s", err)
		return
	}
	var v map[string]interface{}
	err = hcl.Unmarshal(input, &v)
	instances := mushData(v["resource"])

	//tJ(instances)
	//return

	for k, v := range searches {
		outputDuplicates(instances, k, v)
	}
}

func tJ(f interface{}) {
	str, err := json.Marshal(f)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	fmt.Println(string(str))
}

func mushData(v interface{}) map[string]map[string]interface{} {
	things := make(map[string]map[string]interface{})
	for _, resources := range v.([]map[string]interface{}) {
		for resource, instances := range resources {
			if things[resource] == nil {
				things[resource] = make(map[string]interface{})
			}
			for _, instance := range instances.([]map[string]interface{}) {
				for name, fields := range instance {
					things[resource][name] = fields.([]map[string]interface{})[0]
				}
			}
		}
	}
	return things
}

func findDuplicates(instances map[string]interface{}, uniqueFields []string) map[string][]string {
	duplicates := make(map[string][]string)

	for _, fieldName := range uniqueFields {
		currentValues := make(map[string][]string)
		var fieldValue string
		for instanceName, v := range instances {
			fieldValue = v.(map[string]interface{})[fieldName].(string)
			currentValues[fieldValue] = append(currentValues[fieldValue], instanceName)
		}
		if len(currentValues) != len(instances) {
			duplicates[fieldName] = currentValues[fieldValue]
		}
	}
	return duplicates
}

func outputDuplicates(instances map[string]map[string]interface{}, resource string, fields []string) {
	duplicates := findDuplicates(instances[resource], fields)

	if len(duplicates) > 0 {
		fmt.Println("")
		fmt.Println("#########################")
		fmt.Printf("## %s\n", resource)
		fmt.Println("#########################")
		for f, n := range duplicates {
			fmt.Printf(" - %s :: %s\n", f, n)
		}
	}
}
