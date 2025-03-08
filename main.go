package main

import (
	"log"
	"os"
	"os/exec"

	"r00t2.io/gokwallet"
)

func main() {
	r := gokwallet.DefaultRecurseOpts
	r.AllWalletItems = true
	wm, err := gokwallet.NewWalletManager(r, "KeePassUnlocker")
	if err != nil {
		log.Panicln(err)
	}

	wallet := wm.Wallets["kdewallet"]
	folder := wallet.Folders["KeePassXC"]
	pass := []byte(folder.Passwords["Passwords.kdbx"].Value)
	wm.Close()

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
