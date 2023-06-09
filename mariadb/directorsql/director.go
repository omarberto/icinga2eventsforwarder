package directorsql

import (
    "encoding/json"
    "database/sql"
    "fmt"
    "log"

	//_ "os"
	//"bufio"
    "strings"
	"time"

	"wp/icinga2eventsforwarder/config"
	"wp/icinga2eventsforwarder/tornado"
	"wp/icinga2eventsforwarder/nats"
    _ "github.com/go-sql-driver/mysql"
)

var connectionString string
var firstcall bool = true
var refreshInterval time.Duration = 60 * time.Second
var lastDeleteHostTime time.Time = time.Now().UTC()
//time.Date(2023, 04, 14, 11, 30, 32, 0, time.UTC)
//time.Now().UTC()
var lastDeletedServiceTime time.Time = time.Now().UTC()
//time.Date(2022, 04, 14, 11, 30, 32, 0, time.UTC)
//time.Now().UTC()
var lastChangedHostTime time.Time = time.Now().UTC()
var lastChangedServiceTime time.Time = time.Now().UTC()
func InitSQLdata(hostname string, c *config.MYSQL) {
	if firstcall {
		firstcall = false
		connectionString = c.Username + ":" + c.Password + c.Connection

		if c.RefreshInterval == "" {
			//c.RefreshInterval = 60 * time.Second
		} else {
			t, err := time.ParseDuration(c.RefreshInterval)
			if err == nil {
                refreshInterval = t	
			} else{
				//c.RefreshInterval = 60 * time.Second
				//return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tmp.Timeout, err)
			}
		}

		go func() {

			for
			{
				readDeletedHosts(hostname)
				readDeletedServices(hostname)
				readChangedHosts(hostname)
				readChangedServices(hostname)
				time.Sleep(refreshInterval)
			}
		}()
	}
}

func readDeletedHosts(hostname string) {
	log.Println("readDeletedHosts start")
	
    db, err := sql.Open("mysql", connectionString)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query("select object_name,convert_tz(change_time,@@session.time_zone,'+00:00') as change_time_utc from director_activity_log where action_name = 'delete' and object_type = 'icinga_host' and convert_tz(change_time,@@session.time_zone,'+00:00') > '" + sqlFormat(lastDeleteHostTime) + "'")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    for res.Next() {

		//var object_type string
		var object_name string
		//var old_properties string
		var change_time_utc time.Time
        err := res.Scan(&object_name, &change_time_utc)

        if err != nil {
            log.Fatal(err)
			continue
		}
		

		// fmt.Println("object_name", object_name)
		// fmt.Println("change_time_utc", sqlFormat(change_time_utc))
		// fmt.Println("object_name is not null", object_name != "")
		// fmt.Println("\n")


		if(change_time_utc.After(lastDeleteHostTime)){
			lastDeleteHostTime = change_time_utc
		}

		if(object_name != "") {
			var tornadoMessageObj tornado.TORNADOMSG_OBJ_DELETED = tornado.TORNADOMSG_OBJ_DELETED{
				Host: hostname,
				Type: "icinga2eventsformatter.deletedobject",
				Obj: tornado.OBJ_DELETED{
					Host: object_name,
					Type: "icinga_host",
					Service: "",
				},
			}

			tornadoMessageBin, err := json.Marshal(tornadoMessageObj)

			if err == nil {
				log.Println("Writing tornado message: ", string(tornadoMessageBin))
				var natsClient = nats.GetNATSclient()
				natsClient.Write(tornadoMessageBin)
			}
		}
    }

	log.Println("readDeletedHosts end")
}
func sqlFormat(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second())
}
func readDeletedServices(hostname string) {
	log.Println("readDeletedServices start")
	
    db, err := sql.Open("mysql", connectionString)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query("select object_name,old_properties,convert_tz(change_time,@@session.time_zone,'+00:00') as change_time_utc from director_activity_log where action_name = 'delete' and object_type = 'icinga_service' and convert_tz(change_time,@@session.time_zone,'+00:00') > '" + sqlFormat(lastDeletedServiceTime) + "'")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    for res.Next() {

		var object_name string
		var old_properties string
		var change_time_utc time.Time
        err := res.Scan(&object_name, &old_properties, &change_time_utc)

        if err != nil {
            log.Println(err)
			continue
		}
		
		var oldPropertiesObj OLD_PROPERTIES_SERVICE

		err = json.Unmarshal([]byte(old_properties), &oldPropertiesObj)

		if err != nil {
            log.Println(err)
			continue
		}

		// fmt.Println("object_name", object_name)
		// fmt.Println("old_properties", old_properties)
		// fmt.Println("change_time_utc", sqlFormat(change_time_utc))
		// fmt.Println("object_name is not null", object_name != "")
		// fmt.Println("host is not null", oldPropertiesObj.Host != "")
		// fmt.Println("object_type", strings.EqualFold(oldPropertiesObj.Type, "object"))
		// fmt.Println("\n")


		if(change_time_utc.After(lastDeletedServiceTime)){
			lastDeletedServiceTime = change_time_utc
		}

		if(strings.EqualFold(oldPropertiesObj.Type, "object") && oldPropertiesObj.Host != "") {
			
			var tornadoMessageObj tornado.TORNADOMSG_OBJ_DELETED

			if(oldPropertiesObj.Host != "") {
				tornadoMessageObj = tornado.TORNADOMSG_OBJ_DELETED{
					Host: hostname,
					Type: "icinga2eventsformatter.deletedobject",
					Obj: tornado.OBJ_DELETED{
						Host: oldPropertiesObj.Host,
						Type: "icinga_service",
						Service: object_name,
					},
				}
			} else {
				tornadoMessageObj = tornado.TORNADOMSG_OBJ_DELETED{
					Host: hostname,
					Type: "icinga2eventsformatter.deletedobject",
					Obj: tornado.OBJ_DELETED{
						Host: "",
						Type: "icinga_service2",
						Service: object_name,
					},
				}
			}

			tornadoMessageBin, err := json.Marshal(tornadoMessageObj)

			if err == nil {
				log.Println("Writing tornado message: ", string(tornadoMessageBin))
				var natsClient = nats.GetNATSclient()
				natsClient.Write(tornadoMessageBin)
			}
		}
///////////////////////////////////////////
    }

	log.Println("readDeletedServices end")
}

