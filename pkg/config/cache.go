package config

type Cache struct {
	Dir   string
	Rules []*CacheRule
}

type CacheRule struct {
	Name       string
	Conditions []map[string]string
}
