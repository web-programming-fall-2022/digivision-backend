package jobs

import (
	"github.com/arimanius/digivision-backend/internal/bootstrap/job"
	"github.com/arimanius/digivision-backend/internal/cfg"
)

func StartJobs(config cfg.Config) []job.WithGracefulShutDown {
	// TODO: Instantiate your job here and add those that need graceful shutdown to the return value.
	return []job.WithGracefulShutDown{}
}
