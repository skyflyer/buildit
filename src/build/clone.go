package build

import (
	"conf"
	"log"
	"os"
	"util"

	"io/ioutil"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

// CleanBuildDir cleans build directory
func CleanBuildDir() {
	log.Printf("Cleaning build directory\r\n")
	err := os.RemoveAll(conf.BuildDir)
	util.Check(err)
	err = os.Mkdir(conf.BuildDir, 0777)
	util.Check(err)
}

// CloneRepository clones git repository
func CloneRepository(r *git.Repository, cfg conf.Conf) {
	log.Printf("Cloning %s...\r\n", cfg.Repo)
	err := r.Clone(&git.CloneOptions{
		URL:  cfg.Repo,
		Auth: conf.AuthMethod,
	})
	util.Check(err)
}

// UpdateRepository updates repository
func UpdateRepository(r *git.Repository) {
	remotes, err := r.Remotes()
	util.Check(err)

	log.Printf("Pulling from %s...\r\n", remotes[0].Config().Name)
	err = r.Pull(&git.PullOptions{
		RemoteName: remotes[0].Config().Name,
		Auth:       conf.AuthMethod,
	})
	util.Check(err)
}

// GetLastHead returns hash of last head
func GetLastHead() string {
	contents, err := ioutil.ReadFile(conf.LastHeadMarker)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) == 0 {
		return ""
	}
	return lines[0]
}

// SaveLastHead saves last head hash
func SaveLastHead(hash string) {
	f, err := os.OpenFile(conf.LastHeadMarker, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Could not save last status!")
	}
	f.WriteString(hash)
	f.Close()
}
