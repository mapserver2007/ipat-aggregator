package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
)

type SLF4JFormatter struct{}

func (f *SLF4JFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	level := entry.Level.String()
	level = fmt.Sprintf("%-5s", level) // SLF4J形式に合わせてレベルを整列

	var buf bytes.Buffer
	goroutineId := getGoroutineID()
	buf.WriteString(fmt.Sprintf("%s [%s] thread%d - %s\n", timestamp, level, goroutineId, entry.Message))

	return buf.Bytes(), nil
}

func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, err := strconv.Atoi(string(idField))
	if err != nil {
		return -1
	}
	return id
}
