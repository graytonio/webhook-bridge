package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Webhook struct {
    Path string `mapstructure:"path"`
    Method string `mapstructure:"method"`
    URL string `mapstructure:"url"`
    Headers map[string]string `mapstructure:"headers"`
}

type Config struct {
    LogLevel string `mapstructure:"log_level"`
    ListenAddress string `mapstructure:"listen_address"`
    Webhooks []Webhook `mapstructure:"webhooks"`
}

func forwardWebhook(config *Webhook) gin.HandlerFunc {
   return func(c *gin.Context) {
    req, err := http.NewRequest(strings.ToUpper(config.Method), config.URL, nil)
    if err != nil {
        c.AbortWithStatusJSON(500, map[string]interface{}{"error": err.Error()})
        return
    }

    for key, value := range config.Headers {
        req.Header.Add(key, value)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        c.AbortWithStatusJSON(500, map[string]interface{}{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.AbortWithStatusJSON(500, map[string]interface{}{"error": err.Error()})
        return
    }

    c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
   } 
}

var config Config

func init() {
    viper.SetConfigName("config")
    viper.AddConfigPath("/config/")
    viper.AddConfigPath(".")

    viper.BindEnv("LOG_LEVEL")
    viper.SetDefault("log_level", "INFO")

    viper.SetDefault("listen_address", "0.0.0.0:8080")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Fatal("Unable to find config file.") 
        } else {
            log.WithField("error", err.Error()).Fatalf("error reading config")
        }
    } 

    if err := viper.Unmarshal(&config); err != nil {
        log.WithField("error", err.Error()).Fatal("error unmarshalling config")
    }

    level, err := log.ParseLevel(config.LogLevel)
    if err != nil {
        log.WithField("log_level_input", config.LogLevel).Error("Could not parse log level. Using info by default")
        level = log.InfoLevel
    }

    log.SetLevel(level)

    if level == log.DebugLevel {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }

    log.Debugf("Loaded Config: %+v", config)
}

func main() {
    r := gin.Default()

    for _, webhook := range config.Webhooks {
        r.POST(webhook.Path, forwardWebhook(&webhook))
    }

    r.Run(config.ListenAddress)
}
