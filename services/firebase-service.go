package services

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type File interface {
	FirebaseInit() *firebase.App
	Upload(b64Name string, fileName string) 
	Download() 
}

type file struct {
}

func NewFile() File {
	return &file{}
}

func (f *file) FirebaseInit() *firebase.App {
	config := &firebase.Config{
		StorageBucket: os.Getenv("BUCKET"),
	}
	opt := option.WithCredentialsFile(os.Getenv("CREDENTIAL"))
	fb, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Println(err)
	}

	return fb
}

func (f *file) Upload(b64Name string, fileName string) {
	fb := f.FirebaseInit()
	client, err := fb.Storage(context.Background())
	if err != nil {
		log.Println(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Println(err)
	}
	i := strings.Index(b64Name, ",")
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64Name[i+1:]))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	wc := bucket.Object(fileName).NewWriter(ctx)
	_, err = io.Copy(wc, reader)
	if err != nil {
		log.Println(err.Error())
	}

	if err := wc.Close(); err != nil {
		log.Println(err.Error())
	}
}

func (f *file) Download() {
}
