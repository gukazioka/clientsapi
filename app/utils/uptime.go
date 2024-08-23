package utils

import (
	"time"
)

var serverStartTime time.Time

func Initialize() {
	serverStartTime = time.Now()
}

func GetUptime() string {
	return serverStartTime.String()
}
