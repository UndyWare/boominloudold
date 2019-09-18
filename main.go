package main

import (
	"flag"
	"fmt"

	dt "github.com/undyware/boominloud/discord"
)

const (
	defaultPrefix = ""
	defaultToken  = ""
)

type botOptions struct {
	CommandPrefix string
	Token         string
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

	bot, err := dt.NewBot(opts.Token, opts.CommandPrefix, "with my teetee")
	if err != nil {
		fmt.Printf("error creating bot: %v", err)
	}
	bot.Close()
}

func parseFlags() (*botOptions, error) {
	prefix := flag.String("prefix", defaultToken, "bot command prefix to look for")
	token := flag.String("token", defaultToken, "bot token")
	flag.Parse()

	if *prefix == "" || *token == "" {
		return nil, fmt.Errorf("invalid options. must provide token and prefix")
	}

	opts := &botOptions{
		CommandPrefix: *prefix,
		Token:         *token,
	}
	return opts, nil
}
