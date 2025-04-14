package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func main(){
	slaveIPs := []string{
		"3.17.110.160:22",
		"3.137.169.157:22",
		"18.221.160.46:22",   
		"3.148.194.157:22",
		"3.142.42.187:22",
		"3.15.211.135:22", 
	}
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

	for _, slave := range slaveIPs{
		client, err := ssh.Dial("tcp", slave, config)
		if err != nil {
			fmt.Printf("Failed to dial %s: %v\n", slave, err)
			continue
		}

		commands := []string{
			"sudo yum update -y",
			"sudo yum install docker -y",
			"sudo systemctl start docker",
			"sudo systemctl enable docker",
			"sudo usermod -aG docker ec2-user",
		}

		for _, command := range commands{
			session, err := client.NewSession()
			if err != nil {
				fmt.Printf("Failed to create session %s: %v", slave, err)
				continue
			}
			output, err := session.CombinedOutput(command)
			if err != nil{
				fmt.Printf("Command %q failed: %v - Output: %s\n", command, err)
			} else {
				fmt.Printf("Output for %q: - %s\n", command, string(output))
			}
			session.Close()
		}
		client.Close()
	}

}
