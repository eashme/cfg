package coord_cfg

import "os"

func getFromEnv(k string) string {
	return os.Getenv(StandCode(k))
}
