package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cptdom/sasha/cmd/location"
	s3s "cptdom/sasha/cmd/s3"
	"cptdom/sasha/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func RunLoop(s s3s.SessionClient) {
	// initiate wd
	level := location.Location{
		Session: s.Session,
		S3: *s3.New(&s.Session, &aws.Config{
			S3ForcePathStyle: utils.BoolAddr(true),
		}),
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		level.Update() // this sometimes takes notable time
		fmt.Printf("\n%v@%v %v %% ", s.GetAccessID(), s.GetEndpoint(), level.Repr())
		command, _, _ := reader.ReadLine()
		strCommand := strings.Trim(string(command), " ")
		splitCommand := strings.Split(strCommand, " ")
		if len(splitCommand) == 1 {
			switch strCommand {
			case "ls":
				level.LS()
			case "pwd":
				level.PWD()
			case "whoami":
				s.Whoami()
			case "cd":
				level.Reset()
			case "update", "reload":
				level.Update()
			case "exit", "Exit":
				os.Exit(0)
			default:
				fmt.Println("Unknown command.")
			}
		}
		if len(splitCommand) == 2 {
			switch splitCommand[0] {
			case "cd":
				if splitCommand[1] == ".." {
					level.Pop(1)
					continue
				}
				level.Add(strings.Split(strings.Trim(splitCommand[1], "/"), "/")[0])
			case "file":
				level.File(splitCommand[1])
			default:
				fmt.Println("Unknown command.")
			}
		}
	}
}
