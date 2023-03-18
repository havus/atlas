package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/havus/atlas/handler"
	"github.com/joho/godotenv"
)

var discordSession *discordgo.Session

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatalf("Invalid load token: %v", err)
	}

	if discordSession, err = discordgo.New("Bot " + os.Getenv("TOKEN")); err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	commandHandlers := handler.Commands

	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	fmt.Println("Welcome")

	discordSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

		if _, err := discordSession.ChannelMessageSend("1085627994258821221", "Welcome aboard"); err != nil {
			log.Printf("error when send message %v", err)
		}
	})

	if err := discordSession.Open(); err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}
	defer discordSession.Close()

	log.Println("Adding commands...")

	commands := GetCommands()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := discordSession.ApplicationCommandCreate(discordSession.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := discordSession.ApplicationCommandDelete(discordSession.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}
