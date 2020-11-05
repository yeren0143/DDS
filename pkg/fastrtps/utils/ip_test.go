package utils

import (
	"fmt"
	"testing"
)

func TestGetIP4Address(t *testing.T) {
	addrs := GetIP4Address()
	for _, address := range addrs {
		fmt.Println("ipv4 address:", address.Address)
	}
}

func TestGetHost(t *testing.T) {
	host := GetHost()
	fmt.Println("host:", host.id)
}

// func TestGetNetDevStats(t *testing.T) {
// 	getNetDevStats()
// }
