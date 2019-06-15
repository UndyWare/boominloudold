package dt

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	botID string
)

//Bot struct that holds the discord bot session
type Bot struct {
	Prefix string
	dgo.Session
}

//NewBot creates new bot
func NewBot(token string, status string) (*Bot, error) {
	session, err := dgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", "error creating session", err)
	}

	session.AddHandler(func(discord *dgo.Session, ready *dgo.Ready) {
		err = session.UpdateStatus(0, status)
		if err != nil {
			//log error fmt.Errorf("%s: %+v", "error creating session", err)
			return
		}
		servers := session.State.Guilds
		fmt.Printf("bot started on %d servers\n", len(servers))
	})
	session.AddHandler(textHandler)

	err = session.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", "error opening session", err)
	}

	defer session.Close()

	<-make(chan struct{})
	return nil, nil
}

func textHandler(session *dgo.Session, message *dgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Ignore self and other bots
		return
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
	msg := fmt.Sprintf("%s said %s in channel %s", message.Author, message.Content, message.ChannelID)
	fmt.Printf("sending msg: %s\n", msg)
	_, err := session.ChannelMessageSend(message.ChannelID, msg)
	if err != nil {
		fmt.Printf("error sending message: %v", err)
	}
}

func textSend(session *dgo.Session, msg string) {
	fmt.Printf("sending msg: %s\n", msg)
	//discord.Channel
}
