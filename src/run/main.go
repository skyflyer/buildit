package main

import (
	"build"
	"conf"
	"flag"
	"fmt"
	"log"
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
	log.Printf("Running\r\n")
	cfg, err := conf.ReadConfig()
	util.Check(err)

	conf.ConfigureAuth(cfg.Auth)

	r, err := git.NewFilesystemRepository(conf.GitDir)
	util.Check(err)

	remotes, err := r.Remotes()
	util.Check(err)
	if len(remotes) == 0 {
		build.CloneRepository(r, cfg)
	} else {
		build.UpdateRepository(r)
	}

	lastHead := build.GetLastHead()
	branch := "master"

	if cfg.Branch != "" {
		branch = cfg.Branch
	}
	// find the branch
	head, err := r.Reference(
		plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branch)),
		true)
	if err != nil {
		log.Fatalf("Branch %s not found: %s\r\n", cfg.Branch, err)
	}

	if lastHead != head.Hash().String() || runEvenWhenNoChanges {
		build.CleanBuildDir()
		build.CheckoutFiles(r, head)
		build.RunBuildSteps(cfg)
		build.SaveLastHead(head.Hash().String())
	} else {
		log.Println("No changes")
	}
}
