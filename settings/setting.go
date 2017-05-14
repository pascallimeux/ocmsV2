package settings

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/op/go-logging"
)

type Settings struct {
	Version            string
	HttpHostUrl        string
	LogFileName	   string
	LogMode            string
	LogFile		   *os.File
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration

	ChainCodePath      string
	ChainCodeVersion   string
	ChainCodeID	   string
	ChainID		   string
	ProviderName       string
	Repo               string
	StatstorePath      string
	NetworkConfigfile  string
	ChannelConfigFile  string
	Adminusername	   string
	AdminPwd           string


}
var log = logging.MustGetLogger("ocms.settings")

func (s *Settings) ToString() string {
	st :=     "Logger          --> file:" + s.LogFileName + " in " + s.LogMode + " mode \n"
	st = st + "Server          --> url :" + s.HttpHostUrl
	return st
}

func (s *Settings) CloseLogger() {
	s.LogFile.Close()
}

func (s *Settings) InitLogger() (err error){
	s.LogFile, err = InitLogger(s.LogMode, s.LogFileName)
	if err != nil {
		return errors.New("Error logfile!")
	}
	return nil
}


func findConfigFile(configPath, configFileName string) error {
	path := configPath + "/" + configFileName + ".toml"
	if _, err := os.Stat(path); err != nil {
		configPath = os.Getenv("OCMSPATH")
		if configPath == "" {
			return errors.New("no config file found!")
		} else {
			fmt.Println("read config file: " + configPath + "/" + configFileName + ".toml")
			viper.SetConfigName(configFileName)
			viper.AddConfigPath(configPath)
			return nil
		}
	} else {
		fmt.Println("read config file: ", path)
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(configPath)
		return nil
	}
}

func GetSettings(configPath, configFileName string) (Settings, error) {
	var configuration Settings
	err := findConfigFile(configPath, configFileName)
	if err != nil {
		fmt.Println(err.Error())
		return configuration, errors.New("Config file not found...")
	}
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return configuration, errors.New("Config file not found...")
	} else {
		configuration.Version = viper.GetString("OCMSversion.version")

		logMode := viper.GetString("logger.mode")
		logFileName := os.Getenv("OCMSLOGFILE")
		if logFileName == "" {
			logFileName = viper.GetString("logger.logFileName")
		}
		configuration.LogFileName = logFileName
		configuration.LogMode = logMode

		configuration.LogFile, err = InitLogger(logMode, logFileName)
		if err != nil {
			return configuration, errors.New("Error logfile!")
		}


		configuration.HttpHostUrl, err = getHostUrl()
		if err != nil {
			return configuration, err
		}
		configuration.ReadTimeout = viper.GetDuration("server.readTimeout")
		configuration.WriteTimeout = viper.GetDuration("server.writeTimeout")

		configuration.ChainCodePath = viper.GetString("chaincode.chainCodePath")
		configuration.ChainCodeVersion = viper.GetString("chaincode.chainCodeVersion")
		configuration.ChainCodeID = viper.GetString("chaincode.chainCodeID")
		configuration.ChainID = viper.GetString("chaincode.chainID")
		configuration.ProviderName = viper.GetString("chaincode.providerName")

		configuration.Repo = viper.GetString("path.repo")
		configuration.StatstorePath = viper.GetString("path.statStorePath")
		configuration.NetworkConfigfile = viper.GetString("path.networkConfigFile")
		configuration.ChannelConfigFile = viper.GetString("path.channelConfigFile")

		configuration.Adminusername = viper.GetString("admin.adminUsername")
		configuration.AdminPwd = viper.GetString("admin.adminPwd")

		fmt.Println("Application configuration: \n" + configuration.ToString())
		return configuration, nil
	}
}

func getHostUrl() (string, error) {
	ipAddress := viper.GetString("server.httpHostIp")
	ipPort := viper.GetInt("server.httpHostPort")

	var err error
	if ipAddress == "" {
		ipAddress, err = getOutboundIP()
		if err != nil {
			return ipAddress, errors.New(" Error to get local IP address")
		}
	}
	ipAddress = ipAddress + ":" + strconv.Itoa(ipPort)
	return ipAddress, nil
}

func getOutboundIP() (string, error) {
	var localAddr string

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return localAddr, err
	}
	defer conn.Close()

	localAddr = conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx], nil
}

func InitLogger2(logMode, logFilePath string) (*os.File, error) {
	format := logging.MustStringFormatter(
		//`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		` %{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
	)
	f := os.Stderr
	if logFilePath != "" {
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			f, err = os.Create(logFilePath)
			if err != nil {
				return f, err
			}
		}else {
			f, err = os.OpenFile(logFilePath, os.O_APPEND | os.O_WRONLY, 0600)
			if err != nil {
				return f, err
			}
		}
	}

	backend := logging.NewLogBackend(f, "ocms", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	level := logging.ERROR
	switch logMode {
	case "critical":
		level = logging.CRITICAL
	case "error":
		level = logging.ERROR
	case "warning":
		level = logging.WARNING
	case "info":
		level = logging.INFO
	case "debug":
		level = logging.DEBUG
	}
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
	log.Debug("Logger initialized")
	return f, nil
}


func  InitLogger(logMode, logFilePath string) (*os.File, error) {
	f := os.Stderr
	if logFilePath != "" {
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			f, err = os.Create(logFilePath)
			if err != nil {
				return f, err
			}
		}else {
			f, err = os.OpenFile(logFilePath, os.O_APPEND | os.O_WRONLY, 0600)
			if err != nil {
				return f, err
			}
		}
	}
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} [%{module}] %{level:.4s} : %{color:reset} %{message}`,
	)
	logLevel := logging.DEBUG
	backend := logging.NewLogBackend(f, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(logLevel), "ocms")
	return f, nil
}

