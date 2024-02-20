package store

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Store struct {
	db *badger.DB
}

func NewStore(dbPath string) *Store {
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &Store{db: db}
}

func (s *Store) SaveDocument(id string, pdfContent []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(id), pdfContent)
		return err
	})
	return err
}

func (s *Store) GetDocument(id string) ([]byte, error) {
	var pdfContent []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			pdfContent = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return pdfContent, nil
}
