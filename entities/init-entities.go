package entities

type Controller_clientinit struct {
	Client_hostname string `json:"client_hostname" validate:"required"`
}
