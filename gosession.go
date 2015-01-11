//
// ---- Package GoSession ----
// Provides an interface for creating Bolt DB backed sessions.
//
// To use:
// - Open a bold db connection
// - Run the Boltstore reaper to remove expired sessions
//
// db, err = bolt.Open("./sessions.db", 0666, nil)
// if err != nil {
//   panic(err)
// }
// defer db.Close()
// // Invoke a reaper which checks and removes expired sessions periodically.
// defer reaper.Quit(reaper.Run(db, reaper.Options{}))
//

package gosession

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	store "github.com/yosssi/boltstore/store"
	"log"
	"net/http"
	"time"
)

var db *bolt.DB

func handlErr(err error) {
	start := time.Now()
	log.Printf(
		"%s\t%s",
		err,
		time.Since(start),
	)
}

func newStore() *store.Store {
	config := store.Config{}
	newStore, err := store.New(db, config, []byte("#_B43D+=,n8S26L240FCPjc<,P0Fy*HZ2@Hu-aZ-()h3j713oL{uk+S^x3e52er"))
	if err != nil {
		handlErr(err)
	}
	return newStore
}

// New creates and saves a new session on the response and in the database.
func New(w http.ResponseWriter, r *http.Request, name string, values map[interface{}]interface{}) {
	// Create a store.
	str := newStore()

	// Get a session.
	session, err := str.New(r, name)
	if err != nil {
		handlErr(err)
	}

	// Add a value on the session.
	session.Values = values

	// Save the session.
	if err := sessions.Save(r, w); err != nil {
		handlErr(err)
	}
}

// Refresh updates the timestam of the current session.
func Refresh(w http.ResponseWriter, r *http.Request, name string) {
	// Create a store.
	str := newStore()

	// Get a session.
	session, err := str.Get(r, name)
	if err != nil {
		handlErr(err)
	}
	// Save the session.
	if err := sessions.Save(r, w); err != nil {
		handlErr(err)
	}
	log.Println(session)
}

// Future: Update() to pass in new values
