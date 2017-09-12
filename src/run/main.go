package main

import (
	"build"
	"conf"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"util"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	forceRun := flag.Bool("force", false, "force runs the build steps even if there are no changes in the repository")
	watchRepo := flag.Bool("watch", false, "watch starts the agent in a periodic watch mode. It checks every minute for changes. Period can be adjusted by using the -period flag.")
	period := flag.Int("period", 1, "When run in watch mode, defines how often to check the repository. Defaults to 1 minute")
	flag.Parse()

	initial := make(chan bool, 1)
	initial <- true
	quit := make(chan bool, 1)

out:
	for {
		select {
		case <-initial:
			if !(*watchRepo) {
				quit <- true
			}
			runProgram(*forceRun)
		case <-time.After(time.Duration(*period) * time.Minute):
			runProgram(*forceRun)
		case <-quit:
			break out
		}
	}
}

func runProgram(runEvenWhenNoChanges bool) {
	cfg, err := conf.ReadConfig()
	util.Check(err)
	log.Printf("Running. Working directory: %s\n", cfg.WorkingDirectory)

	conf.ConfigureAuth(cfg.Auth)

	branch := "master"
	if cfg.Branch != "" {
		branch = cfg.Branch
	}
	branchReference := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

	var r *git.Repository
	_, err = os.Stat(cfg.WorkingDirectory)
	if err == nil {
		log.Println("Repository directory already exists. Opening repo...")
		r, err = git.PlainOpen(cfg.WorkingDirectory)
		util.Check(err)
		log.Println("Fetching...")
		err = r.Fetch(&git.FetchOptions{
			Auth:       conf.AuthMethod,
			RemoteName: "origin",
			Progress:   os.Stdout,
		})
		util.Check(err)
	} else {
		log.Println("Cloning repository...")
		r, err = git.PlainClone(cfg.WorkingDirectory, false, &git.CloneOptions{
			Auth:          conf.AuthMethod,
			RemoteName:    "origin",
			URL:           cfg.Repo,
			Progress:      os.Stdout,
			ReferenceName: branchReference,
		})

		util.Check(err)
	}

	worktree, err := r.Worktree()
	util.Check(err)

	remoteBranchReference := plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch))

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: remoteBranchReference,
		Force:  true,
	})
	util.Check(err)

	lastHead := build.GetLastHead()

	// find the head
	head, err := r.Head()
	if err != nil {
		log.Fatalf("Head of branch %s not found: %s\n", cfg.Branch, err)
	}
	log.Printf("Current repository head: %s\n", head.Hash().String())

	if lastHead != head.Hash().String() || runEvenWhenNoChanges {
		build.SaveLastHead(head.Hash().String())
		build.RunBuildSteps(cfg)
	} else {
		log.Println("No changes")
	}
}
