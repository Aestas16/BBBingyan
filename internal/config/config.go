package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Configuration struct {
    Server  ServerConfig    `yaml:"server"`
    Redis   RedisConfig     `yaml:"redis"`
}
type ServerConfig struct {
    Port    int         `yaml:"port"`
    Version string      `yaml:"version"`
    JwtKey  string      `yaml:"jwtkey"`
    Admin   AdminConfig `yaml:"admin"`
    Email   EmailConfig `yaml:"email"`
    Dsn     string      `yaml:"dsn"`
}
type RedisConfig struct {
    Host     string `yaml:"host"`
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
}
type AdminConfig struct {
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
}
type EmailConfig struct {
    Host        string  `yaml:"host"`
    Port        int     `yaml:"port"`
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
    Expire      int     `yaml:"expire"`
}

var Config Configuration

func InitConfig() {
    yamlFile, err := os.ReadFile("config.yaml")
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal(yamlFile, &Config)
    if err != nil {
        panic(err)
    }
}