package psutil

import (
	"runtime"
)

func GoMemory() *GoMemoryStat {

	mstat := &runtime.MemStats{}
	runtime.ReadMemStats(mstat)

	return &GoMemoryStat{
		Alloc:     mstat.Alloc,
		Sys:       mstat.Sys,
		HeapAlloc: mstat.HeapAlloc,
		HeapSys:   mstat.HeapSys,
		LastGC:    uint64(mstat.LastGC / 1e9),
		NumGC:     mstat.NumGC,
	}

}
