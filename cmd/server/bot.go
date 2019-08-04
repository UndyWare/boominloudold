package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	ytdl "github.com/rylio/ytdl"
)

var (
	botID     string
	cmdPrefix string
)

//Bot struct that holds the discord bot session
type Bot struct {
	Prefix string
	Token  string
	*dgo.Session
	StreamingSession *dca.StreamingSession
	urlQueue         []string
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

func (bot *Bot) messageHandler(session *dgo.Session, message *dgo.MessageCreate) {
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
		vc.Disconnect()
		return

	case "skip":
		bot.urlQueue = bot.urlQueue[1:]
		err := bot.loadAudio(bot.urlQueue[0], vc)
		if err != nil {
			fmt.Printf("error loading audio: %v\n", err)
		}
		return

	case "shuffle":
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(bot.urlQueue), func(i, j int) {
			bot.urlQueue[i],
				bot.urlQueue[j] = bot.urlQueue[j], bot.urlQueue[i]
		})
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

func fetchVideo(urlstring string) (string, error) {

	fmt.Println("Fetching " + urlstring)
	URL, err1 := url.Parse(urlstring)
	if err1 != nil {
		return "nil", err1
	}
	vid, err2 := ytdl.GetVideoInfoFromURL(URL)
	if err2 != nil {
		return "nil", err2
	}

	title := strings.Replace(vid.Title, " ", "", -1)

	file, err3 := os.Create(title + ".mp4")
	if err3 != nil {
		return "nil", err3
	}
	defer file.Close()
	err := vid.Download(vid.Formats[0], file)
	if err != nil {
		return "nil", err
	}
	return title, nil
}

func convertVideo(title string) string {
	defer os.Remove(title + ".mp4")
	fmt.Println("converting " + title)
	cmd := exec.Command("ffmpeg", "-i", title+".mp4", title+".mp3")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error converting video...%v\n", err)
		return "nil"
	}
	return title + ".mp3"
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

		title, err1 := fetchVideo(bot.urlQueue[0])
		if err1 != nil {
			fmt.Printf("Error Fetching Audio...%v\n", err1)
			return
		}

		mp3 := convertVideo(title)

		err2 := bot.loadAudio(mp3, vc)
		if err2 != nil {
			fmt.Printf("error loading audio: %v\n", err2)
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
				vid, err3 := fetchVideo(bot.urlQueue[0])
				if err3 != nil {
					fmt.Println("Error Fetching Audio...")
					return
				}

				mp3 := convertVideo(vid)

				err4 := bot.loadAudio(mp3, vc)
				if err4 != nil {
					fmt.Printf("error loading audio: %v\n", err4)
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
