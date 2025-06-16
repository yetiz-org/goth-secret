package secret

type Cassandra struct {
	DefaultSecret
	Writer CassandraMeta `json:"writer"`
	Reader CassandraMeta `json:"reader"`
}

type CassandraMeta struct {
	Endpoints []string `json:"endpoints"`
	Keyspace  string   `json:"keyspace"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	CaPath    string   `json:"ca_path"`
}
