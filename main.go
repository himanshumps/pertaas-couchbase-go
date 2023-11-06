package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

func main() {
	//gocb.SetLogger(gocb.VerboseStdioLogger())
	var jobId string
	var message string
	type Message struct {
		KeyTx   string `json:"key_tx"`
		JobId   string `json:"jobId"`
		Message string `json:"message"`
	}
	connectionString := os.Getenv("couchbaseConnectionString")
	username := os.Getenv("couchbaseUsername")
	password := os.Getenv("couchbasePassword")
	bucketName := os.Getenv("couchbaseBucket")
	cluster, err := gocb.Connect("couchbase://"+connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = cluster.WaitUntilReady(30*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(30*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&jobId, "setJobId", "", "")
	flag.StringVar(&message, "setMessage", "", "") // new added
	flag.Parse()
	keyTx := jobId + "::" + uuid.New().String()
	fmt.Printf("keyTx: %v\n", keyTx)
	messageVar := Message{
		KeyTx:   keyTx,
		JobId:   jobId,
		Message: message,
	}
	_, err = bucket.DefaultCollection().Upsert(keyTx, &messageVar, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = cluster.Close(nil)
}
