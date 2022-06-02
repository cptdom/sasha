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
	// flags
	endpointUrl := flag.String(
		"e",
		"",
		`S3 entrypoint URL. Required. Can be stored in env as 'S3_ENTRYPOINT'.
		Example value: https://prop.s3.loc.dom.com`,
	)
	flag.Parse()
	if *endpointUrl == "" {
		if val, ok := os.LookupEnv("S3_ENTRYPOINT"); !ok {
			fmt.Println("Missing entrypoint URL. Check 'sasha --help'.")
			os.Exit(1)
		} else {
			endpointUrl = &val
		}
	}
	// initiate session
	fmt.Printf("Sasha v%v\n", version.Version)
	s := s3s.SessionClient{
		Session: *s3s.GetSession(*endpointUrl),
	}
	// run the loop
	handler.RunLoop(s)
}
