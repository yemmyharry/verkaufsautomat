package resource

import uuid "github.com/satori/go.uuid"

type Product struct {
	ProductID       uuid.UUID `gorm:"primary_key" json:"product_id"`
	AmountAvailable int       `json:"amount_available"`
	Cost            int       `json:"cost"`
	ProductName     string    `json:"product_name"`
	SellerID        uuid.UUID `json:"seller_id" gorm:"foreign_key:UserID"`
}

type User struct {
	UserID   uuid.UUID `gorm:"primary_key" json:"user_id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Deposit  int       `json:"deposit"`
	RoleID   uuid.UUID `json:"role_id" gorm:"foreign_key:RoleId"`
}

type Role struct {
	RoleId   uuid.UUID `gorm:"primary_key" json:"role_id"`
	RoleName string    `json:"role_name"`
}

type Permission struct {
	PermissionID   uuid.UUID `gorm:"primary_key" json:"permission_id"`
	PermissionName string    `json:"permission_name"`
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" gorm:"foreign_key:RoleId"`
	PermissionID uuid.UUID `json:"permission_id" gorm:"foreign_key:PermissionID"`
}
