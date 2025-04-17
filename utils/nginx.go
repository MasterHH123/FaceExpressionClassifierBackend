package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

func main(){
	masterIP := "3.17.179.120:22"
	slaveIPs := []string{
		"3.17.110.160",
		"3.137.169.157",
		"18.221.160.46",   
		"3.148.194.157",
		"3.142.42.187",
		"3.15.211.135",
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

	client, err := ssh.Dial("tcp", masterIP, config)
	if err != nil{
		fmt.Printf("Failed to dial %s : %v", client, err)
	}
	defer client.Close()

	var serverLines []string
	for _, ip := range slaveIPs{
		serverLines = append(serverLines, fmt.Sprintf("		server %s:8000;", ip))	
	}

	upstreamServers := strings.Join(serverLines, "\n")

	nginx_config := fmt.Sprintf(`user nginx;
		worker_processes auto;
		error_log /var/log/nginx/error.log warn;
		pid /var/run/nginx.pid;

		events {
			worker_connections 1024;
		}

		http {
			upstream gin_app {
		%s
			}

			server {
				listen 80;

				location / {
					proxy_pass http://gin_app;
					proxy_set_header Host $host;
					proxy_set_header X-Real-IP $remote_addr;
					proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
					proxy_set_header X-Forwarded-Proto $scheme;
				}
			}
		}`, upstreamServers)

	fmt.Println(nginx_config)

	sftpClient, err := sftp.NewClient(client)
	if err != nil{
		fmt.Printf("Failed to open SFTP session: %v", err)
		return
	}
	defer sftpClient.Close()

	ec2Path := "/home/ec2-user/nginx.conf"
	remoteFile, err := sftpClient.Create(ec2Path)
	if err != nil {
		fmt.Printf("Failed to create remote file: %v", err)
		return
	}
	defer remoteFile.Close() 

	_, err = remoteFile.Write([]byte(nginx_config))
	if err != nil{
		fmt.Printf("Could't write to remote file: %v", err)
	}
	remoteFile.Close()

	commands := []string{
		"sudo mv /home/ec2-user/nginx.conf /etc/nginx/nginx.conf",
		"sudo systemctl restart nginx",
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
	fmt.Println("Updated nginx config successfully!")
}
