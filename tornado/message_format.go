package tornado

import (
    "encoding/json"
	"log"
	//"fmt"
	//"os"
	
	//"wp/icinga2eventsforwarder/nats"
	"wp/icinga2eventsforwarder/mariadb/icingasql"
)

type Message struct {
	formatString1 string
	formatString2 string
	formatString3 string
}

func NewMessage() *Message{
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
	
	var checkResultObj CHECKRESULT

	err := json.Unmarshal([]byte(icinga2json), &checkResultObj)

	if err != nil {
		log.Println("error in unmarhal", icinga2json)
		log.Println(err)
		return nil, err
	}

	if checkResultObj.Service == "" {
		//return nil, errors.New("no service")
		checkResultObj.Type = "host"
		if checkResultObj.CheckResult.ExitStatus > 1 {
			checkResultObj.CheckResult.ExitStatus = 1
		} else {
			checkResultObj.CheckResult.ExitStatus = 0
		}
	} else {
		checkResultObj.Type = "service"
	}

    var hg map[string][]string = icingasql.GetSQLdata()
	checkResultObj.Groups = hg[checkResultObj.Host]

	var tornadoMessageObj TORNADOMSG = TORNADOMSG{
		Host: hostname,
		Type: "icinga2eventsformatter.checkresult",
		Event: checkResultObj,
	}

	tornadoMessage, err := json.Marshal(tornadoMessageObj)

    if err != nil {
		log.Println(err)
		return nil, err
    }
	
	return tornadoMessage, nil
}

type TORNADOMSG struct {
	Host string  `json:"host"`
	Type string  `json:"type"`
	Event 	CHECKRESULT `json:"icinga2event"`
}

type CHECKRESULT struct {
    Host          string `json:"host"`
    Service          string `json:"service"`
	Type          string `json:"type"`
	Groups       []string `json:"groups"`
	CheckResult  CHECKRESULTOUT `json:"check_result"`
}
type CHECKRESULTOUT struct {
    ExitStatus          float64 `json:"exit_status"`
    State               float64 `json:"state"`
    PluginOutput          string `json:"output"`
	PerformanceData       []string `json:"performance_data"`
}

//poi andrebbe scritto meglio, per ora faccio tutto brutale..
type TORNADOMSG_OBJ_DELETED struct {
	Host string  `json:"host"`
	Type string  `json:"type"`
	Obj	 OBJ_DELETED `json:"obj"`
}
type OBJ_DELETED struct {
    Host          string `json:"host"`
    Service          string `json:"service"`
	Type          string `json:"type"`
}


