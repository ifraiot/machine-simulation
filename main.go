package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/cisco/senml"
	"github.com/ifraiot/machine-simulation/mqtt"
)

func init() {
	rand.Seed(time.Now().Unix())
}

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
			var condition = 1.0

			outputQty := float64(randInt(5, 20))
			rejectedOutputQty := float64(randInt(0, 7))

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

// randInt randoms integer between min and max inclusively using "math/rand" package.
func randInt(min, max int) int {
	return rand.Intn((max - min) + 1) + min
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
