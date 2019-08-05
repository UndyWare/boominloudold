package discord

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

const (
	addedMsg   = "%s has been added to the queue"
	invalidMsg = "invalid command given: %s"
	helpmsg    = `BoominLoud Commands:
	play <url> - queue up audio file found at url
	pause - pause player
	resume - resume player
	stop - stop player and clear queue
	skip - skip to next song in queue
	shuffle - shuffle the songs that are currently in the queue
	vol <integer 0-100> - set volume of player`
)

type DiscordBot struct {
	ID               string
	BotSession       *dgo.Session
	StreamingSession *dca.StreamingSession
	VoiceConnection  *dgo.VoiceConnection
	CommandPrefix    string

	queue []string
}

func NewBot(prefix string, queueSize uint8) *DiscordBot {
	bot := &DiscordBot{
		CommandPrefix: prefix,
		queue:         make([]string, queueSize),
	}
	return bot
}
func (bot *DiscordBot) Connect(token, status string) error {
	session, err := dgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("%s: %v", "error creating bot session")
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
		return fmt.Errorf("%s: %v", "failed to open session", err)
	}
	bot.BotSession = session
	return nil
}

func (bot *DiscordBot) messageHandler(session *dgo.Session, message *dgo.MessageCreate) {
	user := message.Author
	if user.ID == bot.ID || user.Bot {
		//Ignore self and other bots
		return
	}

	if message.Content[0:len(bot.CommandPrefix)] == bot.CommandPrefix {
		vs, err := findUserVoiceState(bot.BotSession, user.ID, message.GuildID)
		if err != nil {
			fmt.Printf("Error fetching VoiceState: %v\n", err)
			return
		}

		err = bot.BotSession.ChannelMessageDelete(message.ChannelID, message.ID)
		if err != nil {
			fmt.Printf("Error deleting message: %v\n", err)
		}

		// we just join the channel before any command is processed here
		vc, err := bot.BotSession.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			fmt.Printf("error joining voice channel: %v\n", err)
		}

		//move the split here

		tokens := strings.Split(message.Content, " ")
		command := tokens[0][len(bot.CommandPrefix):]
		options := tokens[1:]
		if returnMsg := bot.commandHandler(command, options); returnMsg != nil {
			_, err := bot.BotSession.ChannelMessageSend(message.ChannelID, *returnMsg)
			if err != nil {
				fmt.Printf("error sending message: %v\n", err)
			}
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

func (bot *DiscordBot) commandHandler(command string, options []string) *string {
	var msg *string
	switch command {
	case "play":
		// This can double as resume if they dont give an audio file to player
		// TODO add error checking to SetPaused
		if 0 == len(options) {
			bot.StreamingSession.SetPaused(false)
			return nil
		}

		// If the queue is empty, then start up the player
		// May want to replace this with a isPlaying bool
		if len(bot.queue) == 0 {
			bot.queue = append(bot.queue, options[0])
			go bot.player(bot.VoiceConnection)
		} else {
			fmt.Printf("Appending %v to queue...\n", options[1])
			bot.queue = append(bot.queue, options[1])
		}

	case "pause":
		bot.StreamingSession.SetPaused(true)

	case "stop":
		bot.queue = bot.queue[len(bot.queue):]

	case "skip":
		bot.queue = bot.queue[1:]
		err := bot.loadAudio(bot.queue[0], bot.VoiceConnection)
		if err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}

	case "shuffle":
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(bot.queue), func(i, j int) {
			bot.queue[i],
				bot.queue[j] = bot.queue[j], bot.queue[i]
		})

	case "queue":
		msg := fmt.Sprintf("URL Queue: %v\n", bot.queue)

	case "resume":
		if bot.StreamingSession.Paused() {
			bot.StreamingSession.SetPaused(false)
		}

	case "help":
		*msg = helpmsg

	default:
		*msg = "invalid command: "
	}
	return nil
}
