package main

import (
	"io/ioutil"
	"os"
	"time"
)

func readFileAsString(path string) string {
	out, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return string(out)
}

func writeFileFromString(path string, content string) {
	e := ioutil.WriteFile(path, []byte(content), 0644)
	if e != nil {
		panic(e)
	}
}
func existAndIsFile(path string) bool {
	stat, e := os.Stat(path)
	if os.IsNotExist(e) {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

func truncateTime(input time.Time) time.Time {
	return time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, input.Location())
}
