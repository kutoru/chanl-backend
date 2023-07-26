package glb

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
)

// Frontend origins
var AllowedOrigins = []string{"http://localhost:5000", "http://192.168.1.12:5000"}
var DB *sql.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
}

// Not sure how this works but in theory it should return current network IP
func GetIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addresses {
			tempIP := address.(*net.IPNet).IP
			if tempIP.DefaultMask().String() == "ffffff00" {
				return tempIP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("could not get the IP")
}
