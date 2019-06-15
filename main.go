package main

import (
	"flag"
	"fmt"

	"github.com/undyware/boominloud/dt"
)

const (
	defaultPrefix   = ""
	defaultToken    = ""
	defaultUserName = "disc-bot"
)

type botOptions struct {
	CommandPrefix string
	//Name          string
	Token string
}

func main() {
	fmt.Println("bot started")
	opts, err := parseFlags()
	if err != nil {
		fmt.Printf("failed parsing opts: %v\n", err)
		return
	}
	fmt.Printf("prefix: %s\n", opts.CommandPrefix)
	fmt.Printf("token: %s\n", opts.Token)

	bot, err := dt.NewBot(opts.Token, "booted and zooted")
	if err != nil {
		fmt.Printf("error creating bot: %v", err)
	}
	bot.Close()
	//discapi.Say("hello")
}

func parseFlags() (*botOptions, error) {
	prefix := flag.String("prefix", defaultToken, "bot command prefix to look for")
	token := flag.String("token", defaultToken, "bot token")
	flag.Parse()

	opts := &botOptions{
		CommandPrefix: *prefix,
		//Name:          *uName,
		Token: *token,
	}
	return opts, nil
}
