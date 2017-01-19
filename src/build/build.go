package build

import (
	"conf"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"util"
)

func streamReader(stream io.ReadCloser, name string) {
	buf := make([]byte, 1024)
	var err error
	var cnt int
	for err == nil {
		cnt, err = stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("[%s] Error: %s\r\n", name, err)
			}
		} else {
			log.Printf("[%s] %s", name, strings.Trim(string(buf[:cnt]), "\r\n\t"))
		}
	}
}

// RunBuildSteps runs all of the configured build steps
func RunBuildSteps(cfg conf.Conf) {
	previousCWD, err := os.Getwd()
	util.Check(err)
	os.Chdir(conf.BuildDir)

	for _, step := range cfg.Steps {
		log.Println("About to run command: ", step)
		args := strings.Split(step, " ")
		if len(args) == 0 {
			continue
		}
		cmd := exec.Command(args[0], args[1:]...)
		outPipe, _ := cmd.StdoutPipe()
		errPipe, _ := cmd.StderrPipe()
		go streamReader(outPipe, "OUT")
		go streamReader(errPipe, "ERR")
		err := cmd.Start()
		if err != nil {
			log.Printf("Error executing command: %v\r\n", err)
			continue
		}
		err = cmd.Wait()
		if err != nil {
			log.Printf("Command exited with error: %v\r\n", err)
			break
		}
	}

	os.Chdir(previousCWD)
}
