package getproviders

import (
	"github.com/apparentlymart/go-versions/versions"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/addrs"
)

type Version = versions.Version

var UnspecifiedVersion Version = versions.Unspecified

type VersionList = versions.List

type Warnings = []string

type Requirements map[addrs.Provider]version.Constraints

func (r Requirements) Merge(other Requirements) Requirements {
	ret := make(Requirements)
	for addr, constraints := range r {
		ret[addr] = constraints
	}
	for addr, constraints := range other {
		ret[addr] = append(ret[addr], constraints...)
	}
	return ret
}
