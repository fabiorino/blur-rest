package config

import "blur-rest/store"

var GlobalConfig struct {
	Fqdn  string
	Port  string
	Store store.Store
}
