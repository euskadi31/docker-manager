package entity

// Registry entity
type Registry struct {
	ID       int    `json:"ID" storm:"id,increment"`
	Server   string `json:"Server" storm:"unique"`
	Username string `json:"Username"`
	Password string `json:"-"`
	Token    string `json:"Token"`
}
