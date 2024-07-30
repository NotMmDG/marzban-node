package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/NotMmDG/marzban-node/restservice"
	"github.com/NotMmDG/marzban-node/rpycservice"
)

var (
	SERVICE_PORT          = 62050 // Replace with actual config or env var
	SERVICE_PROTOCOL      = "rpyc" // Replace with actual config or env var
	SSL_CERT_FILE         = "/var/lib/marzban-node/ssl_cert.pem"
	SSL_KEY_FILE          = "/var/lib/marzban-node/ssl_key.pem"
	SSL_CLIENT_CERT_FILE  = "" // Replace with actual config or env var
	DEBUG                 = true
)

func generateSSLCertificates() error {
	cert, key, err := generateCertificate()
	if err != nil {
		return err
	}
	if err := os.WriteFile(SSL_KEY_FILE, []byte(key), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(SSL_CERT_FILE, []byte(cert), 0600); err != nil {
		return err
	}
	return nil
}

func main() {
	// Ensure SSL certificates exist or generate them
	if _, err := os.Stat(SSL_CERT_FILE); os.IsNotExist(err) {
		if _, err := os.Stat(SSL_KEY_FILE); os.IsNotExist(err) {
			if err := generateSSLCertificates(); err != nil {
				log.Fatalf("Error generating SSL certificates: %v", err)
			}
		}
	}

	if SSL_CLIENT_CERT_FILE == "" {
		log.Println("Warning: Running without SSL_CLIENT_CERT_FILE, this is not secure!")
	} else if _, err := os.Stat(SSL_CLIENT_CERT_FILE); os.IsNotExist(err) {
		log.Fatalf("Client's certificate file specified in SSL_CLIENT_CERT_FILE is missing")
	}

	switch strings.ToLower(SERVICE_PROTOCOL) {
	case "rpyc":
		startRPYCServer()
	case "rest":
		startRESTServer()
	default:
		log.Fatalf("SERVICE_PROTOCOL is not any of (rpyc, rest).")
	}
}

func startRPYCServer() {
	cert, err := tls.LoadX509KeyPair(SSL_CERT_FILE, SSL_KEY_FILE)
	if err != nil {
		log.Fatalf("Failed to load SSL certificates: %v", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", SERVICE_PORT))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", SERVICE_PORT, err)
	}

	server := grpc.NewServer(grpc.Creds(creds))
	rpycservice.RegisterXrayService(server)
	log.Printf("Node service running on :%d", SERVICE_PORT)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func startRESTServer() {
	if SSL_CLIENT_CERT_FILE == "" {
		log.Fatalf("SSL_CLIENT_CERT_FILE is required for REST service.")
	}
	// Configure REST server settings with SSL/TLS
	certManager := autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.com"), // Replace with your domain
	}

	tlsConfig := &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", SERVICE_PORT),
		TLSConfig: tlsConfig,
		Handler:   restservice.NewRouter(), // Replace with your router setup
	}

	log.Printf("Node service running on :%d", SERVICE_PORT)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

func generateCertificate() (cert string, key string, err error) {
	// Replace this with your certificate generation logic
	return "cert-data", "key-data", nil
}
