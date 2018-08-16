package conf

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DefaultConf ...
type DefaultConf struct {
	EnvServerDEV   string
	EnvServerSTAGE string
	EnvServerPROD  string

	ConfServerPORT    int
	ConfServerLOGMODE string
	ConfServerTIMEOUT int
	ConfAPILOGLEVEL   string

	ConfEOSHOST            string
	ConfEOSPORT            int
	ConfEOSAcctContract    string
	ConfEOSCrawlContract   string
	ConfEOSCrawlDurationMS int
	ConfEOSBaseSymbol      string

	ConfDBHOST string
	ConfDBPORT int
	ConfDBUSER string
	ConfDBPASS string
	ConfDBNAME string
}

var defaultConf = DefaultConf{
	EnvServerDEV:           ".env.dev",
	EnvServerSTAGE:         ".env.stage",
	EnvServerPROD:          ".env",
	ConfServerPORT:         2333,
	ConfServerLOGMODE:      "console",
	ConfServerTIMEOUT:      30,
	ConfAPILOGLEVEL:        "debug",
	ConfEOSHOST:            "http://10.100.100.2",
	ConfEOSPORT:            18888,
	ConfEOSAcctContract:    "dollarbillgo",
	ConfEOSCrawlContract:   "eosseieossei",
	ConfEOSCrawlDurationMS: 500,
	ConfEOSBaseSymbol:      "SYS",
	ConfDBHOST:             "www.db4free.net",
	ConfDBPORT:             3306,
	ConfDBUSER:             "eosdaquser",
	ConfDBPASS:             "eosdaqvotmdnjem",
	ConfDBNAME:             "eosdaq",
}

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

// Burgundy ...
var Burgundy ViperConfig

func init() {
	pflag.BoolP("version", "v", false, "Show version number and quit")
	pflag.IntP("port", "p", defaultConf.ConfServerPORT, "burgundy Port")
	pflag.IntP("timeout", "t", defaultConf.ConfServerTIMEOUT, "burgundy Context timeout(sec)")

	pflag.String("db_host", defaultConf.ConfDBHOST, "burgundy's DB host")
	pflag.Int("db_port", defaultConf.ConfDBPORT, "burgundy's DB port")
	pflag.String("db_user", defaultConf.ConfDBUSER, "burgundy's DB user")
	pflag.String("db_pass", defaultConf.ConfDBPASS, "burgundy's DB password")
	pflag.String("db_name", defaultConf.ConfDBNAME, "burgundy's DB name")

	pflag.Parse()

	var err error
	Burgundy, err = readConfig(map[string]interface{}{
		"port":              defaultConf.ConfServerPORT,
		"timeout":           defaultConf.ConfServerTIMEOUT,
		"logmode":           defaultConf.ConfServerLOGMODE,
		"loglevel":          defaultConf.ConfAPILOGLEVEL,
		"profile":           false,
		"profilePort":       6060,
		"eos_host":          defaultConf.ConfEOSHOST,
		"eos_port":          defaultConf.ConfEOSPORT,
		"eos_acctcontract":  defaultConf.ConfEOSAcctContract,
		"eos_crawlcontract": defaultConf.ConfEOSCrawlContract,
		"eos_crawlMS":       defaultConf.ConfEOSCrawlDurationMS,
		"eos_baseSymbol":    defaultConf.ConfEOSBaseSymbol,
		"db_host":           defaultConf.ConfDBHOST,
		"db_port":           defaultConf.ConfDBPORT,
		"db_user":           defaultConf.ConfDBUSER,
		"db_pass":           defaultConf.ConfDBPASS,
		"db_name":           defaultConf.ConfDBNAME,
	})
	if err != nil {
		fmt.Printf("Error when reading config: %v\n", err)
		os.Exit(1)
	}

	Burgundy.BindPFlags(pflag.CommandLine)
	Burgundy.Debug()
}

func readConfig(defaults map[string]interface{}) (ViperConfig, error) {
	// Read Sequence (will overloading)
	// defaults -> config file -> env -> cmd flag
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./")
	v.AddConfigPath("./conf")
	v.AddConfigPath("../conf")
	v.AddConfigPath("../../conf")
	v.AddConfigPath("$HOME/.burgundy")

	v.SetEnvPrefix("eosdaq")
	v.AutomaticEnv()

	switch strings.ToUpper(v.GetString("ENV")) {
	case "DEVELOPMENT":
		fmt.Println("Loading Development Environment...")
		v.SetConfigName(defaultConf.EnvServerDEV)
	case "STAGE":
		fmt.Println("Loading Stage Environment...")
		v.SetConfigName(defaultConf.EnvServerSTAGE)
	case "PRODUCTION":
		fmt.Println("Loading Production Environment...")
		v.SetConfigName(defaultConf.EnvServerPROD)
	default:
		fmt.Println("Loading Production(Default) Environment...")
		v.SetConfigName(defaultConf.EnvServerPROD)
	}

	err := v.ReadInConfig()
	if err != nil {
		return ViperConfig{}, err
	}

	return ViperConfig{v}, nil
}

// APILogLevel string to log level
func (vp ViperConfig) APILogLevel() log.Lvl {
	switch strings.ToLower(vp.GetString("loglevel")) {
	case "off":
		return log.OFF
	case "error":
		return log.ERROR
	case "warn", "warning":
		return log.WARN
	case "info":
		return log.INFO
	case "debug":
		return log.DEBUG
	default:
		return log.DEBUG
	}
}

// SetProfile ...
func (vp ViperConfig) SetProfile() {
	if vp.GetBool("profile") {
		runtime.SetBlockProfileRate(1)
		go func() {
			profileListen := fmt.Sprintf("0.0.0.0:%d", vp.GetInt("profilePort"))
			http.ListenAndServe(profileListen, nil)
		}()
	}
}
