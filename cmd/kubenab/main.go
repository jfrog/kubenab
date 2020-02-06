package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tlsCertFile string
	tlsKeyFile  string
)

func main() {
	// print Version Informations
	log.Printf("Starting kubenab version %s - %s - %s", version, date, commit)

	// check if all required Flags are set and in a correct Format
	checkArguments()

	flag.StringVar(&tlsCertFile, "tls-cert", "/etc/admission-controller/tls/tls.crt", "TLS certificate file.")
	flag.StringVar(&tlsKeyFile, "tls-key", "/etc/admission-controller/tls/tls.key", "TLS key file.")
	flag.Parse()

	os.LookupEnv("PORT")
	port := "4443"
	if envPort, exists := os.LookupEnv("PORT"); exists {
		port = envPort
	}

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(httpRequestsTotal)
	promRegistry.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/ping", healthCheck)
	http.HandleFunc("/mutate", mutateAdmissionReviewHandler)
	http.HandleFunc("/validate", validateAdmissionReviewHandler)
	s := http.Server{
		Addr: ":" + port,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}

	log.Fatal(s.ListenAndServeTLS(tlsCertFile, tlsKeyFile))
}
