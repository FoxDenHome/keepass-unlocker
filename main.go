package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// tpm2_createprimary -c ~/.tpm/primary.ctx -Q
// tpm2_createpolicy -Q --policy-pcr -l sha256:7 -L ~/.tpm/pcr7.policy
// echo 'PASSWORD' | tpm2_create -C ~/.tpm/primary.ctx -L ~/.tpm/pcr7.policy -i- -c ~/.tpm/keepass.ctx -Q
func getPassword() ([]byte, error) {
	proc := exec.Command("/usr/bin/tpm2_unseal", "-c", "/home/doridian/.tpm/keepass.ctx", "-p", "pcr:sha256:7")
	proc.Stderr = os.Stderr
	stdout, err := proc.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = proc.Start()
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	return data, proc.Wait()
}

func main() {
	pass, err := getPassword()
	if err != nil {
		panic(err)
	}

	proc := exec.Command("/usr/bin/keepassxc", "--pw-stdin", "/home/doridian/Sync/KeePass/Passwords.kdbx")
	stdin, err := proc.StdinPipe()
	if err != nil {
		panic(err)
	}
	err = proc.Start()
	if err != nil {
		panic(err)
	}
	_, err = stdin.Write(pass)
	if err != nil {
		panic(err)
	}
	stdin.Close()
	log.Printf("Completed startup and unlock!")
	os.Exit(0)
}
