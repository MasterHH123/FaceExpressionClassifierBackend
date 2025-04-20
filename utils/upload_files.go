package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
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

	pythonFile, err := os.Open("../FEC.py")
	if err != nil {
		fmt.Printf("Failed to open file %v", err)
		return
	}
	defer pythonFile.Close()

	modelFile, err := os.Open("../best_gpu_tweak.pth")
	if err != nil{
		fmt.Printf("Failed to open file %v", err)
		return
	}
	defer modelFile.Close()

	for _, slave := range slaveIPs{
		client, err := ssh.Dial("tcp", slave, config)
		if err != nil {
			fmt.Printf("Failed to dial %s: %v\n", slave, err)
			continue
		}
			sftpClient, err := sftp.NewClient(client)
			if err != nil{
				fmt.Printf("Failed to open SFTP session: %v", err)
				return
			}
			defer sftpClient.Close()


			remotePythonFile, err := sftpClient.Create("/home/ec2-user/FEC.py")
			if err != nil{
				fmt.Printf("Failed to create remote file %v", err)
				return
			}
			defer remotePythonFile.Close()

			remoteModelFile, err := sftpClient.Create("/home/ec2-user/best_gpu_tweak.pth")
			if err != nil {
				fmt.Printf("Failed to create remote file %v", err)
				return
			}
			defer remoteModelFile.Close()

			//resets file pointer before copying to avoid copying EOF content
			_, err = pythonFile.Seek(0, 0)
			if err != nil {
				fmt.Printf("Failed to reset pythonFile pointer: %v\n", err)
				return
			}

			_, err = io.Copy(remotePythonFile, pythonFile)
			if err != nil {
				fmt.Printf("Failed to upload python file: %v", err)
				return
			}

			//resets file pointer before copying to avoid copying EOF content
			_, err = modelFile.Seek(0, 0)
			if err != nil {
				fmt.Printf("Failed to reset pythonFile pointer: %v\n", err)
				return
			}

			_, err = io.Copy(remoteModelFile, modelFile)
			if err != nil {
				fmt.Printf("Failed to upload python file: %v", err)
				return
			}

		client.Close()
		fmt.Printf("Successfully uploaded files to %s\n", slave)
	}

	fmt.Println("Successfully uploaded files to all slaves!")

}
