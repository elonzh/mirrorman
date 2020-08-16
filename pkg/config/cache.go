package config

type Cache struct {
	Backend string
	Options map[string]interface{}
	Rules   []*CacheRule
}

type CacheRule struct {
	Name       string
	Conditions []map[string]string
}
