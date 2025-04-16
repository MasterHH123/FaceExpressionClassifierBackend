package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	masterIP := "3.17.179.120:22:22"
	keyPath := "id_rsa"

	localFilePath := "server.go"
	remoteFilePath := "/home/ec2-user/server.go"       // Destination on the master node.

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Printf("Failed to read SSH key: %v\n", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Printf("Failed to parse SSH key: %v\n", err)
	}

	config := &ssh.ClientConfig{
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", masterIP, config)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Printf("Failed to create SFTP client: %v\n", err)
		os.Exit(1)
	}
	defer sftpClient.Close()

	serverFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("Failed to open local file: %v\n", err)
	}
	defer serverFile.Close()

	remoteServerFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		fmt.Printf("Failed to create remote file %s: %v\n", err)
	}
	defer remoteServerFile.Close()

	bytesCopied, err := io.Copy(remoteServerFile, serverFile)
	if err != nil {
		fmt.Printf("Failed to copy file to remote server: %v\n", err)
	}

	fmt.Printf("Successfully uploaded file to master node. Bytes copied: %d\n", bytesCopied)
}

