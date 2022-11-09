package main

import "runtime/pprof"
import "os"
import "time"

func main() {
	go leakyFunctionInline()
	time.Sleep(500 * time.Millisecond)
	f, _ := os.Create("/tmp/profile.pb.gz")
	defer f.Close()
	pprof.WriteHeapProfile(f)
	pprof.StartCPUProfile(f)
}

func leakyFunctionInline() {
	s := make([]string, 3)
	for i := 0; i < 10000000; i++ {
		s = append(s, "magical pprof time")
	}
}
