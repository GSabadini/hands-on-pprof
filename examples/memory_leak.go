package main

import (
	"database/sql"
	"expvar"
	_ "expvar"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"
)

// Publish the port as soon as the program starts
func init() {
	http.Handle("/leak", LeakyFunction())
	http.Handle("/leak-not-closed", LeakyFunction1())
	//http.Handle("/leak-query", Query(db))

	go http.ListenAndServe(":8080", nil)
}

// Custom struct that will be exported
type Load struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// Function that will be called by expvar
// to export the information from the structure
// every time the endpoint is reached
func AllLoadAvg() interface{} {
	return Load{
		Load1:  loadAvg(0),
		Load5:  loadAvg(1),
		Load15: loadAvg(2),
	}
}

// Aux function to retrieve the load average
// in GNU/Linux systems
func loadAvg(position int) float64 {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		panic(err)
	}
	values := strings.Fields(string(data))

	load, err := strconv.ParseFloat(values[position], 64)
	if err != nil {
		panic(err)
	}

	return load
}

func main() {
	//db, err := newMySQLHandler()
	//if err != nil {
	//	log.Println(err)
	//}

	var (
		numberOfSecondsRunning = expvar.NewInt("system.numberOfSeconds")
		programName            = expvar.NewString("system.programName")
		lastLoad               = expvar.NewFloat("system.lastLoad")
		numberOfLoginsPerUser  = expvar.NewMap("system.numberOfLoginsPerUser")
	)

	// The contents returned by the function will be autoexported in JSON format
	expvar.Publish("system.allLoad", expvar.Func(AllLoadAvg))

	programName.Set(os.Args[0])

	// We will increment this metrics every second
	for {
		numberOfSecondsRunning.Add(1)
		lastLoad.Set(loadAvg(0))
		numberOfLoginsPerUser.Add("foo", 2)
		numberOfLoginsPerUser.Add("bar", 1)
		time.Sleep(1 * time.Second)
	}
}

var (
	count int
)

func Query(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT count FROM stats")
		if err != nil {
			log.Fatal(err)
		}

		//defer rows.Close()
		if rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(count)
		}
	}
}

func LeakyFunction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := make([]string, 3)
		for i := 0; i < 10000000; i++ {
			s = append(s, "magical pprof time")
		}
	}
}

func LeakyFunction1() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := http.NewRequest(http.MethodGet, "https://reqres.in/api/users", nil)
		client := http.Client{}
		_, err := client.Do(req)
		//if resp != nil {
		//	defer resp.Body.Close()
		//}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte("OK"))
	}
}

func newMySQLHandler() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/app")
	if err != nil {
		return &sql.DB{}, err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS stats (name VARCHAR(50) PRIMARY KEY, count INTEGER);"); err != nil {
		return &sql.DB{}, err
	}

	//for x:=1; x < 20000; x++ {
	//	name := fmt.Sprintf("GABRIEL-%d", x)
	//	if _, err := db.Exec("INSERT INTO stats (name, count) VALUES (?, 10);", name); err != nil {
	//		return &sql.DB{}, err
	//	}
	//}

	log.Println("Database connected")

	return db, nil
}
