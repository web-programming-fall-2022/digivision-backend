package jobs

import (
	"github.com/web-programming-fall-2022/digivision-backend/internal/bootstrap/job"
	"github.com/web-programming-fall-2022/digivision-backend/internal/cfg"
)

func StartJobs(config cfg.Config) []job.WithGracefulShutdown {
	// TODO: Instantiate your job here and add those that need graceful shutdown to the return value.
	return []job.WithGracefulShutdown{}
}
