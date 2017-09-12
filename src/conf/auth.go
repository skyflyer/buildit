package conf

import (
	"crypto/x509"
	"encoding/pem"
	"log"

	"golang.org/x/crypto/ssh"

	"io/ioutil"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// AuthMethod for repository authentication
var AuthMethod transport.AuthMethod

// ConfigureAuth configures global auth mechanism
func ConfigureAuth(auth *AuthConf) {
	if auth == nil {
		log.Fatalf("Could not configure authentication parameters!\n")
	}

	if auth.Key != "" {
		key, err := ioutil.ReadFile(auth.Key)
		if err != nil {
			log.Fatalf("Unable to read private key: %v", err)
		}

		block, rest := pem.Decode(key)
		if len(rest) > 0 {
			log.Fatalf("Extra data included in key!\n")
		}

		if x509.IsEncryptedPEMBlock(block) {
			der, err := x509.DecryptPEMBlock(block, []byte(auth.Password))
			if err != nil {
				log.Fatalf("Could not decrypt private key! %v\n", err)
			}
			key = pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der})
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("Unable to parse private key: %v", err)
		}
		log.Println("Using SSH public key auth")
		AuthMethod = &gitssh.PublicKeys{User: auth.Username, Signer: signer}
		return
	}

	if auth.UseSSHAgent {
		return
	}

	log.Println("Using username/password authentication")
	// default is username/password
	AuthMethod = &gitssh.Password{User: auth.Username, Pass: auth.Password}
	return
}
