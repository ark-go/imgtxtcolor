package internal

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"os"
)

func tlsConfig() *tls.Config {
	homedir, _ := os.UserHomeDir()
	crt, err := ioutil.ReadFile(homedir + "/localcert/localcert.pem")
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(homedir + "/localcert/localkey.pem")
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "172.16.172.10",
	}
}
