package main

import (
    "database/sql"
    "fmt"
    "log"
    "strings"

	"os"
	"bufio"
    "sync"
	"time"

    _ "github.com/go-sql-driver/mysql"
)

var lock = &sync.Mutex{}

var hostgroupsMap map[string][]string

func main() {

	go func() {
		for
		{
			readGroupsFake();
			fmt.Println("Sleeping...")
			time.Sleep(1500 * time.Millisecond)
			fmt.Println("ok")
		}
	}()


	reader := bufio.NewReader(os.Stdin)
	for
	{
		reader.ReadString('\n')

		for key, value := range hostgroupsMap {
			fmt.Println(key, value)
		}
	}

}



func getGroups() {
	hostgroupsMapTmp := readGroups()
	
	fmt.Println("getGroups hostgroupsMapTmp: ", len(hostgroupsMapTmp))
	lock.Lock()
	defer lock.Unlock()

	hostgroupsMap = hostgroupsMapTmp;
}

func readGroupsFake() map[string][]string {
	var hostgroupsMapTmp map[string][]string = make(map[string][]string)

	hostgroupsMapTmp["fake-h.wp.lan"] = strings.Split("fake1-g", ",")

	return hostgroupsMapTmp
}

func readGroups() map[string][]string {
	fmt.Print("Mysql read...")
	
	var hostgroupsMapTmp map[string][]string = make(map[string][]string)

    db, err := sql.Open("mysql", "icinga:99QjwQeKBBMQWDnbNXxLpwF6K1helAQt@tcp(127.0.0.1:3306)/icinga")
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    res, err := db.Query("select h.alias as host,GROUP_CONCAT(g.alias) as groups from icinga_hostgroup_members m join icinga_hostgroups g on g.hostgroup_id=m.hostgroup_id join icinga_hosts h on h.host_object_id=m.host_object_id group by h.alias")

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

		hostgroupsMapTmp[host] = strings.Split(groups, ",")
    }

	fmt.Println("end")
	return hostgroupsMapTmp
}



