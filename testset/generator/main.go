package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	eos "github.com/eoscanada/eos-go"
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
		"chainID": "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f",
	})

	cid, _ := hex.DecodeString(v.GetString("chainID"))
	api := eos.New(fmt.Sprintf("http://10.168.0.103:%d", v.GetInt("eosport")), cid)

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
		Table: "tx",
		JSON:  true,
	})

	fmt.Printf("out value [%v]\n", out)
	/*
		var daqRes []*eosdaq.EosdaqTx
		out.BinaryToStructs(&daqRes)
		fmt.Printf("res value [%v]\n", daqRes)
	*/
	data, _ := json.Marshal(out)
	fmt.Printf("data value [%s]\n", string(data))

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
