package dt

import (
	"errors"
	"fmt"
	"strings"
	"math/rand"
	"time"

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
	Token string
	*dgo.Session
	StreamingSession *dca.StreamingSession
	urlQueue []string
}

//NewBot creates new bot
func NewBot(token string, prefix string, status string) (*Bot, error) {
	session, err := dgo.New("Bot " + token)
	bot := Bot{prefix, token, session, nil, nil}
	if err != nil {
		return nil, fmt.Errorf("%s: %+v", "error creating session\n", err)
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
		return nil, fmt.Errorf("%s: %+v", "error opening session\n", err)
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
			fmt.Printf("error joining voice channel: %v\n", err)
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

		// If the queue is empty, then start up the player
		// May want to replace this with a isPlaying bool
		if len(bot.urlQueue) == 0 {
			bot.urlQueue = append(bot.urlQueue, tokens[1])
			go bot.player(vc)
		} else {
			fmt.Printf("Appending %v to queue...\n", tokens[1])
			bot.urlQueue = append(bot.urlQueue, tokens[1])
		}

	case "pause":
		bot.StreamingSession.SetPaused(true)

	case "stop":
		bot.urlQueue = bot.urlQueue[len(bot.urlQueue):]
		bot.urlQueue = append(bot.urlQueue, "./nil.mp3")
		vc.Disconnect()
		err := bot.loadAudio(bot.urlQueue[0], vc)
		if err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}
		return

	case "skip":
		return

	case "shuffle":
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(bot.urlQueue), func(i, j int) {bot.urlQueue[i],
			bot.urlQueue[j] = bot.urlQueue[j], bot.urlQueue[i]})
		return

	case "vol":
		return

	case "queue":
		msg := fmt.Sprintf("URL Queue: %v\n", bot.urlQueue)
		_, err := session.ChannelMessageSend(message.ChannelID, msg)
		if err != nil {
			fmt.Printf("error sending message: %v\n", err)
		}

	case "resume":
		if bot.StreamingSession.Paused() {
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
			fmt.Printf("error sending message: %v\n", err)
		}

	default:
		fmt.Println("Invalid command given.")
		_, err := session.ChannelMessageSend(message.ChannelID, "Invalid command given.")
		if err != nil {
			fmt.Printf("error sending message: %v\n", err)
		}
	}
}


func (bot *Bot) player(vc *dgo.VoiceConnection) {
	// This is here in case the session has ended and is trying to start up again
	finished := false
	if bot.StreamingSession != nil {
		finished, _ = bot.StreamingSession.Finished()
	}

	// If the streaming session is not running yet, then start it up
	if bot.StreamingSession == nil || finished {
		fmt.Println("Starting player...")
		err := bot.loadAudio(bot.urlQueue[0], vc)
		if err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}
	}

	// Continuously check if url queue is not empty
	for len(bot.urlQueue) > 0 {
		// If it isnt empty, continuously check if the current stream is finished
		if finished, _ := bot.StreamingSession.Finished(); finished {
			// If it is finished, then we need to clear that finished url from Queue
			bot.urlQueue = bot.urlQueue[1:]
			// If its empty, then player is done
			if len(bot.urlQueue) == 0 {
				fmt.Println("Player session ending...")
				return
			} else {
				// Now we need to load up the next url to play
				err := bot.loadAudio(bot.urlQueue[0], vc)
				if err != nil {
					fmt.Printf("error loading audio: %v\n", err)
				}
				continue
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
