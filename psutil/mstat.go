package psutil

import (
	"runtime"
)

func GoMemory() *GoMemoryStat {

	mstat := &runtime.MemStats{}
	runtime.ReadMemStats(mstat)

	return &GoMemoryStat{
		Alloc:      mstat.Alloc,
		Sys:        mstat.Sys,
		HeapAlloc:  mstat.HeapAlloc,
		HeapInuse:  mstat.HeapInuse,
		HeapSys:    mstat.HeapSys,
		StackInuse: mstat.StackInuse,
		StackSys:   mstat.StackSys,
		TotalAlloc: mstat.TotalAlloc,
		LastGC:     uint64(mstat.LastGC / 1e9),
		NumGC:      mstat.NumGC,
	}

}
