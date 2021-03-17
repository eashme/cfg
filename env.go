package coord_cfg

import "os"

func GetEnv(k string) string {
	return os.Getenv(StandCode(k))
}
