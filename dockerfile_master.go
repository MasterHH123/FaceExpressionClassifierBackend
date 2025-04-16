package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/sftp"
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

	dockerfile_content := strings.TrimSpace(`
		FROM golang:1.23.0
		WORKDIR /app
		COPY go.mod go.sum server.go ./
		RUN go mod tidy

		RUN CGO_ENABLED=0 GOOS=linux go build -o /server
		EXPOSE 8080
		CMD ["/server"]
	`)

	client, err := ssh.Dial("tcp", masterIP, config)
	if err != nil {
		fmt.Printf("Failed to dial %s: %v\n", masterIP, err)
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
			fmt.Printf("Failed to create session %v", err)
		}
		output, err := session.CombinedOutput(command)
		if err != nil {
			fmt.Printf("Command %q failed %v - Output: %s\n", command, err)
		} else {
			fmt.Printf("Output for %q: - %s\n", command, string(output))
		}
		session.Close()
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil{
		fmt.Printf("Failed to open SFTP session: %v", err)
		return
	}
	defer sftpClient.Close()


	dockerFile, err := sftpClient.Create("/home/ec2-user/Dockerfile")
	if err != nil {
		fmt.Printf("Failed to upload dockerfile: %v", err)
		return
	}
	defer dockerFile.Close()

	_, err = dockerFile.Write([]byte(dockerfile_content))
	if err != nil{
		fmt.Printf("Could't write to remote file on %s: %v\n", masterIP, err)
	} else {
		fmt.Printf("Successfully uploaded Dockerfile to %s\n", masterIP)

	}
	dockerFile.Close()
	sftpClient.Close()
	client.Close()

	fmt.Println("Successfully uploaded dockerfile to master node!")

}
