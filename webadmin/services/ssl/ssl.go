package ssl

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/config"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/utils"
)

type Ssl struct {
	SslUse    bool   `json:"sslUse"`
	SslVerify bool   `json:"sslVerify"`
	SslCert   string `json:"sslCertPath"`
	SslKey    string `json:"sslKeyPath"`
	ConfigDir string `json:"configDir"`
	Host      string `json:"host"`
	log       *admin_log.Log
}

const (
	MODULE    = "ssl"
	RSA_BITS  = 2048
	CERT_NAME = "web_admin_cert.pem"
	KEY_NAME  = "web_admin_key.pem"
)

var (
	ecdsa_curve string = "P521"
)

func New(conf *config.Configuration, log *admin_log.Log) *Ssl {

	s := new(Ssl)
	s.SslUse = conf.WebAdmin.SslUse
	s.SslVerify = conf.WebAdmin.SslVerify
	s.SslCert = conf.WebAdmin.SslCert
	s.SslKey = conf.WebAdmin.SslKey
	s.ConfigDir = conf.WebAdmin.ConfigDir
	s.log = log

	return s
}

func (s *Ssl) Load() bool {

	if !s.SslUse {
		return true
	}

	var noCert bool = false
	var noKey bool = false

	d, err := filepath.Abs(s.ConfigDir)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to calculate absolute path from '%s' relative path ", s.ConfigDir), MODULE, admin_log.LOG_CRITICAL)
	}

	if s.SslCert == "" {

		p := path.Join(d, CERT_NAME)
		if !utils.FileExists(p) {
			noCert = true
		} else {
			s.SslCert = p
			s.log.Write(fmt.Sprintf("private certificate found at '%s'", s.ConfigDir), MODULE, admin_log.LOG_INFO)
		}
	}

	if s.SslKey == "" {

		p := path.Join(d, KEY_NAME)
		if !utils.FileExists(p) {
			noKey = true
		} else {
			s.SslKey = p
			s.log.Write(fmt.Sprintf("public certificate key found at '%s'", s.ConfigDir), MODULE, admin_log.LOG_INFO)
		}
	}

	if noCert || noKey {
		if s.SslVerify {
			s.log.Write(fmt.Sprintf("this server can not generate valid SSL certificates and no certificates where found at '%s' ", s.ConfigDir), MODULE, admin_log.LOG_CRITICAL)
		}

		s.create()
	} else {
		if err := s.validCertificates(); err != nil {
			s.log.Write(fmt.Sprintf("the certificates found at '%s' are invalid. ERR: %s", s.ConfigDir, err.Error()), MODULE, admin_log.LOG_CRITICAL)
			return false
		} else {
			s.log.Write(fmt.Sprintf("the certificates found at '%s' are ok", s.ConfigDir), MODULE, admin_log.LOG_INFO)
		}
	}

	return true
}

func (s *Ssl) create() bool {

	var priv any
	var err error

	s.log.Write("generating self signed certificates", MODULE, admin_log.LOG_INFO)

	switch ecdsa_curve {
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		s.log.Write(fmt.Sprintf("unknown ecdsa curve '%s'", ecdsa_curve), MODULE, admin_log.LOG_CRITICAL)
	}

	if err != nil {
		s.log.Write(fmt.Sprintf("failed to generate private key. ERR: %s", err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}

	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment

	validFrom := time.Now()
	validUntil := validFrom.Add(365 * 24 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to generate certificate serial number. ERR: %s", err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Sandstorm Web Admin"},
		},
		NotBefore: validFrom,
		NotAfter:  validUntil,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	ips, err := s.ips()
	if ips == nil || err != nil || len(*ips) == 0 {
		if err == nil {
			err = fmt.Errorf("no ip addresses found")
		}
		s.log.Write(fmt.Sprintf("failed to get host ip addresses from network interfaces. ERR: %s", err), MODULE, admin_log.LOG_CRITICAL)
	}

	template.IPAddresses = append(template.IPAddresses, *ips...)

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, s.publicKey(priv), priv)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to generate SSL certificates. ERR: %s", err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}

	// Write cert.pem
	s.SslCert = path.Join(s.ConfigDir, CERT_NAME)
	certOut, err := os.Create(s.SslCert)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to generate SSL certificates. ERR: %s", err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		s.log.Write(fmt.Sprintf("failed to write data to '%s'. ERR: %s", s.SslCert, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}
	if err := certOut.Close(); err != nil {
		s.log.Write(fmt.Sprintf("error closing '%s'. ERR: %s", s.SslCert, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}

	// Write key.pem
	s.SslKey = path.Join(s.ConfigDir, KEY_NAME)
	keyOut, err := os.OpenFile(s.SslKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to write data to '%s'. ERR: %s", s.SslKey, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to marshal private key '%s'. ERR: %s", s.SslKey, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		s.log.Write(fmt.Sprintf("failed to write data to '%s'. ERR: %s", s.SslKey, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}
	if err := keyOut.Close(); err != nil {
		s.log.Write(fmt.Sprintf("error closing '%s'. ERR: %s", s.SslKey, err.Error()), MODULE, admin_log.LOG_CRITICAL)
	}

	s.log.Write(fmt.Sprintf("certificates '%s' and '%s' generated successfully", s.SslCert, s.SslKey), MODULE, admin_log.LOG_INFO)

	return false
}

func (s *Ssl) publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func (s *Ssl) ips() (*[]net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, 0)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip)
		}
	}

	return &ips, nil
}

func (s *Ssl) validCertificates() error {

	_, err := tls.LoadX509KeyPair(s.SslCert, s.SslKey)
	return err
}

// func (s *Ssl) hostname() error {

// 	var err error
// 	s.Host, err = os.Hostname()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
