package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

var (
	tlsCertFile string
	tlsKeyFile  string
)

func main() {
	// check if all required Flags are set and in a correct Format
	checkArguments()

	flag.StringVar(&tlsCertFile, "tls-cert", "/etc/admission-controller/tls/tls.crt", "TLS certificate file.")
	flag.StringVar(&tlsKeyFile, "tls-key", "/etc/admission-controller/tls/tls.key", "TLS key file.")
	flag.Parse()

	http.HandleFunc("/ping", healthCheck)
	http.HandleFunc("/mutate", mutateAdmissionReviewHandler)
	http.HandleFunc("/validate", validateAdmissionReviewHandler)
	s := http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}
	log.Fatal(s.ListenAndServeTLS(tlsCertFile, tlsKeyFile))
}
