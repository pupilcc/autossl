package model

type CertCommand struct {
	Domain    string `json:"domain"`
	Algorithm string `json:"algorithm"`
}

type Cert struct {
	Code     string   `json:"code"`
	CertCode string   `json:"-"`
	KeyCode  string   `json:"-"`
	Domain   string   `json:"domain"`
	DNSNames []string `json:"dnsNames"`
	Cert     string   `json:"cert"`
	Key      string   `json:"key"`
}
