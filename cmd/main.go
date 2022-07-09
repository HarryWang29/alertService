package main

import (
	"github.com/HarryWang29/alertService/config"
	"github.com/HarryWang29/alertService/job"
	"os"
)

func main() {
	//加载配置文件
	configPath := ""
	if os.Getenv("alert_service") == "dev" {
		configPath = "./config_dev.yaml"
	} else {
		configPath = "./config.yaml"
	}

	//加载配置文件
	conf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}
	err = config.ParseProvinceCodeFromFile("ProvinceCode.json")
	if err != nil {
		panic(err)
	}
	myApp := config.NewApp()
	err = myApp.Parse(conf)
	if err != nil {
		panic(err)
	}
	go func() {
		job.Start(myApp.Jobs)
	}()
	select {}
}
