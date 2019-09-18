package discord

import (
	"fmt"
	"io"
	"os"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func (bot *Bot) loadAudio(path string, vc *dgo.VoiceConnection) error {
	options := dca.StdEncodeOptions
	es, err := dca.EncodeFile(path, options)
	if err != nil {
		return err
	}
	defer es.Cleanup()
	defer os.Remove(path)
	if err := vc.Speaking(true); err != nil {
		return fmt.Errorf("error enabling speaking: %v", err)
	}
	defer vc.Speaking(false)

	done := make(chan error)
	bot.StreamingSession = dca.NewStream(es, vc, done)
	err = <-done
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
