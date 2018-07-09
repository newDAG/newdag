package store

import (
	"errors"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const newdag_dbFile = "newdag_chain.db"

type BoltDB = bolt.DB
type BoltTx = bolt.Tx

func DbGetvalue(stable string, db *BoltDB, key []byte) ([]byte, error) {
	var itemBytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(stable))
		itemBytes := b.Get(key)
		if itemBytes == nil {
			return errors.New("@KeyNotFound")
		}
		return nil
	})

	return itemBytes, err
}

func OpenKvDatabase(flag int) *BoltDB {
	if flag >= 0 {
		if _, err := os.Stat(newdag_dbFile); os.IsNotExist(err) {
			log.Panic(err)
		}
	} else if flag < 0 {
		//清空数据库
	}
	db, err := bolt.Open(newdag_dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	return db
}

type CallbackIter func(k, v []byte, p1, p2 interface{}) (bool, error)

func EachAllDBView(db *BoltDB, tableName string, callback CallbackIter, p1, p2 interface{}) error {
	err := db.View(func(tx *BoltTx) error {
		b := tx.Bucket([]byte(tableName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ret, err := callback(k, v, p1, p2)
			if ret == true || err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
