package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var isDebug bool
var myLockMap = make(map[string]string)

func debug(s string) {
	if isDebug {
		log.Print(s)
	}
}

func listLocks() (s string) {
	debug(">> listLocks")
	s = ""
	if len(myLockMap) == 0 {
		return "[]"
	}
	for k, v := range myLockMap {
		debug(fmt.Sprintf(">>> k: %s, v: %s", k, v))
		if len(s) == 0 {
			s = fmt.Sprintf("%s", k)
		} else {
			s = fmt.Sprintf("%s,\n  %s", s, k)
		}
	}
	debug(fmt.Sprintf(">>> s: %s", s))
	return fmt.Sprintf("[ %s ]", s)
}

func storeLock(id string) (e error) {
	myLockMap[id] = id
	debug(fmt.Sprintf(">> storeLock: storing %s", id))
	err := getLock(id)
	if err != nil {
		return errors.New("storeLock: unabled to store lock")
	}
	debug(fmt.Sprintf(">> storeLock: lock %s successfully stored", id))
	return nil
}

func getLock(id string) (e error) {
	_, ok := myLockMap[id]
	if !ok {
		debug(fmt.Sprintf(">> getLock: lock %s not found", id))
		return errors.New("getLock: lock not found")
	}
	debug(fmt.Sprintf(">> getLock: lock %s found", id))
	return nil
}

func deleteLock(id string) (e error) {
	delete(myLockMap, id)
	err := getLock(id)
	if err == nil {
		debug(fmt.Sprintf(">> deleteLock: lock %s not deleted", id))
		return errors.New("deleteLock: unable to delete lock")
	}
	debug(fmt.Sprintf(">> deleteLock: lock %s deleted", id))
	return nil
}

func createLockHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	log.Print("creating ", id)
	err := getLock(id)
	if err == nil {
		debug(fmt.Sprintf(">> lock %s already existed", id))
		http.Error(w, "Lock already exists", http.StatusConflict)
	} else {
		err = storeLock(id)
		if err != nil {
			debug(fmt.Sprintf(">> lock %s not created", id))
			http.Error(w, "Unable to create lock", http.StatusInternalServerError)
		} else {
			debug(fmt.Sprintf(">> lock %s created", id))
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Lock %s created\n", id)
		}
	}
}

func deleteLockHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	log.Print("deleting ", id)
	err := deleteLock(id)
	if err != nil {
		debug(fmt.Sprintf("> lock %s not deleted", id))
		http.Error(w, "Unable to delete lock", http.StatusInternalServerError)
	} else {
		debug(fmt.Sprintf("> lock %s deleted", id))
		fmt.Fprintf(w, "Lock %s deleted\n", id)
	}
}

func infoLockHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	log.Print("info ", id)

	if id == "" {
		debug("> listing all locks")
		fmt.Fprintf(w, listLocks())
	} else {
		err := getLock(id)
		if err != nil {
			debug(fmt.Sprintf("> no such lock %s", id))
			http.Error(w, "Lock not found", http.StatusNotFound)
		} else {
			debug(fmt.Sprintf("> lock info %s", id))
			fmt.Fprintf(w, id)
		}
	}
}

// route lock operations
func lockRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		infoLockHandler(w, r)
	} else if r.Method == "PUT" {
		createLockHandler(w, r)
	} else if r.Method == "DELETE" {
		deleteLockHandler(w, r)
	} else {
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func main() {
	flag.BoolVar(&isDebug, "d", false, "-d to activate debug logs")
	flag.Parse()

	// register URL handlers
	http.HandleFunc("/", lockRouter)

	// start the server
	log.Print("starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
