package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func main(){
	masterIP := "3.17.179.120:22"
	keyPath := "id_rsa"

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Printf("Couldn't read key %v\n", err)	
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Printf("Couldn't parse key %v\n", err)
	}

	config := &ssh.ClientConfig {
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", masterIP, config)
	if err != nil {
		fmt.Printf("Failed to dial %v\n", err)
	}

	commands := []string{
		"sudo docker build -t gin_app .",
		"sudo docker run -d -p 8080:8080 gin_app",
	}

	for _, cmd := range commands{
		session, err := client.NewSession()
		if err != nil{
			fmt.Printf("Failed to create session: %v", err)
		}

		output, err := session.CombinedOutput(cmd)
		if err != nil {
			fmt.Printf("Command %q failed: %v - Output: %s", cmd, err, string(output))
		} else {
			fmt.Printf("Output for %q: %s\n", cmd, string(output))
			fmt.Println("Docker service on master is successful!")
		}
		session.Close()
	}
}
