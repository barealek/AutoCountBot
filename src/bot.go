package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	_ "github.com/barealek/autocountbot/devenv"
)

var (
	SleepAmount  int
	ListenID     string
	CooldownTime int = 0
)

func main() {

	// Parse the sleep amount from the environment variable
	SleepAmt, err := strconv.Atoi(os.Getenv("COOLDOWN_TIME"))
	if err == nil {
		SleepAmount = SleepAmt
	} else {
		fmt.Println("couldn't parse COOLDOWN_TIME: ", err)
		fmt.Println("using default value of 1 second")
	}

	var Token string = os.Getenv("TOKEN")
	ListenID = os.Getenv("LISTENER_ID")

	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("couldn't make a discord session: ", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsAll

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("couldn't make a connection: ", err)
		return
	}

	fmt.Printf("Client is running on %s\n", dg.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Logging out...")
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if m.ChannelID != ListenID {
		return
	}

	parsed_num, err := strconv.Atoi(m.Content)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(CooldownTime) * time.Second)
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(parsed_num+1))
}
