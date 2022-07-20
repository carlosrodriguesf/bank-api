package model

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
)

type (
	Options struct {
		Logger logger.Logger
	}

	Response struct {
		Data  interface{}        `json:"data,omitempty"`
		Error *apierror.ApiError `json:"error,omitempty"`
	}
)
