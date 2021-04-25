package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	keyProfile = "profile"
	// each profile
	keyWorkGroup      = "work-group"
	keyOutputLocation = "output-location"
	keyDatabaseName   = "database-name"
	keyCatalogName    = "catalog-name"
)

// find a sub for given account-id
func profileViper(v *viper.Viper, aid string) *viper.Viper {
	findSub := func(parent, key string) *viper.Viper {
		pkey := fmt.Sprintf("%s.%s", parent, key)
		v := v.Sub(pkey)
		if v == nil {
			log.Fatalf("no found profile, %v", pkey)
		}
		logger.Printf("use a profile, %v", pkey)
		return v
	}

	// use it if given explicitly.
	if p := v.GetString(keyProfile); p != "" {
		return findSub("profiles", p)
	}

	// New empty viper.
	return findSub("accounts", aid)
}

type Profile struct {
	v  *viper.Viper
	pv *viper.Viper
}

// return b if a is empty.
// return def if b is empty.
func (p *Profile) either(key, def string) string {
	a := p.v.GetString(key)
	b := p.pv.GetString(key)
	if a != "" {
		return a
	}
	if b != "" {
		return b
	}
	return def
}

func (p *Profile) WorkGroup() string {
	return p.either(keyWorkGroup, "primary")
}

func (p *Profile) CatalogName() string {
	return p.either(keyCatalogName, "AwsDataCatalog")
}

func (p *Profile) DatabaseName() string {
	return p.either(keyDatabaseName, "")
}

func (p *Profile) OutputLocation() string {
	return p.either(keyOutputLocation, "")
}
