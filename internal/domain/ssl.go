package domain

type Cert struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type CertCommand struct {
	Domain    string `json:"domain"`
	Algorithm string `json:"algorithm"`
}

type CertDTO struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
