package infrastructure

import (
	"middleware/internal/config"

	http_client_infratructure "middleware/infrastructure/http/http_client"
	"middleware/infrastructure/logger/logrus"
)

type Infrastructure struct {
	Config     *config.Config
	Logger     logrus.LogrusInterface
	HttpClient http_client_infratructure.HttpClientInterface
}
