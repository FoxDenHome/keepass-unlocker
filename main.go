package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// systemd-ask-password -n | systemd-creds --tpm2-pcrs=7+8 --user --with-key=host+tpm2 encrypt - ~/.tpm/keepassxc
// systemd-creds --user decrypt ~/.tpm/keepassxc
func getPassword() ([]byte, error) {
	proc := exec.Command("/usr/bin/systemd-creds", "--user", "decrypt", "/home/doridian/.tpm/keepassxc")
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
