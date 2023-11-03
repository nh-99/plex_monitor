package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"plex_monitor/internal/database"
	"plex_monitor/internal/discord"
	"plex_monitor/internal/secrets"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Bot parameters
var (
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var secretManager secrets.SecretManager
	secretManager = secrets.NewEnvSecretManager()

	botToken, err := secretManager.GetSecret("DISCORD_BOT_TOKEN")
	if err != nil {
		logrus.Fatalf("Invalid bot parameters: %v", err)
		panic(fmt.Errorf("invalid bot parameters: %v", err))
	}
	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		logrus.Fatalf("Invalid bot parameters: %v", err)
	}

	database.InitDB(secretManager.GetSecretOrDefault("DATABASE_URL", ""), secretManager.GetSecretOrDefault("DATABASE_NAME", ""))
}

var (
	commands        = discord.GetCommands()
	commandHandlers = discord.GetCommandHandlers()
)

func init() {
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logrus.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		logrus.Fatalf("Cannot open the session: %v", err)
		cleanup()
		panic(err)
	}

	logrus.Println("Adding commands...")
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
	if err != nil {
		logrus.Fatalf("Cannot add commands: %v", err)
		cleanup()
		panic(err)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logrus.Println("Press Ctrl+C to exit")
	<-stop

	cleanup()

	logrus.Println("Gracefully shutting down.")
}

func cleanup() {
	if *RemoveCommands {
		logrus.Println("Removing commands...")
		// We need to fetch the commands, since deleting requires the command ID.
		// We are doing this from the returned commands on line 375, because using
		// this will delete all the commands, which might not be desirable, so we
		// are deleting only the commands that we added.
		registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
		if err != nil {
			logrus.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
			if err != nil {
				logrus.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
