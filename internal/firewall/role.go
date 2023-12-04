package firewall

type Role struct {
	Exist bool   `json:"exist"`
	Name  string `json:"name"`
	Super bool   `json:"super"`
}
