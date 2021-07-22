package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

// `tail -F` like reader
// from https://stackoverflow.com/questions/31120987/tail-f-like-generator
type tailReader struct {
	io.ReadCloser
}

func (t tailReader) Read(b []byte) (int, error) {
	for {
		n, err := t.ReadCloser.Read(b)
		if n > 0 {
			return n, nil
		} else if err != io.EOF {
			return n, err
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func newTailReader(fileName string) (tailReader, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return tailReader{}, err
	}

	if _, err := f.Seek(0, 2); err != nil {
		return tailReader{}, err
	}
	return tailReader{f}, nil
}

// remove log prefix
// from `[15:27:33] [Server thread/INFO]: lnazzz lost connection: Disconnected`
// to   `lnazzz lost connection: Disconnected`
func removePrefix(message string) string {
	prefixPattern := regexp.MustCompile(`^\[.*\]:\s`)
	return prefixPattern.ReplaceAllString(message, "")
}

type discordNotifier struct {
	session *discordgo.Session
}

func newDiscordNotifier(token string) (discordNotifier, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return discordNotifier{}, err
	}
	err = session.Open()
	if err != nil {
		return discordNotifier{}, err
	}
	return discordNotifier{session: session}, nil
}

func (dn discordNotifier) notify(message string) {
	dn.session.ChannelMessageSend(os.Getenv("DISCORD_CHANNEL_ID"), message)
}

func (dn *discordNotifier) close() {
	dn.session.Close()
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic("args length must be 1")
	}
	reader, err := newTailReader(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	notifier, err := newDiscordNotifier(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}
	defer notifier.close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		notifier.notify(removePrefix(scanner.Text()))
	}
}
