package domain

type Cert struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type CertDTO struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
