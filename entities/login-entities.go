package entities

type Loginmobile struct {
	Username string `json:"username" form:"username" validate:"required"`
	Name     string `json:"nama" form:"nama" validate:"required"`
	Device   string `json:"typedevice" form:"typedevice" validate:"required"`
}
