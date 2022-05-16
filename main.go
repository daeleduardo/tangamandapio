package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"time"
    "strings"
	"github.com/go-redis/redis/v8"
)

var (
	hostname string
    err error
)
func main() {
	hostname , err = shell_command("hostname")
	go sub()
	pub()
}

func shell_command(command string) (string, error) {

	cmd := exec.Command("bash", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil

}

func pub() {

	ctx := context.Background()

	rdb := conn()

	defer rdb.Close()

	err := rdb.Publish(ctx, hostname, random()).Err()

	if err != nil {
		panic(err)
	}

	fmt.Println("pub alright")
}

func conn() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func sub() {

	ctx := context.Background()

	rdb := conn()

	defer rdb.Close()

	sub := rdb.Subscribe(ctx, hostname)
	worker(sub)

	fmt.Println("sub alright")
}

func worker(pubsub *redis.PubSub) {
	chn := pubsub.Channel()
	msg := <-chn
	fmt.Println(msg.Payload)

}

func random() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}
