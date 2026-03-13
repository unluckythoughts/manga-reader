package main

import (
	"flag"
	"fmt"
	"os"
)

func readFlags() {
	var name string
	flag.StringVar(&name, "name", "book-reader", "Service name")
	if err := os.Setenv("SERVICE_NAME", name); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_NAME, err: %+v", err))
	}

	var enableDB bool
	flag.BoolVar(&enableDB, "enable-db", true, "Enable DB")
	if err := os.Setenv("SERVICE_ENABLE_DB", fmt.Sprintf("%t", enableDB)); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_ENABLE_DB, err: %+v", err))
	}

	var dbType string
	flag.StringVar(&dbType, "db-type", "sqlite", "DB Type")
	if err := os.Setenv("SERVICE_DB_TYPE", dbType); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_DB_TYPE, err: %+v", err))
	}

	var enableCache bool
	flag.BoolVar(&enableCache, "enable-cache", false, "Enable Cache")
	if err := os.Setenv("SERVICE_ENABLE_CACHE", fmt.Sprintf("%t", enableCache)); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_ENABLE_CACHE err: %+v", err))
	}

	var enableBus bool
	flag.BoolVar(&enableBus, "enable-bus", false, "Enable Bus")
	if err := os.Setenv("SERVICE_ENABLE_BUS", fmt.Sprintf("%t", enableBus)); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_ENABLE_BUS err: %+v", err))
	}

	var enableAuth bool
	flag.BoolVar(&enableAuth, "enable-auth", false, "Enable Auth")
	if err := os.Setenv("SERVICE_AUTH_ENABLE", fmt.Sprintf("%t", enableAuth)); err != nil {
		panic(fmt.Errorf("error while setting env var: SERVICE_AUTH_ENABLE err: %+v", err))
	}

	var dbFilePath string
	flag.StringVar(&dbFilePath, "db-file-path", "./data/book_reader.db", "DB File Path")
	if err := os.Setenv("DB_FILE_PATH", dbFilePath); err != nil {
		panic(fmt.Errorf("error while setting env var: DB_FILE_PATH err: %+v", err))
	}

	var dbDebug bool
	flag.BoolVar(&dbDebug, "db-debug", false, "DB Debug")
	if err := os.Setenv("DB_DEBUG", fmt.Sprintf("%t", dbDebug)); err != nil {
		panic(fmt.Errorf("error while setting env var: DB_DEBUG err: %+v", err))
	}

	var port int
	flag.IntVar(&port, "port", 8080, "DB Debug")
	if err := os.Setenv("WEB_PORT", fmt.Sprintf("%d", port)); err != nil {
		panic(fmt.Errorf("error while setting env var: WEB_PORT err: %+v", err))
	}

	var socketPath string
	flag.StringVar(&socketPath, "socket-path", "/socket", "Socket Path")
	if err := os.Setenv("WEB_SOCKET_PATH", socketPath); err != nil {
		panic(fmt.Errorf("error while setting env var: WEB_SOCKET_PATH err: %+v", err))
	}

	var workerCount int
	flag.IntVar(&workerCount, "worker-count", 20, "Worker Count")
	if err := os.Setenv("WEB_WORKER_COUNT", fmt.Sprintf("%d", workerCount)); err != nil {
		panic(fmt.Errorf("error while setting env var: WEB_WORKER_COUNT err: %+v", err))
	}

	var enableCors bool
	flag.BoolVar(&enableCors, "enable-cors", true, "Enable CORS")
	if err := os.Setenv("WEB_CORS", fmt.Sprintf("%t", enableCors)); err != nil {
		panic(fmt.Errorf("error while setting env var: WEB_CORS err: %+v", err))
	}

	var enableProxy bool
	flag.BoolVar(&enableProxy, "enable-proxy", false, "Enable Proxy")
	if err := os.Setenv("WEB_PROXY", fmt.Sprintf("%t", enableProxy)); err != nil {
		panic(fmt.Errorf("error while setting env var: WEB_PROXY err: %+v", err))
	}

	var logLevel string
	flag.StringVar(&logLevel, "log-level", "debug", "Log Level")
	if err := os.Setenv("LOG_LEVEL", logLevel); err != nil {
		panic(fmt.Errorf("error while setting env var: LOG_LEVEL err: %+v", err))
	}

	var sessionKey string
	flag.StringVar(&sessionKey, "session-key", "your-secret-key-change-this-in-production-32bytes", "Session Secret Key")
	if err := os.Setenv("SESSION_SECRET_KEY", sessionKey); err != nil {
		panic(fmt.Errorf("error while setting env var: SESSION_SECRET_KEY err: %+v", err))
	}

	flag.Parse()
}
