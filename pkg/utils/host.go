package utils

import (
	"sync"
)

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
			hostInstance = &Host{}
			addrs := GetIP4Address()

			if len(addrs) == 0 {
				hostInstance.id |= 127 << 8
				hostInstance.id |= 1
			} else {
				// var data []byte
				// for _, loc := range addrs {
				// 	data = append(data, loc.Address[0:16]...)
				// }
				// digest := md5.Sum(data)
				// for i := 0; i < len(digest); i += 2 {
				// 	hostInstance.id ^= ((uint16(digest[i]) << 8) | uint16(digest[i+1]))
				// }
				md5 := NewMd5()
				for _, loc := range addrs {
					var input []byte
					for i := 0; i < len(loc.Address); i += 1 {
						input = append(input, loc.Address[i])
					}
					md5.Update(input, uint32(len(input)))
				}
				md5.Finalize()
				for i := 0; i < len(md5.digest); i += 2 {
					a := uint16(md5.digest[i]) << 8
					b := uint16(md5.digest[i+1])
					hostInstance.id ^= (a | b)
				}
			}
		}
	}
	return hostInstance
}
