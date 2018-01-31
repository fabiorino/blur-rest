package config

import "blur-rest/store"

// GlobalConfig is a variable accessible from everywhere in the program
var GlobalConfig struct {
	Fqdn  string
	Port  string
	Store store.Store
}
