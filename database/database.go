package database

import (
	"aino-spring.com/aino_site/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connection struct {
  Database *gorm.DB
  Config *config.Config
}

func NewConnetion(conf *config.Config) (*Connection, error) {
  connection := new(Connection)
  connection.Config = conf
  var err error
  connection.Database, err = gorm.Open(mysql.Open(conf.MysqlDsn))
  if err != nil {
    return nil, err
  }
  return connection, nil
}

func (connection *Connection) Migrate() {
  connection.Database.AutoMigrate(&Setting{}, &Page{}, &Post{}, &User{})
}

