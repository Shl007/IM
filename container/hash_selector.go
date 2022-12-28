package container

import (
	"im"
	"im/wire/pkt"
)

// HashSelector HashSelector
type HashSelector struct {
}

// Lookup a server
func (s *HashSelector) Lookup(header *pkt.Header, srvs []im.Service) string {
	ll := len(srvs)
	code := HashCode(header.ChannelId)
	return srvs[code%ll].ServiceID()
}
