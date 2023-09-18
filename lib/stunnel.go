package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func serverAliveCheck(client *ssh.Client) (err error) {
	// This is ported version of Open SSH client server_alive_check function
	// see: https://github.com/openssh/openssh-portable/blob/b5e412a8993ad17b9e1141c78408df15d3d987e1/clientloop.c#L482
	_, _, err = client.SendRequest("keepalive@openssh.com", true, nil)
	return
}

// StartKeepalive starts sending server keepalive messages until done channel
// is closed.
func StartKeepalive(client *ssh.Client, interval time.Duration, countMax int, done <-chan struct{}) {
	t := time.NewTicker(interval)
	defer t.Stop()

	n := 0
	for {
		select {
		case <-t.C:
			if err := serverAliveCheck(client); err != nil {
				n++
				if n >= countMax {
					client.Close()
					return
				}
			} else {
				n = 0
			}
		case <-done:
			return
		}
	}
}

// 转发
func sForward(serverAddr string, remoteAddr string, localConn net.Conn, config *ssh.ClientConfig, closeTunnel func()) {
	// 设置sshClientConn
	sshClientConn, err := ssh.Dial("tcp", serverAddr, config)
	if err != nil {
		fmt.Printf("ssh.Dial failed1: %s", err)
		return
	}
	// 设置Connection
	sshConn, err := sshClientConn.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("ssh.Dial failed2: %s", err)
		defer sshClientConn.Close()
		return
	}

	go StartKeepalive(sshClientConn, 30*time.Second, 2, make(<-chan struct{}))

	// 将localConn.Reader复制到sshConn.Writer
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			fmt.Printf("io.Copy failed: %v", err)
			defer sshClientConn.Close()
			defer sshConn.Close()
			return
		}
	}()

	// 将sshConn.Reader复制到localConn.Writer
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			fmt.Printf("io.Copy failed: %v", err)
			defer sshClientConn.Close()
			defer sshConn.Close()
			return
		}
	}()
}

func publicKeyAuthFunc(keyPath string) ssh.AuthMethod {
	// keyPath, err := homedir.Expand(kPath)
	// if err != nil {
	// 	log.Fatal("find key's home dir failed", err)
	// }

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func Tunnel(username string, privateKeyPath string, serverAddr string, remoteAddr string, localAddr string, start func(), closeTunnel func()) {
	// 设置SSH配置
	// fmt.Printf("%s，服务器：%s；远程：%s；本地：%s\n", "设置SSH配置", serverAddr, remoteAddr, localAddr)
	sshConfig := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	config := &ssh.ClientConfig{
		Config: sshConfig,
		User:   username,
		Auth: []ssh.AuthMethod{
			// ssh.Password(password),
			ssh.AuthMethod(publicKeyAuthFunc(privateKeyPath)),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// 设置本地监听器
	localListener, err := net.Listen("tcp", localAddr)
	if err != nil {
		closeTunnel()
		log.Fatalf("net.Listen failed: %v\n", err)
		return
	}
	start()
	go func() {
		for {
			// 设置本地
			var localConn net.Conn
			localConn, err = localListener.Accept()
			if err != nil {
				closeTunnel()
				log.Fatalf("localListener.Accept failed: %v\n", err)
				return
			}
			go sForward(serverAddr, remoteAddr, localConn, config, closeTunnel)
		}
	}()

}
