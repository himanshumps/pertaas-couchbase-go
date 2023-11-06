package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"log"
	"os"
	"time"
)

func main() {
	//gocb.SetLogger(gocb.VerboseStdioLogger())
	var jobId string
	var message string
	flag.StringVar(&jobId, "setJobId", "", "")
	flag.StringVar(&message, "setMessage", "", "") // new added
	flag.Parse()
	keyTx := jobId + "::" + time.Now().Format(time.RFC3339)
	type Message struct {
		KeyTx   string `json:"key_tx"`
		JobId   string `json:"jobId"`
		Message string `json:"message"`
	}
	connectionString := os.Getenv("COUCHBASE_CONNECTION_STRING")
	username := os.Getenv("COUCHBASE_USER")
	password := os.Getenv("COUCHBASE_PASSWORD")
	bucketName := os.Getenv("COUCHBASE_BUCKET")
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
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
