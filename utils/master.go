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
		fmt.Println("Couldn't read key :/", err)	
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println("Couldn't parse key :/", err)
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
		fmt.Println("Failed to dial: %v", err)
	}
	defer client.Close()

	commands := []string{
		"sudo yum update -y",
		"sudo yum install nginx -y",
		"sudo systemctl start nginx",
		"sudo systemctl enable nginx",
	}

	for _, command := range commands{
		session, err := client.NewSession()
		if err != nil {
			fmt.Println("Failed to create session: %v", err)
			continue
		}
		output, err := session.CombinedOutput(command)
		if err != nil{
			fmt.Println("Command %q failed: %v - Output: %s\n", command, err)
		} else {
			fmt.Println("Output for %q: - %s\n", command, string(output))

		}
		session.Close()
	}


}
