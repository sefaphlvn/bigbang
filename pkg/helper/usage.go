package helper

import (
	"fmt"
	"runtime"
)

// MemoryStats gösterir
func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

// CPU Gösterir
func PrintCPUUsage() {
	numCPU := runtime.NumCPU()
	numGoroutine := runtime.NumGoroutine()
	fmt.Printf("CPU = %d, Goroutines = %d\n", numCPU, numGoroutine)
}

// Byte to MB dönüşümü için yardımcı fonksiyon
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
