package model

// Stat is the user of this system.
type Stat struct {
	BaseModel
	URLCount uint  `form:"UrlCount"`
	ResCount uint  `form:"ResCount"`
	ErrCount uint  `form:"ErrCount"`
	Status   uint8 `form:"Status"`
}
