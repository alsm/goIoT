package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

type Platform interface {
	LedsOn()
	LedsOff()
	LedsToggle()
	LedsCycle(int)
	Close()
}

func NewPlatform(host embd.Host) Platform {
	switch host {
	case embd.HostBBB:
		return NewBeagleBone()
	case embd.HostRPi:
		return NewRPi()
	}
	return nil
}
