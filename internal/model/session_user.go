package model

type SessionRole struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type SessionUser struct {
	ID              int64         `json:"id"`
	FirstName       string        `json:"firstName"`
	LastName        string        `json:"lastName"`
	PasswordDefault bool          `json:"passwordDefault"`
	Roles           []SessionRole `json:"roles"`
	Menus           SessionMenu   `json:"menus"`
}
