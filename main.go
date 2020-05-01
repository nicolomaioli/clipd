package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"

	"github.com/nicolomaioli/clipd/middleware"
)

// Clip represents a message in a register
type Clip struct {
	Reg     string `json:"reg,omitempty"`
	Content string `json:"clip"`
}

// Mem holds clipd's memory
var Mem = map[string]string{}

// MemMut is used to lock and unlock Mem
var MemMut = &sync.RWMutex{}

// DebugLogger is used to log debugs to console
var DebugLogger = log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime)

// ErrorLogger is used to log errors to console
var ErrorLogger = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime)

const (
	// DefaultRegister is the default clipboard register
	DefaultRegister = "0"
)

func yank(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Print("error reading body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yr := &Clip{}
	json.Unmarshal(body, yr)
	if yr.Reg == "" {
		yr.Reg = DefaultRegister
	}

	MemMut.Lock()
	Mem[yr.Reg] = yr.Content
	MemMut.Unlock()
	w.WriteHeader(http.StatusOK)
}

func paste(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reg := p.ByName("reg")
	if reg == "" {
		reg = DefaultRegister
	}

	var clip string

	MemMut.RLock()
	if v, ok := Mem[reg]; ok {
		clip = v
	}
	MemMut.RUnlock()

	if clip == "" {
		DebugLogger.Printf("clip not found with reg %q", reg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	yr := &Clip{
		Reg:     reg,
		Content: clip,
	}

	b, err := json.Marshal(yr)
	if err != nil {
		ErrorLogger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func main() {
	router := httprouter.New()
	router.POST("/yank", yank)
	router.GET("/paste/:reg", paste)
	router.GET("/paste", paste)

	addr := ":8080"
	l := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	l.Printf("Server listening at %s", addr)

	log.Fatal(
		http.ListenAndServe(
			addr,
			middleware.RequestLogger{Next: router, Logger: l},
		),
	)
}
