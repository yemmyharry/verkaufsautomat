package resource

type Product struct {
	ProductID       uint   `json:"product_id" gorm:"primary_key;autoIncrement"`
	AmountAvailable int    `json:"amount_available"`
	Cost            int    `json:"cost"`
	ProductName     string `json:"product_name"`
	SellerID        uint   `json:"seller_id" gorm:"foreignKey:UserID"`
}

type User struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	RoleID   uint   `json:"role_id" gorm:"foreignKey:RoleId"`
}

type Role struct {
	RoleId   uint   `json:"roleId" gorm:"primaryKey;autoIncrement"`
	RoleName string `json:"role_name"`
}

type Permission struct {
	PermissionId   uint   `json:"permissionId" gorm:"primaryKey;autoIncrement"`
	PermissionName string `json:"permission_name"`
}

type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"foreignKey:RoleId"`
	PermissionID uint `json:"permission_id" gorm:"foreignKey:PermissionId"`
}
