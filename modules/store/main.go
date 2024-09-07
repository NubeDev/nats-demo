package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Enable JetStream
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error enabling JetStream: %v", err)
	}

	// Create or open an existing Object Store bucket
	objStore, err := js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket:  "myobjbucket",    // Name of the object store bucket
		Storage: nats.FileStorage, // Ensure files are stored on disk
	})
	if err != nil {
		objStore, err = js.ObjectStore("myobjbucket")
	}

	// File to upload
	filePath := "/home/aidan/tinygo_0.30.0_amd64.deb"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Get the file name
	fileName := filepath.Base(filePath)

	// Upload the file to the Object Store using an io.Reader
	obj, err := objStore.Put(&nats.ObjectMeta{Name: fileName}, file)
	if err != nil {
		log.Fatalf("Error uploading file to object store: %v", err)
	}

	// Print confirmation of upload
	fmt.Printf("File %s uploaded to Object Store with digest: %s\n", fileName, obj.Digest)
}
