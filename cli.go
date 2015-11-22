package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nlopes/slack"
	"github.com/upamune/go-esa/esa"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

type Esa struct {
	Token    string `toml:"token"`
	Team     string `toml:"team"`
	Category string `toml:"category"`
}

type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

type Config struct {
	Slack Slack `toml:"slack"`
	Esa   Esa   `toml:"esa"`
}

func (c *Config) loadConfig(filename string) error {
	_, err := toml.DecodeFile(filename, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getSlackHitorigoto() ([]slack.SearchMessage, error) {
	var roomMessages []slack.SearchMessage
	api := slack.New(c.Slack.Token)
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	todayStr = "2015-11-21"
	query := "on:" + todayStr + " "

	searchParam := slack.NewSearchParameters()

	messages, err := api.SearchMessages(query, searchParam)
	if err != nil {
		return nil, err
	}
	for _, message := range messages.Matches {
		if message.Channel.Name == c.Slack.Channel {
			roomMessages = append(roomMessages, message)
		}
	}

	return roomMessages, nil
}

func unixTStoHourMinute(unixStr string) (string, error) {
	unixStr = unixStr[:strings.Index(unixStr, ".")]
	unixTimeStamp, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return "", err
	}
	t := time.Unix(unixTimeStamp, 0)
	hourStr := strconv.Itoa(t.Hour())
	minuteStr := strconv.Itoa(t.Minute())
	hourMinute := hourStr + ":" + minuteStr

	return hourMinute, nil
}

func (c *Config) postEsaNippo(messages []slack.SearchMessage) error {
	client := esa.NewClient(c.Esa.Token)
	post := esa.Post{}

	today := time.Now()
	todayStr := today.Format("2006/01/02")
	post.Name = c.Esa.Category + "/" + todayStr

	for _, message := range messages {
		hourMinute, err := unixTStoHourMinute(message.Timestamp)
		if err != nil {
			return err
		}

		header := "[" + hourMinute + "]" + "(" + message.Permalink + ") "
		body := message.Text

		post.BodyMd += header + body + "\n"
	}
	_, err := client.Post.Create(c.Esa.Team, post)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validate() error {
	var errorMessages []string
	if c.Esa.Token == "" {
		errorMessages = append(errorMessages, "EsaのTokenを指定してください")
	}
	if c.Esa.Team == "" {
		errorMessages = append(errorMessages, "EsaのTeamを指定してください")
	}
	if c.Esa.Category == "" {
		errorMessages = append(errorMessages, "EsaのCategoryを指定してください")
	}
	if c.Slack.Token == "" {
		errorMessages = append(errorMessages, "SlackのTokenを指定してください")
	}
	if c.Slack.Channel == "" {
		errorMessages = append(errorMessages, "SlackのChannelを指定してください")
	}

	if len(errorMessages) != 0 {
		errorMessage := ""

		for _, message := range errorMessages {
			errorMessage += message + "\n"
		}

		return errors.New(errorMessage)
	}

	return nil
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		c string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&c, "config", "", "config file")
	flags.StringVar(&c, "c", "", "(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	_ = c

	config := Config{}
	err := config.loadConfig(c)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}
	err = config.validate()
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}
	messages, err := config.getSlackHitorigoto()
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}
	err = config.postEsaNippo(messages)
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	return ExitCodeOK
}
