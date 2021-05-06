package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cisco/senml"
	"github.com/ifraiot/machine-simulation/mqtt"
)

type Machine struct {
	MachineName string `json:"name"`
	Config      struct {
		OrganizationID string `json:"organizationId"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	} `json:"config"`
	Condition       int `json:"condition"`
	OutputQTY       int `json:"outputQty"`
	RejectOutputQTY int `json:"rejectOutputQty"`
}

func main() {
	machines, err := loadMcConfig("machines.js")
	if err != nil {
		fmt.Println(err)
	}

	for _, mc := range machines {

		m := mqtt.NewMQTT(
			"staging.mqtt.ifra.io",
			"1883",
			mc.Config.Username,
			mc.Config.Password,
		)

		for {

			condition := 1.0
			outputQty := 10.0
			rejectedOutputQty := 0.0

			s := senml.SenML{
				Records: []senml.SenMLRecord{
					senml.SenMLRecord{Value: &condition, Unit: "-", Name: "condition"},
					senml.SenMLRecord{Value: &outputQty, Unit: "-", Name: "output_qty"},
					senml.SenMLRecord{Value: &rejectedOutputQty, Unit: "-", Name: "reject_output_qty"},
				},
			}

			dataOut, err := senml.Encode(s, senml.JSON, senml.OutputOptions{})
			if err != nil {
				fmt.Println("Encode of SenML failed")
			} else {
				fmt.Println(string(dataOut))
			}

			messages := string(dataOut)

			m.Pub(mc.Config.OrganizationID, messages)
			time.Sleep(5 * time.Second)
		}

	}

}

func loadMcConfig(fileName string) ([]Machine, error) {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var machines []Machine

	err = json.Unmarshal(byteValue, &machines)
	if err != nil {
		fmt.Println(err)
	}

	return machines, nil
}
