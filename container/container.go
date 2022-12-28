package container

import (
	"im"
	"im/logger"
	"im/naming"
	"sync"
)

const (
	stateUninitialized = iota
	stateInitialized
	stateStarted
	stateClosed
)

// Container Container
type Container struct {
	sync.RWMutex
	Naming     naming.Naming
	Srv        im.Server
	state      uint32
	srvclients map[string]ClientMap
	selector   Selector
	dialer     im.Dialer
	deps       map[string]struct{}
}

var log = logger.WithField("module", "container")

// Default Container
var c = &Container{
	state:    0,
	selector: &HashSelector{},
	deps:     make(map[string]struct{}),
}

// Default Default
func Default() *Container {
	return c
}
