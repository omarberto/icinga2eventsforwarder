package tornado

import (
	"encoding/json"
	"log"
	"strings"
	"wp/icinga2eventsforwarder/mariadb/icingasql"
)

type Message struct {
	formatString1 string
	formatString2 string
	formatString3 string
}

func NewMessage() *Message {
	return &Message{
		formatString1: "{\"host\":\"",
		formatString2: "\",\"type\":\"icinga2eventsformatter.checkresult\",\"icinga2event\":",
		formatString3: "}",
	}
}

func (mess *Message) Append(hostname string, icinga2json string) string {
	return mess.formatString1 + hostname + mess.formatString2 + icinga2json + mess.formatString3
}

func (mess *Message) Parse(hostname string, icinga2json string) ([]byte, error) {

	var eventDataMap map[string]interface{}

	err := json.Unmarshal([]byte(icinga2json), &eventDataMap)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var Type string = eventDataMap["type"].(string)

	var tornadoMessageObj TORNADOMSG = TORNADOMSG{
		Host: hostname,
		Type: "icinga2eventsformatter." + strings.ToLower(Type),
	}

	switch Type {
	case "CheckResult":
		{
			var checkResultMap map[string]interface{} = make(map[string]interface{})
			checkResultMap["host"] = eventDataMap["host"].(string)
			checkResultMap["service"] = eventDataMap["service"]

			var checkResultObj map[string]interface{} = eventDataMap["check_result"].(map[string]interface{})

			// Check if the service field is empty
			if eventDataMap["service"] == nil || eventDataMap["service"] == "" {
				checkResultMap["type"] = "host"
				if checkResultObj["exit_status"].(float64) > 1 {
					checkResultObj["exit_status"] = 1
				} else {
					checkResultObj["exit_status"] = 0
				}
			} else {
				checkResultMap["type"] = "service"
			}

			checkResultMap["check_result"] = map[string]interface{}{
				"exit_status":      checkResultObj["exit_status"],
				"state":            checkResultObj["state"],
				"output":           checkResultObj["output"],
				"performance_data": checkResultObj["performance_data"],
			}

			var hg map[string][]string = icingasql.GetSQLdata()
			checkResultMap["groups"] = hg[eventDataMap["host"].(string)]

			// Add the check_result to the tornadoMessageObj
			tornadoMessageObj.Event = checkResultMap
		}
	default:
		{
			tornadoMessageObj.Event = eventDataMap
		}
	}

	tornadoMessage, err := json.Marshal(tornadoMessageObj)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return tornadoMessage, nil
}

type TORNADOMSG struct {
	Host string `json:"host"`
	Type string `json:"type"`
	//Event CHECKRESULT `json:"icinga2event"`
	Event map[string]interface{} `json:"icinga2event"`
}

// type CHECKRESULT struct {
// 	Host        string         `json:"host"`
// 	Service     string         `json:"service"`
// 	Type        string         `json:"type"`
// 	Groups      []string       `json:"groups"`
// 	CheckResult CHECKRESULTOUT `json:"check_result"`
// }
// type CHECKRESULTOUT struct {
// 	ExitStatus      float64     `json:"exit_status"`
// 	State           float64     `json:"state"`
// 	PluginOutput    string      `json:"output"`
// 	PerformanceData interface{} `json:"performance_data"`
// }

// poi andrebbe scritto meglio, per ora faccio tutto brutale..
type TORNADOMSG_OBJ_DELETED struct {
	Host string      `json:"host"`
	Type string      `json:"type"`
	Obj  OBJ_DELETED `json:"obj"`
}
type OBJ_DELETED struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Type    string `json:"type"`
}
