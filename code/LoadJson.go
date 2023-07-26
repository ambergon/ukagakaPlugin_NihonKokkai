package main

import (
    "fmt"
    "os"
    "encoding/json"
)

type KokkaiConfig struct {
    StartSec         int
    IntervalSec      int 
}
var Config KokkaiConfig

func LoadJson(){
	JsonFile, err := os.Open( Directory + "/Config.json")
	if err != nil {
        fmt.Println( err )
	}
	defer JsonFile.Close()
    decoder := json.NewDecoder( JsonFile )
    err     = decoder.Decode( &Config )
	if err != nil {
        fmt.Println(  err  )
    }
}

