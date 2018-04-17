package code

import "github.com/nanobox-io/golang-scribble"

func getDatabase() *scribble.Driver {
	db, err := scribble.New("./database/", nil)
	Check(err)
	return db
}