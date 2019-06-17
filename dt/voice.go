package dt

import (
	"fmt"
	"io"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func loadAudio(path string, vc *dgo.VoiceConnection) error {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	es, err := dca.EncodeFile(path, options)
	if err != nil {
		return err
	}
	defer es.Cleanup()
	time.Sleep(250 * time.Millisecond)
	if err := vc.Speaking(true); err != nil {
		return fmt.Errorf("error enabling speaking: %v", err)
	}
	defer vc.Speaking(false)

	done := make(chan error)
	dca.NewStream(es, vc, done)
	err = <-done
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
