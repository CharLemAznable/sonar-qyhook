package main

import (
    "fmt"
    "github.com/kataras/golog"
    "go.etcd.io/bbolt"
)

var db *bbolt.DB

const PayloadBucket = "Payload"

func init() {
    _db, err := bbolt.Open("./sonar.db", 0666, nil)
    if err != nil {
        golog.Errorf("DB create error: %s", err.Error())
        panic("DB create error")
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
        golog.Errorf("DB init error: %s", err.Error())
        panic("DB init error")
    }
}
