package utils

import "sync"

type Host struct {
	id uint16
}

var hostInstance *Host
var hostMutex sync.Mutex

func GetHost() *Host {
	if hostInstance == nil {
		hostMutex.Lock()
		defer hostMutex.Unlock()

		if hostInstance == nil {
			hostInstance = Host{}
		}
	}
	return hostInstance
}
