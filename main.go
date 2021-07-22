package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
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

type notifier interface {
	notify(message string)
}

type stdNotifier struct {
}

func (sn stdNotifier) notify(message string) {
	fmt.Println(message)
}

func main() {
	fp, err := newTailReader("/tmp/sample.log")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	mainNotifier := stdNotifier{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		mainNotifier.notify(removePrefix(scanner.Text()))
	}
}
