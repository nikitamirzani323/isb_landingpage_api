package entities

type Model_domain struct {
	Domain_name string `json:"domain_name"`
}
type Controller_domain struct {
	Client_hostname string `json:"client_hostname" validate:"required"`
}
