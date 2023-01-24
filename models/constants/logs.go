package constants

import "github.com/rs/zerolog"

const (
	LogFileName      = "fileName"
	LogCorrelationId = "correlationId"
	LogUserId        = "userId"
	LogJobId         = "jobId"
	LogServerId      = "serverId"

	LogLevelFallback = zerolog.InfoLevel
)
