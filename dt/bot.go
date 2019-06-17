package dt

import (
	"errors"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	botID     string
	cmdPrefix string
)

//Bot struct that holds the discord bot session
type Bot struct {
	Prefix string
	dgo.Session
}

//NewBot creates new bot
func NewBot(token string, prefix string, status string) (*Bot, error) {
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
	cmdPrefix = prefix

	<-make(chan struct{})
	return nil, nil
}

func messageHandler(session *dgo.Session, message *dgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Ignore self and other bots
		return
	}

	if message.Content[0:len(cmdPrefix)] == cmdPrefix {
		vs, err := findUserVoiceState(session, user.ID, message.GuildID)
		if err != nil {
			fmt.Printf("Error fetching VoiceState: %v\n", err)
			return
		}

		err = session.ChannelMessageDelete(message.ChannelID, message.ID)
		if err != nil {
			fmt.Printf("Error deleting message: %v\n", err)
		}

		vc, err := session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		defer vc.Disconnect()
		if err != nil {
			fmt.Printf("error joining voice channel: %v", err)
		}
		if err := loadAudio("./airhorn.mp3", vc); err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}
	}
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
