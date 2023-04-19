package main

import (
	"fmt"
	"flag"
	//"log"
	"os"
	//"sort"
	"strings"

	//"bufio"
	//"encoding/json"
	
	//"config/config"
	"wp/icinga2eventsforwarder/config"
	"wp/icinga2eventsforwarder/httpclient"
	"wp/icinga2eventsforwarder/nats"

	"wp/icinga2eventsforwarder/mariadb/icingasql"
	"wp/icinga2eventsforwarder/mariadb/directorsql"
	//"github.com/urfave/cli/v2"
	//"github.com/influxdata/telegraf/config"
)


//import "github.com/nats-io/nats.go"

// func runApp(args []string, c I2EFConfig, m App) error {
	
// 	return m.Run(args[0])
// }

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

	agent := I2EF{ }
	agent.loadConfiguration(configFile)

    _ = icingasql.InitSQLdata(agent.conf.Icinga)
    directorsql.InitSQLdata(agent.conf.Agent.Host, agent.conf.Director)

    var natsClient = nats.GetNATSclient()
	natsClient.Connect(agent.conf.Nats)

	fmt.Println("natsClient.Connected")
	agent.runReader()

	// err := runApp(os.Args, c, &agent)
	// if err != nil {
	// 	log.Fatalf("E! %s", err)
	// }
}

type I2EF struct {
	pprofErr <-chan error

	conf *config.Config
}

// func (t *I2EF) Init(pprofErr <-chan error) {
// 	t.pprofErr = pprofErr


// }

func (t *I2EF) loadConfiguration(configFilePath string) (error) {
	// If no other options are specified, load the config file and run.
	t.conf = config.LoadConfiguration(configFilePath)
	if t.conf.Agent == nil {
		t.conf.Agent = &config.AGENT { Host: "" }
	} 
	if t.conf.Agent.Host == "" {

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		t.conf.Agent.Host = strings.ToLower(hostname)
	}
    fmt.Println("config host", t.conf.Agent.Host)

	return nil
}

func (t *I2EF) runReader() {
	var stream httpclient.HTTPstream
	stream.Init(t.conf.EventsStream)

	stream.Run(t.conf.Agent.Host)
	fmt.Println("run Reader ended")
}
// func (t *I2EF) runAgent(ctx context.Context, c *config.Config) error {
// 	var err error

// 	ag := agent.NewAgent(c)

// 	// // Notify systemd that telegraf is ready
// 	// // SdNotify() only tries to notify if the NOTIFY_SOCKET environment is set, so it's safe to call when systemd isn't present.
// 	// // Ignore the return values here because they're not valid for platforms that don't use systemd.
// 	// // For platforms that use systemd, telegraf doesn't log if the notification failed.
// 	// _, _ = daemon.SdNotify(false, daemon.SdNotifyReady)



// 	return ag.Run(ctx)
// }


















type App interface {
	Init(<-chan error)
	Run(configFile string) error
}
