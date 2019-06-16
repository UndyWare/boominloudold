package dt

import (
	"fmt"
	"errors"
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
	session, err := dgo.New("Bot " + token)
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
	session.AddHandler(messageHandler)

	err = session.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", "error opening session", err)
	}

	defer session.Close()

	<-make(chan struct{})
	return nil, nil
}


func findUserVoiceState(session *dgo.Session, userid string, guildID string) (*dgo.VoiceState, error) {
	guild, err := session.Guild(guildID)
	if err != nil {
		fmt.Printf("Error fetching guild: %v\n", err)
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userid {
			return vs, nil
		}
	}
	return nil, errors.New("Could not find user's voice state")
}


func messageHandler(session *dgo.Session, message *dgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Ignore self and other bots
		return
	}

	// Check if message starts with bl.
	// If it doesnt, we just ignore the message
	if message.Content[0:3] == "bl." {
		// Get voice state of user's channel, or an error if the user is not in one
		vs, err := findUserVoiceState(session, user.ID, message.GuildID)
		if err != nil {
			fmt.Printf("Error fetching VoiceState: %v\n", err)
			return
		}

		fmt.Println(vs)
		fmt.Println("Sender ID: " + user.ID)
		err1 := session.ChannelMessageDelete(message.ChannelID, message.ID)
		if err != nil {
			fmt.Printf("Error deleting message: %v\n", err1)
		}

	} else {
		fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
		msg := fmt.Sprintf("%s said %s in channel %s", message.Author, message.Content, message.ChannelID)
		fmt.Printf("Sending: %s\n", msg)
		_, err := session.ChannelMessageSend(message.ChannelID, msg)
		if err != nil {
			fmt.Printf("error sending message: %v", err)
		}
	}
}

func textSend(session *dgo.Session, msg string) {
	fmt.Printf("sending msg: %s\n", msg)
	//discord.Channel
}
