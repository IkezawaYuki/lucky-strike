package getproviders

import (
	"context"
	"github.com/hashicorp/terraform/addrs"
)

type Source interface {
	AvailableVersions(ctx context.Context, provider addrs.Provider) (VersionList, Warnings, error)
	PackageMeta(ctx context.Context, provider addrs.Provider, version Version, target Platform) (PackageMeta, error)
	ForDisplay(provider addrs.Provider)
}
