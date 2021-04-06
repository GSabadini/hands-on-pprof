package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
)

func main() {
	http.Handle("/leak", leak1())
	http.Handle("/not-leak", not_leak1())
	log.Println(http.ListenAndServe(":6060", nil))
}

// Goroutine é criado, que bloqueia na linha 29 esperando para receber um valor do canal.
// Enquanto esse Goroutine está esperando, a leak função retorna.
// Neste ponto, nenhuma outra parte do programa pode enviar um sinal pelo canal.
// Isso deixa o Goroutine bloqueado na linha 29 esperando indefinidamente.
// A fmt.Printlnchamada na linha 30 nunca acontecerá.
func leak() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(runtime.NumGoroutine())))
		ch := make(chan int)

		go func() {
			val := <-ch
			fmt.Println("We received a value:", val)
		}()
	}
}

func not_leak() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(runtime.NumGoroutine())))
		ch := make(chan int)

		go func (ch chan int) {
			ch <- 1
		}(ch)

		go func(ch chan int) {
			val := <- ch
			fmt.Println("We received a value:", val)
		}(ch)

	}
}

func leak1() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(runtime.NumGoroutine())))

		go func() {
			var wg sync.WaitGroup
			wg.Add(1)
			wg.Wait()
		}()
	}
}

func not_leak1() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(runtime.NumGoroutine())))

		go func() {
			var wg sync.WaitGroup
			wg.Add(1)
			wg.Done()
			wg.Wait()
		}()
	}
}