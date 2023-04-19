package icingasql

import (
    "database/sql"
    "fmt"
    "log"
    "strings"

	//_ "os"
	//"bufio"
    "sync"
	"time"

	"wp/icinga2eventsforwarder/config"
    _ "github.com/go-sql-driver/mysql"
)
var lock = &sync.Mutex{}

var connectionString string
var firstcall bool = true
var refreshInterval time.Duration = 60 * time.Second
var Hostgroups map[string][]string = make(map[string][]string)
func InitSQLdata(c *config.MYSQL) map[string][]string {
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
				getGroups();
				fmt.Println("Sleeping...")
				time.Sleep(refreshInterval)
				fmt.Println("ok")
			}
		}()
	}
    return Hostgroups
}

func GetSQLdata() map[string][]string {
    return Hostgroups
}

func getGroups() {
	HostgroupsTmp := readGroups()
	//HostgroupsTmp := readGroupsFake()
	
	fmt.Println("getGroups HostgroupsTmp: ", len(HostgroupsTmp))
	lock.Lock()
	defer lock.Unlock()

	Hostgroups = HostgroupsTmp;
}

func readGroupsFake() map[string][]string {
	var HostgroupsTmp map[string][]string = make(map[string][]string)

	HostgroupsTmp["fake-h.wp.lan"] = strings.Split("fake1-g", ",")

	return HostgroupsTmp
}

func readGroups() map[string][]string {
	fmt.Print("Mysql read...")
	
	var HostgroupsTmp map[string][]string = make(map[string][]string)

    db, err := sql.Open("mysql", connectionString)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

//    res, err := db.Query("select h.alias as host,GROUP_CONCAT(g.alias) as groups from icinga_hostgroup_members m join icinga_hostgroups g on g.hostgroup_id=m.hostgroup_id join icinga_hosts h on h.host_object_id=m.host_object_id group by h.alias")
    res, err := db.Query("select oh.name1 as host,GROUP_CONCAT(og.name1) as groups from icinga_hostgroup_members m join icinga_hostgroups g on g.hostgroup_id=m.hostgroup_id join icinga_objects og on og.object_id=g.hostgroup_object_id join icinga_hosts h on h.host_object_id=m.host_object_id join icinga_objects oh on oh.object_id=h.host_object_id where og.is_active=1 and oh.is_active=1 group by oh.name1")

    defer res.Close()

    if err != nil {
        log.Fatal(err)
    }

    for res.Next() {

		var host string
		var groups string
        err := res.Scan(&host, &groups)

        if err != nil {
            log.Fatal(err)
			continue
		}

		HostgroupsTmp[host] = strings.Split(groups, ",")
    }

	fmt.Println("end")
	return HostgroupsTmp
}
