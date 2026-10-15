package config

//import "time"

type ICINGADBQUERIES struct {
	ReadGroupsQuery string        `toml:"readgroupsquery"`
}

type ICINGAMYSQL struct {
	Connection string        `toml:"connection"`
	Username string        `toml:"username"`
	Password string        `toml:"password"`
	RefreshInterval string        `toml:"refreshinterval"`
	Queries ICINGADBQUERIES        `toml:"queries"`
}

type DIRECTORDBQUERIES struct {
	 DeletedHostsQuery string        `toml:"deletedhostsquery"`
	 ChangedHostsQuery string        `toml:"changedhostsquery"`
	 DeletedServicesQuery string        `toml:"deletedservicesquery"`
	 ChangedServicesQuery string        `toml:"changedservicesquery"`
}

type DIRECTORMYSQL struct {
	Connection string        `toml:"connection"`
	Username string        `toml:"username"`
	Password string        `toml:"password"`
	RefreshInterval string        `toml:"refreshinterval"`
	Queries DIRECTORDBQUERIES        `toml:"queries"`
}
