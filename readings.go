package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	loadAvgFile string = "/proc/loadavg"
	cpuTempFile string = "/sys/class/hwmon/hwmon0/device/temp1_input"
	cpuinfoFile string = "/proc/cpuinfo"
)

func getLoadAvg() (float64, float64, float64) {
	rawData, err := ioutil.ReadFile(loadAvgFile)
	if err != nil {
		return 0, 0, 0
	}
	fields := strings.Fields(string(rawData))
	One, _ := strconv.ParseFloat(fields[0], 64)
	Five, _ := strconv.ParseFloat(fields[1], 64)
	Fifteen, _ := strconv.ParseFloat(fields[2], 64)

	return One, Five, Fifteen
}

func getCPUTemp() float64 {
	rawData, err := ioutil.ReadFile(cpuTempFile)
	if err != nil {
		return 0
	}
	cpuTempMilli, _ := strconv.ParseInt(strings.TrimSpace(string(rawData)), 0, 0)
	return (float64(cpuTempMilli) / 1000)
}
