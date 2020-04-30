package main

import (
    "fmt"
    "github.com/CharLemAznable/gokits"
    "go.etcd.io/bbolt"
)

var db *bbolt.DB

const PayloadBucket = "Payload"

func init() {
    _db, err := bbolt.Open("./sonar.db", 0666, nil)
    if err != nil {
        gokits.LOG.Crashf("DB create error: %s", err.Error())
    }
    db = _db

    err = db.Update(func(tx *bbolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(PayloadBucket))
        if err != nil {
            return fmt.Errorf("create bucket "+PayloadBucket+": %s", err.Error())
        }
        return nil
    })
    if err != nil {
        gokits.LOG.Crashf("DB init error: %s", err.Error())
    }
}
