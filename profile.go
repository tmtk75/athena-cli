package main

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	keyProfile = "profile"
	// each profile
	keyProfileWorkGroup      = "work-group"
	keyProfileOutputLocation = "output-location"
	keyProfileDatabaseName   = "database-name"
	keyProfileCatalogName    = "catalog-name"
)

// find a sub for given account-id
func profileViper(v *viper.Viper, aid string) *viper.Viper {
	findSub := func(parent, key string) *viper.Viper {
		pkey := fmt.Sprintf("%s.%s", parent, key)
		v := v.Sub(pkey)
		if v == nil {
			//log.Fatalf("no found profile, %v", pkey)
			return nil
		}
		logger.Printf("use a profile with the key, '%v'", pkey)
		return v
	}

	// use it if given explicitly.
	p := v.GetString(keyProfile)
	logger.Printf("given profile name: %v", p)
	if p != "" {
		v := findSub("profiles", p)
		if v != nil {
			return v
		}
		logger.Printf("profile name was given but no corresponding profile.")
	}

	logger.Printf("no given profile name.")
	logger.Printf("will find a profile correspoinding to your AWS account ID.")
	av := findSub("accounts", aid)
	if av != nil {
		logger.Printf("found a profile correspoinding to your AWS account ID.")
		return av
	}

	logger.Printf("no profile. returns default.")
	return viper.GetViper()
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
	return p.either(keyProfileWorkGroup, "primary")
}

func (p *Profile) CatalogName() string {
	return p.either(keyProfileCatalogName, "AwsDataCatalog")
}

func (p *Profile) DatabaseName() string {
	return p.either(keyProfileDatabaseName, "")
}

func (p *Profile) OutputLocation() string {
	return p.either(keyProfileOutputLocation, "")
}
