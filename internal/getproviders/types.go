package getproviders

import (
	"fmt"
	"github.com/apparentlymart/go-versions/versions"
	"github.com/apparentlymart/go-versions/versions/constraints"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/addrs"
	"runtime"
	"sort"
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

var CurrentPlatform = Platform{
	OS:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

type PackageMeta struct {
	Provider         addrs.Provider
	Version          Version
	ProtocolVersions VersionList
	TargetPlatform   Platform
	Filename         string
	Location         PackageLocation
	Authentication   PackageAuthentication
}

func (m PackageMeta) LessThan(other PackageMeta) bool {
	switch {
	case m.Provider != other.Provider:
		return m.Provider.LessThan(other.Provider)
	case m.TargetPlatform != other.TargetPlatform:
		return m.TargetPlatform.LessThan(other.TargetPlatform)
	case m.Version != other.Version:
		return m.Version.LessThan(other.Version)
	default:
		return false
	}
}

func (m PackageMeta) UnpackedDirectoryPath(baseDir string) string {
	return UnpackedDirectoryPathForPackage(baseDir, m.Provider, m.Version, m.TargetPlatform)
}

func (m PackageMeta) PackedFilePath(baseDir string) string {
	return PackedFilePathForPatckage(baseDir, m.Provider, m.Version, m.TargetPlatform)
}

func (m PackageMeta) AcceptableHashes() []Hash {
	auth, ok := m.Authentication.(PackageAuthenticationHashes)
	if !ok {
		return nil
	}
	return auth.AccepableHashes()
}

type PackageLocation interface {
	packageLocation()
	String() string
}

type PackageLocalArchive string

func (p PackageLocalArchive) packageLocation() {}

func (p PackageLocalArchive) String() string {
	return string(p)
}

type PackageLocalDir string

func (p PackageLocalDir) packageLocation() {}

func (p PackageLocalDir) String() string { return string(p) }

type PackageHTTPURL string

func (p PackageHTTPURL) packageLocation() {}

func (p PackageHTTPURL) String() string { return string(p) }

type PackageMetaList []PackageMeta

func (l PackageMetaList) Len() int {
	return len(l)
}

func (l PackageMetaList) Less(i, j int) bool {
	return l[i].LessThan(l[j])
}

func (l PackageMetaList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l PackageMetaList) Sort() {
	sort.Stable(l)
}

func (l PackageMetaList) FilterPlatform(target Platform) PackageMetaList {
	var ret PackageMetaList
	for _, m := range l {
		if m.TargetPlatform == target {
			ret = append(ret, m)
		}
	}
	return ret
}
