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
			"sudo docker build -t fastapi_app .",
			"sudo docker run -d -p 8000:8000 fastapi_app",
		}

		for _, cmd := range commands{
			session, err := client.NewSession()
			if err != nil{
				fmt.Printf("Failed to create session: %v", err)
				continue
			}

			output, err := session.CombinedOutput(cmd)
			if err != nil {
				fmt.Printf("Command %q failed: %v - Output: %s", cmd, err, string(output))
			} else {
				fmt.Printf("Output for %q: %s\n", cmd, string(output))
			}
			session.Close()
		}
	}
	fmt.Println("Docker service on each slave is succesful!")
}
