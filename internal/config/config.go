package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Configuration struct {
    Server  ServerConfig    `yaml:"server"`
}
type ServerConfig struct {
    Port    int         `yaml:"port"`
    Version string      `yaml:"version"`
    JwtKey  string      `yaml:"jwtkey"`
    Admin   AdminConfig `yaml:"admin"`
    Email   EmailConfig `yaml:"email"`
    Dsn     string      `yaml:"dsn"`
}
type AdminConfig struct {
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
}
type EmailConfig struct {
    Interval    int64   `yaml:"interval"`
    Host        string  `yaml:"host"`
    Port        int     `yaml:"port"`
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
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