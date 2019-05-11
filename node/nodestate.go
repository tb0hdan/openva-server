package node

import (
	"sync"

	"github.com/tb0hdan/openva-server/api"
)

type State struct {
	PlayerState       api.PlayerStateMessage
	SystemInformation api.SystemInformationMessage
	LastUpdatedTS     int64
}

type StateType struct {
	nodeState map[string]*State
	m         sync.RWMutex
}

func (mc *StateType) Get(key string) (value *State, ok bool) {
	mc.m.RLock()
	value, ok = mc.nodeState[key]
	mc.m.RUnlock()
	return
}

func (mc *StateType) Set(key string, value *State) {
	mc.m.Lock()
	mc.nodeState[key] = value
	mc.m.Unlock()
}

func (mc *StateType) Len() (stateSize int) {
	stateSize = len(mc.nodeState)
	return
}

func (mc *StateType) All() map[string]*State {
	return mc.nodeState
}

func (mc *StateType) Delete(key string) {
	delete(mc.nodeState, key)
}

func New() (nodeState *StateType) {
	nodeState = &StateType{nodeState: make(map[string]*State)}
	return
}
