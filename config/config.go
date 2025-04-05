package config

import (
	"os"
	"user-service/common/util"

	"github.com/sirupsen/logrus"
)

var Config AppConfig

type AppConfig struct {
	Port                   int      `json:"port"`
	AppName                string   `json:"appName"`
	AppEnv                 string   `json:"appEnv"`
	SignatureKey           string   `json:"signatureKey"`
	Database               Database `json:"database"`
	RateLimiterMaxRequests int      `json:"rateLimiterMaxRequests"`
	RateLimiterTimeSeconds int      `json:"rateLimiterTimeSeconds"`
	JwtSecretKey           string   `json:"jwtSecretKey"`
	JwtExpirationTime      int      `json:"jwtExpirationTime"`
}

type Database struct {
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	Name                   string `json:"name"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	MaxOpenConnections     int    `json:"maxOpenConnections"`
	MaxIdleConnections     int    `json:"maxIdleConnections"`
	MaxLifetimeConnections int    `json:"maxLifetimeConnections"`
	MaxIdleTime            int    `json:"maxIdleTime"`
}

func init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config: %v", err)
		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
