package model

type CertCommand struct {
	Domain    string `json:"domain"`
	Algorithm string `json:"algorithm"`
}

type Cert struct {
	Code   string `json:"code"`
	Domain string `json:"domain"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}
