package service

import (
	"github.com/speedland/go/gae/service/asynctask"
)

type asynctaskEndpoint struct {
	path   string
	config *asynctask.Config
}
