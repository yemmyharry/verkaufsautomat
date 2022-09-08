package resource

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type MachineRepositoryDB struct {
	db *gorm.DB
}

func (m MachineRepositoryDB) HealthCheck() error {
	return nil
}

var err = godotenv.Load("verkaufsautomat.env")

var DbUsername = os.Getenv("DB_USER")
var DbPassword = os.Getenv("DB_PASS")
var DbName = os.Getenv("DB_NAME")
var DbHost = os.Getenv("DB_HOST")
var DbPort = os.Getenv("DB_PORT")

func NewMachineRepositoryDB() *MachineRepositoryDB {
	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "charset=utf8mb4&parseTime=True&loc=Local"
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//client.AutoMigrate(&domain.Machine{})
	return &MachineRepositoryDB{client}
}
