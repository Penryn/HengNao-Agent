package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env string

	Hertz Hertz `yaml:"hertz"`
	MySQL MySQL `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
	Api   Api   `yaml:"api"`
}

type MySQL struct {
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	DBName   string `yaml:"db_name"`

	DSN string `yaml:"dsn"`
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns"`
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns"`
	// ConnMaxLifetime 连接最大存活时间
	ConnMaxLifetime int `yaml:"conn_max_lifetime"` // 秒
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	DB       int    `yaml:"db"`
}

type Hertz struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
}

type Api struct {
	Url                   string `yaml:"url"`
	Key                   string `yaml:"key"`
	Secret                string `yaml:"secret"`
	GetRelevantHighlights string `yaml:"get_relevant_highlights"`
	GetMeetingMinutes     string `yaml:"get_meeting_minutes"`
	ChatMeeting           string `yaml:"chat_meeting"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}

	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		hlog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.Error("validate config error - %v", err)
		panic(err)
	}

	conf.Env = GetEnv()

	pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() hlog.Level {
	level := GetConf().Hertz.LogLevel
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
