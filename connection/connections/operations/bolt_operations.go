package operations

// unfinished, new update

//func createBoltBucketIfNotExists(u *models.HubUser, tx *bolt.Tx) error {
//	_, err := tx.CreateBucketIfNotExists([]byte(minio_operations.GetUserBucketName(u)))
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func GetUserQuota(u *models.HubUser) (uint64, error) {
//	err := connections.Database.Update(func(tx *bolt.Tx) error {
//		err := createBoltBucketIfNotExists(u, tx)
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return -1, err
//	}
//
//	quotaType: string
//	quotaSpecificValue: uint64()
//	err = connections.Database.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(minio_operations.GetUserBucketName(u)))
//		quotaTypeBytes := b.Get([]byte("quotaType"))
//		if len(quotaTypeBytes) == 0 {
//			err = b.Put([]byte("quotaType"), []byte("default")) // there can be a quota type keyword or just a number in a string form
//			if err != nil {
//				return err
//			}
//		}
//
//		quotaSpecificValueBytes := b.Get([]byte("quotaSpecific")) // can be not set, only if really needed
//		quotaSpecificValue, error := strconv.ParseUint(string(quotaSpecificValueBytes), 10, 64)
//		if err != nil {
//			log.Println("Found a non-uint64 value at db key quotaSpecific for user id =", u.Id)
//			return err
//		}
//
//		result = GetExactQuotaForType(quotaType)
//
//		return nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return nil
//}
