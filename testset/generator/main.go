package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

func readConf(defaults map[string]interface{}) *viper.Viper {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./")
	v.AutomaticEnv()
	v.SetConfigName("genconf")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(1)
	}
	return v
}

func AC(in string) eos.AccountName {
	return eos.AccountName(in)
}

func main() {
	fmt.Println("eosdaq contract generator")

	v := readConf(map[string]interface{}{
		"eosport": 18888,
	})

	api := eos.New(fmt.Sprintf("http://localhost:%d", v.GetInt("eosport")))

	//api.Debug = true
	//eos.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		"5HtZU5SArLEK3WDNntrK9fRCU8GFm9Ga4EAt9omGuYwiiFxMRyd",
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	api.SetSigner(keyBag)

	out, _ := api.GetTableRows(eos.GetTableRowsRequest{
		Scope: "eosdaq",
		Code:  "eosdaq",
		Table: "bidmatch",
		JSON:  true,
	})

	data, _ := json.Marshal(out)
	fmt.Println(string(data))
	/*
		actionResp, err := api.SignPushActions(

			eosdaq.SyncVerify(
				eos.AccountName("eosdaq"),
			),
		)
		if err != nil {
			fmt.Println("ERROR calling :", err)
		} else {
			fmt.Println("RESP:", actionResp)
		}
	*/

}
