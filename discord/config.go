package discord

const (
	defaultAddress = "127.0.0.1:9000"
	defaultToken   = ";;"
)

type opts struct {
	Addr string
}

/*
func parseFlags() (*server, error) {
	addr := flag.String("addr", defaultAddress, "streamer address")
	//token := flag.String("token", defaultToken, "bot token")
	flag.Parse()

	if *addr == "" {
		return nil, fmt.Errorf("invalid options. must provide token and addr")
	}

	opts := &opts{
		Addr: *addr,
	}
	return opts, nil
}
*/
