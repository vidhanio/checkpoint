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
	schools []string
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

	for k := range c.students {
		c.schools = append(c.schools, k)
	}

	s.AddHandler(c.handler)

	return c
}

func (c *Checkpoint) Start() error {
	err := c.session.Open()
	if err != nil {
		return err
	}

	for _, s := range c.schools {
		commands[0].Options[0].Choices = append(
			commands[0].Options[0].Choices,
			&discordgo.ApplicationCommandOptionChoice{
				Name:  s,
				Value: s,
			},
		)
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
