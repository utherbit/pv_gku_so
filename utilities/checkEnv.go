package utilities

import (
	"github.com/joho/godotenv"
	"os"
)

func CheckEnvFile() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func LookupEnv(out *string, key string, defaultVal ...string) {
	val, exist := os.LookupEnv(key)
	if !exist {
		if len(defaultVal) > 0 {
			*out = defaultVal[0]
		} else {
			panic(key + " not found in .env file")
		}
	}
	*out = val
}
