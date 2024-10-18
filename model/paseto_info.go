package model

type PasetoConfig struct {
	PubKey   string `mapstructure:"pub_key" json:"pub_key"`
	Implicit string `mapstructure:"implicit" json:"implicit"`
}
