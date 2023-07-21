package entities

type Login struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Ipaddress string `json:"ipaddress" validate:"required"`
	Timezone  string `json:"timezone" validate:"required"`
}
type Home struct {
	Page string `json:"page"`
}
