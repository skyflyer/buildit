package conf

import (
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
		return
	}

	if auth.Key != "" {
		key, err := ioutil.ReadFile(auth.Key)
		if err != nil {
			log.Fatalf("Unable to read private key: %v", err)
		}

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("Unable to parse private key: %v", err)
		}
		AuthMethod = &gitssh.PublicKeys{User: auth.Username, Signer: signer}
		return
	}

	// default is username/password
	AuthMethod = &gitssh.Password{User: auth.Username, Pass: auth.Password}
	return
}
