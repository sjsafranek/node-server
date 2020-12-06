package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

const (
	DEFAULT_HTTP_PORT         int64  = 8080
	DEFAULT_DATABASE_ENGINE   string = "bolt"
	DEFAULT_DATABASE_DATABASE string = ""
	DEFAULT_DATABASE_PASSWORD string = ""
	DEFAULT_DATABASE_USERNAME string = ""
	DEFAULT_DATABASE_HOST     string = ""
	DEFAULT_DATABASE_PORT     int64  = 0
	DEFAULT_DATABASE_FILE     string = "db.bolt"
)

var (
	HTTP_PORT         = DEFAULT_HTTP_PORT
	DATABASE_ENGINE   = DEFAULT_DATABASE_ENGINE
	DATABASE_DATABASE = DEFAULT_DATABASE_DATABASE
	DATABASE_PASSWORD = DEFAULT_DATABASE_PASSWORD
	DATABASE_USERNAME = DEFAULT_DATABASE_USERNAME
	DATABASE_HOST     = DEFAULT_DATABASE_HOST
	DATABASE_PORT     = DEFAULT_DATABASE_PORT
	DATABASE_FILE     = DEFAULT_DATABASE_FILE
)

func Open(file string) (Configuration, error) {
	var conf Configuration
	return conf, conf.Fetch(file)
}

func New() Configuration {
	var conf Configuration
	conf.Server.Port = HTTP_PORT
	conf.Database.Type = DATABASE_ENGINE
	conf.Database.Database = DATABASE_DATABASE
	conf.Database.Password = DATABASE_PASSWORD
	conf.Database.Username = DATABASE_USERNAME
	conf.Database.Host = DATABASE_HOST
	conf.Database.Port = DATABASE_PORT
	conf.Database.File = DATABASE_FILE
	return conf
}

type Configuration struct {
	Server   ServerConfiguration   `toml:"server"`
	Database DatabaseConfiguration `toml:"database"`
}

type ServerConfiguration struct {
	// Host
	Port int64 `toml:"port"`
}

type DatabaseConfiguration struct {
	Type     string `toml:"type"`
	Database string `toml:"database"`
	Password string `toml:"password"`
	Username string `toml:"username"`
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	File     string `toml:"file"`
}

func (self *Configuration) Fetch(file string) error {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return err
	}
	return self.Unmarshal(string(b))
}

func (self *Configuration) Save(file string) error {
	contents, err := self.Marshal()
	if nil != err {
		return err
	}
	return ioutil.WriteFile(file, []byte(contents), 0644)
}

func (self *Configuration) Unmarshal(data string) error {
	return toml.Unmarshal([]byte(data), self)
}

func (self Configuration) Marshal() (string, error) {
	b, err := toml.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

func (self *DatabaseConfiguration) GetDatabaseConnectionString() (string, error) {
	switch self.Type {
	case "postgres":
		return self.getPostgresDatabaseConnectionString(), nil
	case "bolt":
		return self.getBoltDatabaseConnectionString(), nil
	default:
		return "", errors.New("Unsupported database type")
	}
	return "", nil
}

func (self *DatabaseConfiguration) getBoltDatabaseConnectionString() string {
	return self.File
}

func (self *DatabaseConfiguration) getPostgresDatabaseConnectionString() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable", self.Type, self.Username, self.Password, self.Host, self.Port, self.Database)
}
