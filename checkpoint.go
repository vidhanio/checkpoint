package checkpoint

import (
	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/checkpoint/guild"
	"github.com/vidhanio/checkpoint/student"
)

type Checkpoint struct {
	session  *discordgo.Session
	students student.Students
	guilds   *guild.GuildStore

	guildID string
}

func New(botToken string, guildID string, studentsFilename string, guildsPathFilename string) *Checkpoint {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		panic(err)
	}

	c := &Checkpoint{
		session:  s,
		students: student.MakeStudents(studentsFilename),
		guilds:   guild.NewGuilds(guildsPathFilename),
		guildID:  guildID,
	}

	s.AddHandler(c.handler)

	return c
}

func (c *Checkpoint) Start() error {
	err := c.session.Open()
	if err != nil {
		return err
	}

	_, err = c.session.ApplicationCommandBulkOverwrite(c.session.State.User.ID, c.guildID, commands)
	if err != nil {
		return err
	}

	return c.guilds.Open()
}

func (c *Checkpoint) Stop() error {
	return c.session.Close()
}
