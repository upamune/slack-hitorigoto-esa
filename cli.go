package main

import (
	"flag"
	"fmt"
	"io"
	"github.com/BurntSushi/toml"
	"github.com/nlopes/slack"
	"github.com/upamune/go-esa/esa"
	"log"
	"errors"
	"time"
	"strconv"
	"strings"
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

type Config struct{
	Slack string `toml:"SLACK_TOKEN"`
	Esa string `toml:"ESA_TOKEN"`
	Room string `toml:SLACK_ROOM"`
}

func (c *Config) loadConfig(filename string) (error){
	config := Config{}

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return err
	}

	if c.Slack == "" {
		c.Slack = config.Slack
	}
	if c.Esa == "" {
		c.Esa = config.Esa
	}

	if c.Room == "" {
		c.Room = config.Room
	}

	return nil
}

func (c *Config) getSlackHitorigoto() ([]slack.SearchMessage, error){
	var roomMessages []slack.SearchMessage
	api := slack.New(c.Slack)
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	query := "on:" + todayStr + " "

	searchParam := slack.NewSearchParameters()

	messages, err := api.SearchMessages(query, searchParam)
	if err != nil {
		return nil, err
	}
	for _, message := range messages.Matches {
		if message.Channel.Name == c.Room {
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

func (c *Config) postEsaNippo(messages []slack.SearchMessage) (error){
	client := esa.NewClient(c.Esa)
	post := esa.Post{}

	today := time.Now()
	todayStr := today.Format("2006/01/02")
	post.Name = "日報/" + todayStr

	for _, message := range messages {
		hourMinute, err := unixTStoHourMinute(message.Timestamp)
		if err != nil {
			return err
		}

		header := "[" + hourMinute + "]" + "(" + message.Permalink +") "
		body := message.Text

		post.BodyMd += header + body + "\n"
	}
	_, err := client.Post.Create("", post)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validate() (error){
	if c.Esa == "" {
		return errors.New("EsaのTokenを指定してください")
	}
	if c.Slack == "" {
		return errors.New("SlackのTokenを指定してください")
	}
	if c.Room == "" {
		return errors.New("Roomを指定してください")
	}

	return nil
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		s string
		e string
		c string
		r string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&s, "slack", "", "slack access config")
	flags.StringVar(&s, "s", "", "(Short)")

	flags.StringVar(&e, "esa", "", "esa access config")
	flags.StringVar(&e, "e", "", "(Short)")

	flags.StringVar(&r, "room", "", "slack room name")
	flags.StringVar(&r, "r", "", "(Short)")

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

	_ = s

	_ = e

	_ = c

	_ = r


	config := Config{Slack: s, Esa: e, Room: r}
	if c != "" {
		err := config.loadConfig(c)
		if err != nil {
			log.Fatal(err)
			return ExitCodeError
		}
	}
	err := config.validate()
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
