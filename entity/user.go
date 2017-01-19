package entity

// User entity
type User struct {
	ID        int    `storm:"ID,increment"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email" storm:"unique"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
	Role      int    `json:"Role"`
}
