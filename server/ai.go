package server

import (
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

type AI struct {
	URL    string `json:"url" env:"url"`
	ApiKey string `json:"apiKey" env:"api_key"`
	Model  string `json:"model" env:"model"`
}

func getAiConfig(c *gin.Context) {
	aiConfig := AI{
		URL:    conf.Conf.AI.URL,
		ApiKey: conf.Conf.AI.ApiKey,
		Model:  conf.Conf.AI.Model,
	}
	common.SuccessResp(c, aiConfig)
}
