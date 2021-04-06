package main

import (
	_ "net/http/pprof"
	"net/http"
)

func main() {
	http.ListenAndServe(":6061", nil)
}
