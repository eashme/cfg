package cfg

import "os"

func GetEnv(k string) string {
	return os.Getenv(StandCode(k))
}
