package connections

import (
	"github.com/boltdb/bolt"
)

var Database *bolt.DB

func StartBoltDb() error {
	db, err := bolt.Open("storageUsers.db", 0600, nil)
	if err != nil {
		return err
	}
	Database = db
	//defer func(db *bolt.DB) {
	//	err = db.Close()
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//}(db)

	return nil
}
