package api

import "hq.0xa1.red/axdx/scheduler/internal/logging"

var logger *logging.Logger

func init() {
	logger = logging.MustNew()
}
