package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func RedisConn() *redis.Client {
	e := godotenv.Load("registerLogin/conf.env")
	if e != nil {
		fmt.Print(e)
	}

	host := os.Getenv("REDIS_ADDR")
	port := os.Getenv("REDIS_PORT")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       db,
	})
	return client
}

// func RsshConn() (net.Conn, error) {
// 	const SSH_ADDRESS = "159.89.206.107:22"
// 	const SSH_USERNAME = "root"
// 	const SSH_PASSWORD = "$A86qiP2$mFd"
// 	const REDIS = "localhost:6379"

// 	sshConfig := &ssh.ClientConfig{
// 		User:            SSH_USERNAME,
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(SSH_PASSWORD),
// 		},
// 	}

// 	netConn, err := net.Dial("tcp", SSH_ADDRESS)
// 	if err != nil {
// 		return nil, err
// 	}
// 	clientConn, chans, reqs, err := ssh.NewClientConn(netConn, SSH_ADDRESS, sshConfig)
// 	if err != nil {
// 		netConn.Close()
// 		return nil, err
// 	}
// 	client := ssh.NewClient(clientConn, chans, reqs)
// 	conn, err := client.Dial("tcp", REDIS)
// 	if err != nil {
// 		client.Close()
// 		return nil, err
// 	}
// 	return conn, err
// }
