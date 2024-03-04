package initializers

import (
	"example.com/ramen/models/Agency"
	"example.com/ramen/models/Company"
	"example.com/ramen/models/channel"
	"example.com/ramen/models/file"
	_map "example.com/ramen/models/map"
	"example.com/ramen/models/reference"
	"example.com/ramen/models/resources"
	"example.com/ramen/models/role"
	"example.com/ramen/models/user"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	//DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	err = DB.AutoMigrate(
		&user.User{},
		&Agency.Agency{},
		&Company.Company{},
		&role.Role{},
		&file.File{},
		&reference.Reference{},
		&_map.Map{},
		&_map.RoleMap{},
		&_map.AgencyMap{},
		&channel.Channel{},
		&resources.Resources{},
	)
	if err != nil {
		log.Fatal("Migration Failed:  \n", err.Error())
		os.Exit(1)
	}

	log.Println("🚀 Connected Successfully to the Database")
}
