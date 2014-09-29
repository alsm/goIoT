package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/bbb"
	"time"
)

type BeagleBone struct {
	LEDs []int
}

func NewBeagleBone() *BeagleBone {
	embd.InitLED()
	return &BeagleBone{LEDs: []int{0, 1, 2, 3}}
}

func (b *BeagleBone) LedsOn() {
	for _, l := range b.LEDs {
		embd.LEDOn(l)
	}
}

func (b *BeagleBone) LedsOff() {
	for _, l := range b.LEDs {
		embd.LEDOff(l)
	}
}

func (b *BeagleBone) LedsToggle() {
	for _, l := range b.LEDs {
		embd.LEDToggle(l)
	}
}

func (b *BeagleBone) LedsCycle(repeat int) {
	b.LedsOff()
	for i := 0; i < repeat; i++ {
		for _, l := range b.LEDs {
			embd.LEDOn(l)
			time.Sleep(50 * time.Millisecond)
			embd.LEDOff(l)
		}
	}
}

func (b *BeagleBone) Close() {
	embd.CloseLED()
}

func (b *BeagleBone) CPUTempFile() string {
	return "/sys/class/hwmon/hwmon0/device/temp1_input"
}
