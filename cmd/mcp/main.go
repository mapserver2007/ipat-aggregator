package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/mapserver2007/ipat-aggregator/di"
	"github.com/sirupsen/logrus"
)

const (
	logFilePath = "/tmp/mcp.log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&SLF4JFormatter{})

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	logger.SetOutput(io.MultiWriter(logFile))

	mcpServer := di.NewMcpServer(logger)

	if err := mcpServer.Load(ctx); err != nil {
		logger.Fatalf("Failed to load master data: %v", err)
		return
	}

	if err := mcpServer.Start(); err != nil {
		logger.Fatalf("Failed to start MCP server: %v", err)
	}
}

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
