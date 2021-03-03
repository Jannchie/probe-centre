package model

// User is the user of this system.
type User struct {
	BaseModel
	Name string
	Mail string
}
