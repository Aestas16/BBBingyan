package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Configuration struct {
    Server  ServerConfig    `yaml:"server"`
    SQL     SQLConfig       `yaml:"postgresql"`
}
type ServerConfig struct {
    Port    int         `yaml:"port"`
    Version string      `yaml:"version"`
    JwtKey  string      `yaml:"jwtkey"`
    Admin   AdminConfig `yaml:"admin"`
    Email   EmailConfig `yaml:"email"`
}
type AdminConfig struct {
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
}
type EmailConfig struct {
    Interval    string  `yaml:"interval"`
    Host        string  `yaml:"host"`
    Port        string  `yaml:"port"`
    Username    string  `yaml:"username"`
    Password    string  `yaml:"password"`
}
type SQLConfig struct {
    User        string  `yaml:"user"`
    Password    string  `yaml:"password"`
    Port        int     `yaml:"port"`
    DBName      string  `yaml:"dbname"`
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