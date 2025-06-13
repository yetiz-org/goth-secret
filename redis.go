package secret

type Redis struct {
	DefaultSecret
	Master RedisMeta `json:"master"`
	Slave  RedisMeta `json:"slave"`
}

type RedisMeta struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}