func readChangedHosts(hostname string) {
	log.Println("readChangedHosts start")

    db, err := sql.Open("mysql", connectionString)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query("select object_name,old_properties,new_properties,convert_tz(change_time,@@session.time_zone,'+00:00') as change_time_utc from director_activity_log where action_name = 'modify' and object_type = 'icinga_host' and convert_tz(change_time,@@session.time_zone,'+00:00') > '" + sqlFormat(lastChangedHostTime) + "'")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    for res.Next() {

		var object_name string
		var old_properties string
		var new_properties string
		var change_time_utc time.Time
        err := res.Scan(&object_name, &old_properties, &new_properties, &change_time_utc)

        if err != nil {
            log.Println(err)
			continue
		}
		
		var oldPropertiesObj PROPERTIES_HOST
		var newPropertiesObj PROPERTIES_HOST

		err = json.Unmarshal([]byte(old_properties), &oldPropertiesObj)

		if err != nil {
            log.Println(err)
			continue
		}
		
		err = json.Unmarshal([]byte(new_properties), &newPropertiesObj)

		if err != nil {
            log.Println(err)
			continue
		}

		if(change_time_utc.After(lastChangedHostTime)){
			lastChangedHostTime = change_time_utc
		}

		if(object_name != "" && (newPropertiesObj.Name != oldPropertiesObj.Name || newPropertiesObj.Address != oldPropertiesObj.Address)) {
			var tornadoMessageObj tornado.TORNADOMSG_OBJ_DELETED = tornado.TORNADOMSG_OBJ_DELETED{
				Host: hostname,
				Type: "icinga2eventsformatter.modifiedobject",
				Obj: tornado.OBJ_DELETED{
					Host: oldPropertiesObj.Name,
					Type: "icinga_host",
					Service: "",
				},
			}

			tornadoMessageBin, err := json.Marshal(tornadoMessageObj)

			if err == nil {
				log.Println("Writing tornado message: ", string(tornadoMessageBin))
				var natsClient = nats.GetNATSclient()
				natsClient.Write(tornadoMessageBin)
			}
		}
///////////////////////////////////////////
    }

	log.Println("readChangedHosts end")
}
func readChangedServices(hostname string) {
	log.Println("readChangedServices start")
	
    db, err := sql.Open("mysql", connectionString)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query("select object_name,old_properties,new_properties,convert_tz(change_time,@@session.time_zone,'+00:00') as change_time_utc from director_activity_log where action_name = 'modify' and object_type = 'icinga_service' and convert_tz(change_time,@@session.time_zone,'+00:00') > '" + sqlFormat(lastChangedServiceTime) + "'")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    for res.Next() {

		var object_name string
		var old_properties string
		var new_properties string
		var change_time_utc time.Time
        err := res.Scan(&object_name, &old_properties, &new_properties, &change_time_utc)

        if err != nil {
            log.Println(err)
			continue
		}
		
		var oldPropertiesObj PROPERTIES_SERVICE
		var newPropertiesObj PROPERTIES_SERVICE

		err = json.Unmarshal([]byte(old_properties), &oldPropertiesObj)

		if err != nil {
            log.Println(err)
			continue
		}
		
		err = json.Unmarshal([]byte(new_properties), &newPropertiesObj)

		if err != nil {
            log.Println(err)
			continue
		}

		if(change_time_utc.After(lastChangedServiceTime)){
			lastChangedServiceTime = change_time_utc
		}

		if(strings.EqualFold(oldPropertiesObj.Type, "object") && (newPropertiesObj.Name != oldPropertiesObj.Name)) {
			var tornadoMessageObj tornado.TORNADOMSG_OBJ_DELETED = tornado.TORNADOMSG_OBJ_DELETED{
				Host: hostname,
				Type: "icinga2eventsformatter.modifiedobject",
				Obj: tornado.OBJ_DELETED{
					Host: oldPropertiesObj.Host,
					Type: "icinga_service",
					Service: oldPropertiesObj.Name,
				},
			}

			tornadoMessageBin, err := json.Marshal(tornadoMessageObj)

			if err == nil {
				log.Println("Writing tornado message: ", string(tornadoMessageBin))
				var natsClient = nats.GetNATSclient()
				natsClient.Write(tornadoMessageBin)
			}
		}
///////////////////////////////////////////
    }

	log.Println("readChangedHosts end")
}

type OLD_PROPERTIES_SERVICE struct {
	Host string  `json:"host"`
	Type string  `json:"object_type"`
}

type PROPERTIES_SERVICE struct {
	Type string  `json:"object_type"`
	Name    string `json:"object_name"`
	Host string  `json:"host"`
}
type PROPERTIES_HOST struct {
	Type    string `json:"object_type"`
	Name    string `json:"object_name"`
	Address string `json:"address"`
}