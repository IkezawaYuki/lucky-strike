package getproviders

import (
	"fmt"
	"github.com/apparentlymart/go-versions/versions"
	"github.com/apparentlymart/go-versions/versions/constraints"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/addrs"
	"strings"
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

type Selections map[addrs.Provider]Version

func ParseVersion(str string) (Version, error) {
	return versions.ParseVersion(str)
}

func MustParseVersion(str string) Version {
	ret, err := ParseVersion(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseVersionConstraints(str string) (VersionConstraints, error) {
	return constraints.ParseRubyStyleMulti(str)
}

func MustParseVersionConstraints(str string) VersionConstraints {
	ret, err := ParseVersionConstraints(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func MeetingConstraints(vc VersionConstraints) VersionSet {
	return versions.MeetingConstraints(vc)
}

type Platform struct {
	OS, Arch string
}

func (p Platform) String() string {
	return p.OS + "_" + p.Arch
}

func (p Platform) LessThan(other Platform) bool {
	switch {
	case p.OS != other.OS:
		return p.OS < other.OS
	default:
		return p.Arch < other.Arch
	}
}

func ParsePlatform(str string) (Platform, error) {
	underPos := strings.Index(str, "_")
	if underPos < 1 || underPos >= len(str)-2 {
		return Platform{}, fmt.Errorf("")
	}
	os, arch := str[:underPos], str[underPos+1:]
	if strings.ContainsAny(arch, " \t\n\r") {
		return Platform{}, fmt.Errorf()
	}
	if strings.ContainsAny(arch, " \t\n\r") {
		return Platform{}, fmt.Errorf("")
	}

	return Platform{
		OS:   os,
		Arch: arch,
	}, nil
}
