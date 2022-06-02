package location

import (
	"fmt"
	"os"
	"sort"
	"strings"

	s3s "cptdom/sasha/cmd/s3"
	"cptdom/sasha/utils"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var header = fmt.Sprintf("\n%s %+59s\n%v\n", "Name", "Last modified", utils.CreateLine())

type Location struct {
	Session       session.Session
	S3            s3.S3
	Level         []string
	Bucket        string
	BucketsBuffer map[string]*s3.Bucket
	DirBuffer     map[string]struct{}
	FileBuffer    map[string]struct{}
}

func (Location *Location) Add(newLevel string) {
	if Location.Bucket == "" {
		if _, ok := Location.BucketsBuffer[newLevel]; !ok {
			fmt.Printf("%v: No such bucket.\n", newLevel)
			return
		}
		Location.Bucket = newLevel
		return
	}
	var absDirPath string
	if len(Location.Level) == 0 {
		absDirPath = newLevel + "/"
	} else {
		absDirPath = strings.Join(Location.Level, "/") + "/" + newLevel + "/"
	}
	if _, ok := Location.DirBuffer[absDirPath]; !ok {
		fmt.Printf("%v: No such directory.\n", newLevel)
		return
	}
	Location.Level = append(Location.Level, newLevel)
}

// steps attribute for cd ../.. option
func (Location *Location) Pop(steps int) {
	if len(Location.Level) != 0 {
		Location.Level = Location.Level[:len(Location.Level)-steps]
		return
	}
	Location.Bucket = ""
}

func (Location *Location) Repr() string {
	if Location.Bucket == "" {
		return "~/"
	}
	if len(Location.Level) == 0 {
		return "/B:" + Location.Bucket + strings.Join(Location.Level, "/") + "/"
	}
	return "/B:" + Location.Bucket + "/" + strings.Join(Location.Level, "/") + "/"
}

func (Location *Location) PWD() {
	locationRepr := Location.Repr()
	fmt.Println(locationRepr)
}

func (Location *Location) LS() {
	fmt.Printf(header)
	if Location.Bucket != "" {
		prefix := strings.Join(Location.Level, "/")
		if prefix != "/" {
			prefix = prefix + "/"
		}
		for name, _ := range Location.DirBuffer {
			fmt.Printf("%-50s %s\n", strings.TrimPrefix(name, prefix), "DIR")
		}
		s3s.ListFiles(
			&Location.S3,
			&Location.Bucket,
			&prefix,
		)
		return
	}
	// bucket level
	var names []string
	for name, _ := range Location.BucketsBuffer {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%-50s %s\n", name, Location.BucketsBuffer[name].CreationDate.Format(utils.DateFormat))
	}
}

func (Location *Location) File(filename string) {
	if Location.Bucket == "" {
		fmt.Println("Not in a bucket.")
		return
	}
	fullPath := strings.Join(Location.Level, "/") + "/" + filename
	err := s3s.DescribeObject(&Location.S3, &Location.Bucket, fullPath)
	if err != nil {
		return
	}
}

func (Location *Location) Reset() {
	Location.Level = []string{}
	Location.Bucket = ""
}

// updates the current level in case of expected changes
func (Location *Location) Update() {
	if Location.Bucket != "" {
		// dirs
		prefix := strings.Join(Location.Level, "/")
		if prefix != "/" {
			prefix = prefix + "/"
		}
		Location.BucketsBuffer = make(map[string]*s3.Bucket)
		dirBuff, err := s3s.ListDirs(
			&Location.S3,
			&Location.Bucket,
			&prefix,
		)
		if err != nil {
			fmt.Printf("Failed to fetch directories: %v\n", err)
			return
		}
		Location.DirBuffer = dirBuff

	} else {
		// buckets update, serves as the initial health check
		buff := make(map[string]*s3.Bucket)
		res, err := s3s.GetAllBuckets(&Location.S3)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, b := range res {
			buff[*b.Name] = b
		}
		Location.BucketsBuffer = buff
		Location.DirBuffer = make(map[string]struct{})
		Location.FileBuffer = make(map[string]struct{})
	}
}
