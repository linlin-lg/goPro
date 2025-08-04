package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	pem, err := ioutil.ReadFile("./config/ca.crt")
	if err != nil {
		log.Fatalf("read cert file err: %v", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(pem)
	if !ok {
		panic(fmt.Sprintf("AppendCerts err"))
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:8888", tlsConfig)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	for {
		readBuf := make([]byte, 1024)
		bsCount, err := conn.Read(readBuf)
		if err != nil {
			fmt.Println("conn read err =", err)
			return
		}
		fmt.Println("Read suc =", string(readBuf[:bsCount]))
	}
}