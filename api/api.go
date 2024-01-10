package api

import (
	. "MusicPlayer/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
	"log"
	"os"
)

// https://www.thepolyglotdeveloper.com/2020/02/interacting-with-a-graphql-api-with-golang/

// MetadataRequest fetches metadata from the graphbrainz API
func MetadataRequest(client *graphql.Client, request *graphql.Request) {
	qFilePath := "qFile.json"

	qFile, err := os.OpenFile(qFilePath, os.O_CREATE|os.O_EXCL, 0666)
	if os.IsNotExist(err) {
		log.Fatalf("Error creating %v: %v", qFile.Name(), err)
	} else if os.IsExist(err) {
		Info("Did not create file, %v already exists.", qFilePath)
		qFile.Close()
		return
	}
	defer qFile.Close()

	var response interface{}
	if err := client.Run(context.Background(), request, &response); err != nil {
		log.Fatal(err)
	}

	jsonData, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(qFile, "%s", string(jsonData))
}
