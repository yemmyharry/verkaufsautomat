package resource

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"verkaufsautomat/internal/core/domain/resource"
)

type MachineRepositoryDB struct {
	db *gorm.DB
}

var err = godotenv.Load("verkaufsautomat.env")

var DbUsername = os.Getenv("MYSQL_USER")
var DbPassword = os.Getenv("MYSQL_PASSWORD")
var DbName = os.Getenv("MYSQL_DATABASE")
var DbHost = os.Getenv("MYSQL_DB_HOST")
var DbPort = os.Getenv("MYSQL_DB_PORT")

func NewMachineRepositoryDB() *MachineRepositoryDB {
	dsn := DbUsername + ":" + DbPassword + "@tcp" + "(" + DbHost + ":" + DbPort + ")/" + DbName + "?" + "charset=utf8mb4&parseTime=True&loc=Local"
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	client.AutoMigrate(&resource.Product{}, &resource.User{}, &resource.Role{}, &resource.Permission{}, &resource.RolePermission{})

	autoPopulateRoleTable := MachineRepositoryDB{db: client}
	autoPopulateRoleTable.AutoPopulateRoleTable()
	autoPopulatePermissionTable := MachineRepositoryDB{db: client}
	autoPopulatePermissionTable.AutoPopulatePermissionTable()
	assignPermissionToRole := MachineRepositoryDB{db: client}
	assignPermissionToRole.AssignPermissionToRole()
	assignPermissionToRole2 := MachineRepositoryDB{db: client}
	assignPermissionToRole2.AssignPermissionToRole2()

	return &MachineRepositoryDB{client}
}
