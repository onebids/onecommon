package kvconfig

import (
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
)

type CommonConfig struct {
	Env   string
	Kitex Kitex `yaml:"kitex"`
	MySQL MySQL `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kitex struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	MetricsPort     string `yaml:"metrics_port"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
}

func GetCommonConfig(registryAddr string) (*CommonConfig, error) {
	client, err := api.NewClient(&api.Config{Address: registryAddr})
	if err != nil {
		fmt.Println("Error creating Consul client:", err)
		return nil, err
	}
	//获取配置
	content, _, err := client.KV().Get("onebids/common", nil)
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	conf := new(CommonConfig)
	err = yaml.Unmarshal(content.Value, &conf)
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}

	return conf, nil
}
