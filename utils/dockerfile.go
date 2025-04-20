package main

import (
	"fmt"
	"io/ioutil"
	"strings"

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

	dockerfile_content := strings.TrimSpace(`
		FROM pytorch/pytorch
		WORKDIR /app
		COPY FEC.py /app/FEC.py
		COPY best_gpu_tweak.pth /app/best_gpu_tweak.pth

		RUN pip install torch torchvision matplotlib tqdm fastapi uvicorn python-multipart pillow numpy 

		CMD ["uvicorn", "mnist_train:app", "--host", "0.0.0.0", "--port", "8000"]`)

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


		dockerFile, err := sftpClient.Create("/home/ec2-user/Dockerfile")
		if err != nil {
			fmt.Printf("Failed to upload dockerfile: %v", err)
			return
		}
		defer dockerFile.Close()

		_, err = dockerFile.Write([]byte(dockerfile_content))
		if err != nil{
			fmt.Printf("Could't write to remote file on %s: %v\n", slave, err)
		} else {
			fmt.Printf("Successfully uploaded Dockerfile to %s\n", slave)

		}
		dockerFile.Close()
		sftpClient.Close()
		client.Close()
	}

	fmt.Println("Successfully uploaded dockerfile to all slaves!")

}
