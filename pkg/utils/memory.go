package utils

import (
	"log"
	"runtime"
	"syscall"
)

func LogMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Memory used: Alloc = %v KiB\tTotalAlloc = %v KiB\tSys = %v KiB\tNumGC = %v\n",
		bToKb(m.Alloc),
		bToKb(m.TotalAlloc),
		bToKb(m.Sys),
		m.NumGC,
	)
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

func GetSysTotalMemory() uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	return uint64(in.Totalram) * uint64(in.Unit)
}

func GetSysFreeMemory() uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	return uint64(in.Freeram) * uint64(in.Unit)
}
