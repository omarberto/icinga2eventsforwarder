//applicativo che manda un check_rest di test a cadenza regolare, per vedere che succede lato neteye
package main

import (
	"fmt"
	"flag"
	"bufio"
	"os"
	//"sort"
	//"strings"


    "encoding/json"

	"wp/icinga2eventsforwarder/config"
	//"wp/icinga2eventsforwarder/nats"
)



func main() {

	var configFileArg bool
	var configFile string
	flag.BoolVar(&configFileArg, "config", false, "configuration file path")
    flag.Parse()

    if !configFileArg {
		fmt.Println("No configuration file")
		return
	}

	args := os.Args
    for index, arg := range args {
		if arg == "-config" || arg == "--config" {
			if index + 1 == len(args) { 
				fmt.Println("No configuration file")
				return
			}else{
				configFile = args[index + 1]
			}
		}
	}
	
	var conf = config.LoadConfiguration(configFile)
	fmt.Println("Struct is:", conf.Nats)

    //var natsClient = nats.GetNATSclient()
	//natsClient.Connect(conf.Nats)

	reader := bufio.NewReader(os.Stdin)

	var checkResult = `{
		"check_result": {
		  "active": false,
		  "check_source": "beomM1.wp.lan",
		  "command": null,
		  "execution_end": 1677248383.333,
		  "execution_start": 1677248383.333,
		  "exit_status": 0,
		  "output": "SRS Report rendering is OK",
		  "performance_data": [
			"'rendering SRSReportServerWarmup.AutoDesign'=0.607s;1;5",
			"'rendering SysUserLicenseCountReport.Report'=0.697s;1;5"
		  ],
		  "schedule_end": 1677248383.334313,
		  "schedule_start": 1677248383.334313,
		  "state": 0,
		  "ttl": 0,
		  "vars_after": {
			"attempt": 1,
			"reachable": true,
			"state": 0,
			"state_type": 1
		  },
		  "vars_before": {
			"attempt": 1,
			"reachable": true,
			"state": 0,
			"state_type": 1
		  },
		  "type": "Check_Result"
		},
		"host": "replicatest.wp.lan",
		"groups": ["fake1-test"],
		"timestamp": 1677248383.334386
	  }`

	var checkResultObj CHECKRESULT

	err := json.Unmarshal([]byte(checkResult), &checkResultObj)

	if err != nil {
		fmt.Println(err)
	}

	tornadoMessage, err := json.Marshal(checkResultObj)

    if err != nil {
        fmt.Println(err)
    }
	
	fmt.Println("Struct is:", string(tornadoMessage))

	//natsClient.Write([]byte(tornadoMessage))

	reader.ReadString('\n')
}

type CHECKRESULT struct {
    Host          string `json:"host"`
    Service          string `json:"service"`
	Groups       []string `json:"groups"`
	CheckResult  CHECKRESULTOUT `json:"check_result"`
}
type CHECKRESULTOUT struct {
    ExitStatus          int `json:"exit_status"`
    PluginOutput          string `json:"output"`
	PerformanceData       []string `json:"performance_data"`
}