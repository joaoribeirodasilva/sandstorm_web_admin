package server

type Server struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	SslUse    int    `json:"sslUse"`
	SslVerify int    `json:"sslVerify"`
	SslCert   int    `json:"sslCert"`
	SslKey    int    `json:"sslKey"`
}
