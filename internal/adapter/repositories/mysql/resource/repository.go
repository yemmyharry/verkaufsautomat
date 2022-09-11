package resource

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"verkaufsautomat/internal/core/domain/resource"
	"verkaufsautomat/internal/core/logger"
)

func ComparePassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m MachineRepositoryDB) AutoPopulateRoleTable() {
	role := resource.Role{
		RoleName: "buyer",
	}
	m.db.Create(&role)
	role2 := resource.Role{
		RoleName: "seller",
	}
	m.db.Create(&role2)
}

func (m MachineRepositoryDB) AutoPopulatePermissionTable() {
	permission := resource.Permission{
		PermissionName: "create_product",
	}
	m.db.Create(&permission)
	permission2 := resource.Permission{
		PermissionName: "delete_product",
	}
	m.db.Create(&permission2)
	permission3 := resource.Permission{
		PermissionName: "update_product",
	}
	m.db.Create(&permission3)
	permission4 := resource.Permission{
		PermissionName: "buy_product",
	}
	m.db.Create(&permission4)
	permission5 := resource.Permission{
		PermissionName: "deposit_money",
	}
	m.db.Create(&permission5)
}

func (m MachineRepositoryDB) AssignPermissionToRole() {
	var role resource.Role
	m.db.Where("role_name = ?", "buyer").First(&role)
	var permission resource.Permission
	m.db.Where("permission_name = ?", "buy_product").First(&permission)
	var permission2 resource.Permission
	m.db.Where("permission_name = ?", "deposit_money").First(&permission2)
	rolePermission := resource.RolePermission{
		RoleID:       role.RoleId,
		PermissionID: permission.PermissionId,
	}
	m.db.Create(&rolePermission)
	rolePermission2 := resource.RolePermission{
		RoleID:       role.RoleId,
		PermissionID: permission2.PermissionId,
	}
	m.db.Create(&rolePermission2)
}

func (m MachineRepositoryDB) AssignPermissionToRole2() {

	var role resource.Role
	m.db.Where("role_name = ?", "seller").First(&role)
	var permission resource.Permission
	m.db.Where("permission_name = ?", "create_product").First(&permission)
	var permission2 resource.Permission
	m.db.Where("permission_name = ?", "delete_product").First(&permission2)
	var permission3 resource.Permission
	m.db.Where("permission_name = ?", "update_product").First(&permission3)
	rolePermission := resource.RolePermission{
		RoleID:       role.RoleId,
		PermissionID: permission.PermissionId,
	}
	m.db.Create(&rolePermission)
	rolePermission2 := resource.RolePermission{
		RoleID:       role.RoleId,
		PermissionID: permission2.PermissionId,
	}
	m.db.Create(&rolePermission2)
	rolePermission3 := resource.RolePermission{
		RoleID:       role.RoleId,
		PermissionID: permission3.PermissionId,
	}
	m.db.Create(&rolePermission3)
}

func (m MachineRepositoryDB) CreateProduct(product *resource.Product) error {
	m.db.Create(&product)
	return nil
}

func (m MachineRepositoryDB) DeleteProduct(product resource.Product) error {
	m.db.Delete(&product)
	return nil
}

func (m MachineRepositoryDB) UpdateProduct(product resource.Product) error {
	m.db.Save(&product)
	return nil
}

func (m MachineRepositoryDB) GetProductById(id int) (resource.Product, error) {
	var product resource.Product
	m.db.Where("product_id = ?", id).First(&product)
	return product, nil
}

func (m MachineRepositoryDB) GetProducts() ([]resource.Product, error) {
	var products []resource.Product
	m.db.Find(&products)
	return products, nil
}

func (m MachineRepositoryDB) Register(user *resource.User) error {
	var user2 resource.User
	m.db.Where("username = ?", user.Username).First(&user2)
	if user2.Username != "" {
		logger.Error("User already exists")
		return fmt.Errorf("user already exists")
	}
	result := m.db.Create(user)
	return result.Error
}

func (m MachineRepositoryDB) Login(user *resource.User) error {

	InputPassword := user.Password
	fmt.Println("InputPassword: ", InputPassword)
	result := m.db.Where("username = ?", user.Username).First(user)
	if result.Error != nil {
		logger.Error("User does not exist")
		return result.Error
	}

	_, err := ComparePassword(user.Password, InputPassword)
	if err != nil {
		logger.Error("Password is incorrect")
		return err
	}

	return nil
}

func (m MachineRepositoryDB) HealthCheck() error {
	return nil
}

func (m MachineRepositoryDB) GetUserIdAndRoleId(username string) (uint, uint, error) {
	var user resource.User
	m.db.Where("username = ?", username).First(&user)
	return user.UserID, user.RoleID, nil
}
