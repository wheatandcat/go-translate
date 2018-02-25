package lib

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"
)

// Word Word
type Word struct {
	Words      []string  `json:"words"`
	Translates []string  `json:"translates"`
	Lang       string    `json:"lang"`
	Updated    time.Time `json:"-"`
	Created    time.Time `json:"-"`
}

func NewWord() *Word {
	w := &Word{
		Words:      []string{},
		Translates: []string{},
		Lang:       "",
		Created:    time.Now(),
	}

	return w
}

func (w Word) Put(c context.Context) error {
	w.Updated = time.Now()

	key := datastore.NewIncompleteKey(c, "Word", nil)

	if _, err := datastore.Put(c, key, &w); err != nil {
		return err
	}

	return nil
}
