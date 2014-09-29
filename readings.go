package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func getLoadAvg() (float64, float64, float64) {
	rawData, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, 0
	}
	fields := strings.Fields(string(rawData))
	One, _ := strconv.ParseFloat(fields[0], 64)
	Five, _ := strconv.ParseFloat(fields[1], 64)
	Fifteen, _ := strconv.ParseFloat(fields[2], 64)

	return One, Five, Fifteen
}

func getCPUTemp(cpuTempFile string) float64 {
	rawData, err := ioutil.ReadFile(cpuTempFile)
	if err != nil {
		return 0
	}
	cpuTempMilli, _ := strconv.ParseInt(strings.TrimSpace(string(rawData)), 0, 0)
	return (float64(cpuTempMilli) / 1000)
}
