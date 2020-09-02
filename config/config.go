package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type app struct {
	Name     string `yaml:"name"`
	Auth     auth
	Server   server
	Upload   upload
	Log      sysLog
	Database database
	Redis    redis
}

type auth struct {
	Key     string `yaml:"key"`
	Expires int    `yaml:"expires"`
}

type server struct {
	Port         int           `yaml:"port"`
	Mode         string        `yaml:"mode"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type image struct {
	Path            string `yaml:"path"`
	MaxSize         int    `yaml:"maxSize"`
	AllowExtensions string `yaml:"allowExtensions"`
}

type upload struct {
	Path  string `yaml:"path"`
	Image image
}

type sysLog struct {
	Path       string `yaml:"path"`
	Extension  string `yaml:"extension"`
	DateFormat string `yaml:"dateFormat"`
}

type database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Prefix   string `yaml:"prefix"`
	Charset  string `yaml:"charset"`
	Collate  string `yaml:"collate"`
}

type redis struct {
	Host        string `yaml:"app.redis.host"`
	Port        int    `yaml:"app.redis.port"`
	Password    string `yaml:"app.redis.password"`
	MaxIdle     int    `yaml:"app.redis.maxIdle"`
	MaxActive   int    `yaml:"app.redis.maxActive"`
	IdleTimeout int    `yaml:"app.redis.idleTimeout"`
}

var App = &app{}

const conf string = "E:/go-gin/config/config.yml"

func Init() {
	fileContent, readFileErr := ioutil.ReadFile(conf)
	if readFileErr != nil {
		log.Fatalf("config.Init, failed to read '%s': %v", conf, readFileErr)
	}

	decodeError := yaml.Unmarshal(fileContent, App)
	if decodeError != nil {
		log.Fatalf("config.Init, yaml.Unmarshal failed: %v", decodeError)
	}
}
