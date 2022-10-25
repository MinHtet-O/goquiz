package cmd

import (
	"errors"
	"os"
	"strconv"
)

const (
	DB_Inmemory    = "inmemory"
	DB_Postgres    = "postgres"
	Transport_REST = "rest"
	Transport_GRPC = "grpc"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func GetenvStr(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func GetenvBool(key string, fallback bool) bool {
	s := GetenvStr(key, "")
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fallback
	}
	return v
}

func GetenvInt(key string, fallback int) int {
	int := GetenvStr(key, "")
	v, err := strconv.Atoi(int)
	if err != nil {
		return fallback
	}
	return v
}
