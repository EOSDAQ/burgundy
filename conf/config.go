package conf

import (
	"burgundy/util"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	ConfEOSManageContract  string
	ConfEOSCrawlExclude    string
	ConfEOSCrawlDurationMS int
	ConfEOSBaseSymbol      string

	ConfDBHOST string
	ConfDBPORT int
	ConfDBUSER string
	ConfDBPASS string
	ConfDBNAME string

	ConfAWSRegion string
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
	ConfEOSAcctContract:    "eosdaq555555",
	ConfEOSManageContract:  "eosdaqmanage",
	ConfEOSCrawlExclude:    "",
	ConfEOSCrawlDurationMS: 500,
	ConfEOSBaseSymbol:      "SYS",
	ConfDBHOST:             "www.db4free.net",
	ConfDBPORT:             3306,
	ConfDBUSER:             "eosdaquser",
	ConfDBPASS:             "eosdaqvotmdnjem",
	ConfDBNAME:             "eosdaq",
	ConfAWSRegion:          "ap-northeast-2",
}

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
	ssmsvc      *ssm.SSM
	cacheString map[string]string
	cacheInt    map[string]int
}

// Burgundy ...
var Burgundy *ViperConfig

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

	defaultExclude := []string{}
	for i := 0; i < 22; i++ {
		base := util.ConvertBase(i, 6)
		contract := strings.Replace(fmt.Sprintf("eosdaq%06s", base), "0", "o", -1)
		if contract == "eosdaqoooo2o" {
			continue
		}
		defaultExclude = append(defaultExclude, contract)
	}

	var err error
	Burgundy, err = readConfig(map[string]interface{}{
		"port":               defaultConf.ConfServerPORT,
		"timeout":            defaultConf.ConfServerTIMEOUT,
		"logmode":            defaultConf.ConfServerLOGMODE,
		"loglevel":           defaultConf.ConfAPILOGLEVEL,
		"profile":            false,
		"profilePort":        6060,
		"eos_host":           defaultConf.ConfEOSHOST,
		"eos_port":           defaultConf.ConfEOSPORT,
		"eos_acctcontract":   defaultConf.ConfEOSAcctContract,
		"eos_managecontract": defaultConf.ConfEOSManageContract,
		"eos_crawlexclude":   strings.Join(defaultExclude, ","),
		"eos_crawlMS":        defaultConf.ConfEOSCrawlDurationMS,
		"eos_baseSymbol":     defaultConf.ConfEOSBaseSymbol,
		"db_host":            defaultConf.ConfDBHOST,
		"db_port":            defaultConf.ConfDBPORT,
		"db_user":            defaultConf.ConfDBUSER,
		"db_pass":            defaultConf.ConfDBPASS,
		"db_name":            defaultConf.ConfDBNAME,
		"aws_region":         defaultConf.ConfAWSRegion,
		"env":                "production",
	})
	if err != nil {
		fmt.Printf("Error when reading config: %v\n", err)
		os.Exit(1)
	}
	err = Burgundy.InitAWSSSM()
	if err != nil {
		fmt.Printf("Error when Init AWS SSM: %v\n", err)
		os.Exit(1)
	}

	Burgundy.BindPFlags(pflag.CommandLine)
}

func readConfig(defaults map[string]interface{}) (*ViperConfig, error) {
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

	v.AutomaticEnv()

	switch strings.ToUpper(v.GetString("ENV")) {
	case "DEVEL":
		fmt.Println("Loading Development Environment...")
		v.SetConfigName(defaultConf.EnvServerDEV)
		v.Debug()
	case "STAGE":
		fmt.Println("Loading Stage Environment...")
		v.SetConfigName(defaultConf.EnvServerSTAGE)
	case "PROD":
		fmt.Println("Loading Production Environment...")
		v.SetConfigName(defaultConf.EnvServerPROD)
	default:
		fmt.Println("Loading Production(Default) Environment...")
		v.SetConfigName(defaultConf.EnvServerPROD)
	}

	err := v.ReadInConfig()
	switch err.(type) {
	default:
		return &ViperConfig{}, err
	case viper.ConfigFileNotFoundError:
		fmt.Printf("Warn: %s\n", err)
	}

	return &ViperConfig{
		v,
		nil,
		make(map[string]string),
		make(map[string]int),
	}, nil
}

func (vp *ViperConfig) InitAWSSSM() (err error) {
	region := vp.Viper.GetString("aws_region")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	vp.ssmsvc = ssm.New(sess, aws.NewConfig().WithRegion(region))
	return nil
}

func (vp *ViperConfig) GetString(key string) string {
	if v, ok := vp.cacheString[key]; ok {
		return v
	}
	keyname := fmt.Sprintf("/eosdaq/%s/%s", vp.Viper.GetString("ENV"), key)
	withDecryption := true
	param, err := vp.ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		fmt.Printf("GetString cannot get parameter keyname[%s] err[%s]\n", keyname, err)
		vp.cacheString[key] = vp.Viper.GetString(key)
	} else {
		vp.cacheString[key] = *param.Parameter.Value
	}
	return vp.cacheString[key]
}

func (vp *ViperConfig) GetInt(key string) int {
	if v, ok := vp.cacheInt[key]; ok {
		return v
	}
	keyname := fmt.Sprintf("/eosdaq/%s/%s", vp.Viper.GetString("ENV"), key)
	withDecryption := true
	param, err := vp.ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		fmt.Printf("GetInt cannot get parameter keyname[%s] err[%s]\n", keyname, err)
		vp.cacheInt[key] = vp.Viper.GetInt(key)
	} else {
		v, err := strconv.Atoi(*param.Parameter.Value)
		if err != nil {
			fmt.Printf("GetInt parse error keyname[%s] param[%s] err[%s]\n", keyname, param, err)
			vp.cacheInt[key] = vp.Viper.GetInt(key)
		} else {
			vp.cacheInt[key] = v
		}
	}
	return vp.cacheInt[key]
}

// APILogLevel string to log level
func (vp *ViperConfig) APILogLevel() log.Lvl {
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
func (vp *ViperConfig) SetProfile() {
	if vp.GetBool("profile") {
		runtime.SetBlockProfileRate(1)
		go func() {
			profileListen := fmt.Sprintf("0.0.0.0:%d", vp.GetInt("profilePort"))
			http.ListenAndServe(profileListen, nil)
		}()
	}
}
