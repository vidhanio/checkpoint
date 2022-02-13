package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/vidhanio/checkpoint"
)

func main() {
	studentsFilename := flag.String("students", "students.json", "Path to students.json")
	guildsPathFilename := flag.String("guilds", "guilds.json", "Path to guilds.json")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file.")
	}

	c := checkpoint.New(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_GUILD_ID"), *studentsFilename, *guildsPathFilename)

	err = c.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting bot.")
	}

	log.Info().Msg("Bot started. Press Ctrl+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc

	err = c.Stop()
	if err != nil {
		log.Fatal().Err(err).Msg("Error stopping bot.")
	}

	log.Info().Msg("Bot stopped.")
}
