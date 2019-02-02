package hugoio

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

// RefreshHugoBlog calls the hugo exectutable rebuild a static site with modified blog posts
func RefreshHugoBlog(hugoDir string) (outStr string, errStr string) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("hugo")
	cmd.Dir = hugoDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		fmt.Println(outStr)
		fmt.Println(errStr)
		log.Fatalln(err)
	}
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	return outStr, errStr
}
