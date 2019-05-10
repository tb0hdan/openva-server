package node

import (
	"sync"

	"github.com/tb0hdan/openva-server/api"
)

type NodeState struct {
	PlayerState       api.PlayerStateMessage
	SystemInformation api.SystemInformationMessage
	LastUpdatedTS     int64
}

type NodeStateType struct {
	nodeState map[string]*NodeState
	m         sync.RWMutex
}

func (mc *NodeStateType) Get(key string) (value *NodeState, ok bool) {
	mc.m.RLock()
	value, ok = mc.nodeState[key]
	mc.m.RUnlock()
	return
}

func (mc *NodeStateType) Set(key string, value *NodeState) {
	mc.m.Lock()
	mc.nodeState[key] = value
	mc.m.Unlock()
}

func (mc *NodeStateType) Len() (stateSize int) {
	stateSize = len(mc.nodeState)
	return
}

func (mc *NodeStateType) All() map[string]*NodeState {
	return mc.nodeState
}

func (mc *NodeStateType) Delete(key string) {
	delete(mc.nodeState, key)
}

func New() (nodeState *NodeStateType) {
	nodeState = &NodeStateType{nodeState: make(map[string]*NodeState)}
	return
}
