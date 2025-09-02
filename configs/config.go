package configs

import (
	"fmt"
	rlog "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

var (
	globalConfig GlobalConfig
	once         sync.Once
)

type GlobalConfig struct {
	AppConfig        AppConf        `yaml:"app" mapstructure:"app"`
	LogConfig        LogConf        `yaml:"log" mapstructure:"log"`
	DbConfig         DbConf         `yaml:"db" mapstructure:"db"`
	ScheduleConfig   ScheduleConf   `yaml:"scheduler" mapstructure:"scheduler"`
	CredentialConfig CredentialConf `yaml:"credential" mapstructure:"credential"`
	AlphaConfig      AlphaConf      `yaml:"alpha" mapstructure:"alpha"`
	ResultConfig     ResultConf     `yaml:"result" mapstructure:"result"`
}

type ResultConf struct {
	NeedStoreRecords bool  `yaml:"need_store_records" mapstructure:"need_store_records"`
	MaxRetryNum      int64 `yaml:"max_retry_num" mapstructure:"max_retry_num"`
}
type AlphaConf struct {
	RetryNum       int64 `yaml:"retry_num" mapstructure:"retry_num"`
	ChannelLen     int64 `yaml:"channel_len" mapstructure:"channel_len"`
	ScanIdeaSecond int64 `yaml:"scan_idea_second" mapstructure:"scan_idea_second"`
}
type CredentialConf struct {
	UserName string `yaml:"user_name" mapstructure:"user_name"`
	Password string `yaml:"password" mapstructure:"password"`
	Token    string `yaml:"token" mapstructure:"token"`
}
type AppConf struct {
	AppName     string `yaml:"app_name" mapstructure:"app_name"`
	Version     string `yaml:"version" mapstructure:"version"`
	Port        int    `yaml:"port" mapstructure:"port"`
	RunMod      string `yaml:"run_mod"  mapstructure:"run_mod"`
	Concurrency int64  `yaml:"concurrency" mapstructure:"concurrency"`
}

type LogConf struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"`
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`
	SaveDays   uint   `yaml:"save_days" mapstructure:"save_days"`
	Level      string `yaml:"level" mapstructure:"level"`
}

type DbConf struct {
	Host                     string `yaml:"host" mapstructure:"host"`
	Port                     int    `yaml:"port" mapstructure:"port"`
	User                     string `yaml:"user" mapstructure:"user"`
	Password                 string `yaml:"password" mapstructure:"password"`
	Dbname                   string `yaml:"dbname" mapstructure:"db_name"`
	MaxIdleConn              int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"`
	MaxOpenConn              int    `yaml:"max_open_conn" mapstructure:"max_open_conn"`
	MaxIdleTime              int    `yaml:"max_idle_time" mapstructure:"max_idle_time"`
	SlowThresholdMillisecond int    `yaml:"max_idle_time" mapstructure:"slow_threshold_millisecond"`
}

type ScheduleConf struct {
	TimeSecond int `yaml:"time_second" mapstructure:"time_second"`
}

func GetGlobalConfig() *GlobalConfig {
	once.Do(readConf)
	return &globalConfig
}

func readConf() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	//从程序工作目录开始
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		panic("read config error " + err.Error())
	}
	err = viper.Unmarshal(&globalConfig)

	if err == nil {
		fmt.Println("read config success")
	}

}

func InitGlobalConfig() {
	config := GetGlobalConfig()
	level, err := log.ParseLevel(config.LogConfig.Level)
	if err != nil {
		panic("parse log level error " + err.Error())
	}
	log.SetFormatter(&logFormatter{
		TextFormatter: log.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	})

	log.SetReportCaller(true)
	log.SetLevel(level)
	switch globalConfig.LogConfig.LogPattern {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		logger, err := rlog.New(
			config.LogConfig.LogPath+".%Y%m%d",
			rlog.WithRotationTime(time.Hour*24),
			rlog.WithRotationCount(config.LogConfig.SaveDays))
		if err != nil {
			panic("log conf err" + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log init err")
	}

}
