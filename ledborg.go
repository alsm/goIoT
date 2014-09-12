package main

import (
	"fmt"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"time"
)

type RPi struct {
	LEDs []string
}

func NewRPi() *RPi {
	pi := &RPi{LEDs: []string{"GPIO_17", "GPIO_27", "GPIO_22"}}
	err := embd.InitGPIO()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, l := range pi.LEDs {
		err := embd.SetDirection(l, embd.Out)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return pi
}

func (r *RPi) LedsOn() {
	for _, l := range r.LEDs {
		embd.DigitalWrite(l, embd.High)
	}
}

func (r *RPi) LedsOff() {
	for _, l := range r.LEDs {
		embd.DigitalWrite(l, embd.Low)
	}
}

func (r *RPi) LedsToggle() {

}

func (r *RPi) LedsCycle(repeat int) {
	r.LedsOff()
	for i := 0; i < repeat; i++ {
		for _, l := range r.LEDs {
			embd.DigitalWrite(l, embd.High)
			time.Sleep(50 * time.Millisecond)
			embd.DigitalWrite(l, embd.Low)
		}
	}
}

func (r *RPi) Close() {
	embd.CloseGPIO()
}
