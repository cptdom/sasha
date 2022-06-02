package main

import (
	"flag"
	"fmt"
	"os"

	"cptdom/sasha/cmd/handler"
	s3s "cptdom/sasha/cmd/s3"
	"cptdom/sasha/version"
)

func main() {
	// get endpoint
	endpointUrl := flag.String("e", "", "S3 endpoint URL. Required. Example value: https://prop.s3.loc.dom.com")
	flag.Parse()
	if *endpointUrl == "" {
		fmt.Println("Missing endpoint URL. Check 'sasha --help'.")
		os.Exit(1)
	}
	// initiate session
	fmt.Printf("Sasha v%v\n", version.Version)
	s := s3s.SessionClient{
		Session: *s3s.GetSession(*endpointUrl),
	}
	// run the loop
	handler.RunLoop(s)
}
