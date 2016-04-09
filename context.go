package crossbot

import (
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	data  = make(map[State]interface{})
	datat = make(map[State]int64)
)

// Set stores a value for a given key in a given request.
func Set(s State, key, val interface{}) {
	mutex.Lock()
	if data[s] == nil {
		datat[s] = time.Now().Unix()
	}
	data[s] = val
	mutex.Unlock()
}

// Get returns a value stored for a given key in a given request.
func Get(s State) interface{} {
	mutex.RLock()
	if ctx := data[s]; ctx != nil {
		value := ctx
		mutex.RUnlock()
		return value
	}
	mutex.RUnlock()
	return nil
}
