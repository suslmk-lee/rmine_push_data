package common

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfigProperties map[string]string

var ConfInfo AppConfigProperties

// init initializes the configuration
func init() {
	path, _ := os.Getwd()
	println(path)
	_, err := ReadPropertiesFile("config.properties")
	if err != nil {
		path, _ := os.Getwd()
		println(path)
		return
	}
}

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	ConfInfo = AppConfigProperties{}

	if len(filename) == 0 {
		return ConfInfo, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				ConfInfo[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return ConfInfo, nil
}

func RandomString(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

const timestampFile = "last_checked_timestamp.txt"

// SaveLastCheckedTime saves the last checked time to a file
func SaveLastCheckedTime(t time.Time) error {
	millis := t.UnixMilli() / int64(time.Millisecond)
	return ioutil.WriteFile(timestampFile, []byte(fmt.Sprintf("%d", millis)), 0644)
}

func LoadLastCheckedTime() (time.Time, error) {
	data, err := ioutil.ReadFile(timestampFile)
	if err != nil {
		return time.Time{}, err
	}

	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, millis*int64(time.Millisecond)), nil
}
