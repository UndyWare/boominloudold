package dt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jonas747/dca"
	dgo "github.com/bwmarrin/discordgo"
)

var (
	botID     string
	cmdPrefix string
)

//Bot struct that holds the discord bot session
type Bot struct {
	Prefix string
	*dgo.Session
	StreamingSession *dca.StreamingSession
	isPaused bool
}

//NewBot creates new bot
func NewBot(token string, prefix string, status string) (*Bot, error) {
	session, err := dgo.New("Bot " + token)
	bot := Bot{prefix, session, nil, true}
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
	session.AddHandler(bot.messageHandler)

	err = session.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", "error opening session", err)
	}

	defer session.Close()
	cmdPrefix = prefix

	<-make(chan struct{})
	return &bot, nil
}

func (bot *Bot) messageHandler (session *dgo.Session, message *dgo.MessageCreate) {
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

		// we just join the channel before any command is processed here
		vc, err := session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			fmt.Printf("error joining voice channel: %v", err)
		}

		bot.commandHandler(message, session, vc)
	}
}

func (bot *Bot) commandHandler(message *dgo.MessageCreate, session *dgo.Session, vc *dgo.VoiceConnection) {
	switch tokens := strings.Split(message.Content, " "); tokens[0][len(cmdPrefix):] {
	case "play":
		// This can double as resume if they dont give an audio file to player
		// TODO add error checking to SetPaused
		if len(tokens) < 2 {
			bot.StreamingSession.SetPaused(false)
			return
		}
		err := bot.loadAudio(tokens[1], vc)
		if err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}

	case "pause":
		bot.StreamingSession.SetPaused(true)
		bot.isPaused = true


	case "stop":
		return

	case "skip":
		return

	case "shuffle":
		return

	case "vol":
		return

	case "queue":
		return

	case "resume":
		if bot.isPaused {
			bot.StreamingSession.SetPaused(false)
		}

	case "help":
		helpmsg := `BoominLoud Commands:
	bl.play <url> - queue up audio file found at url
	bl.pause - pause player
	bl.resume - resume player
	bl.stop - stop player and clear queue
	bl.skip - skip to next song in queue
	bl.shuffle - shuffle the songs that are currently in the queue
	bl.vol <integer 0-100> - set volume of player`

		_, err := session.ChannelMessageSend(message.ChannelID, helpmsg)
		if err != nil {
			fmt.Printf("error sending message: %v", err)
		}

	default:
		fmt.Println("Invalid command given.")
		_, err := session.ChannelMessageSend(message.ChannelID, "Invalid command given.")
		if err != nil {
			fmt.Printf("error sending message: %v", err)
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
