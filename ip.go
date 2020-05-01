package echo

import "net"

type ipChecker struct {
	trustLoopback    bool
	trustLinkLocal   bool
	trustPrivateNet  bool
	trustExtraRanges []*net.IPNet
}
