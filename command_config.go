package checkpoint

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Checkpoint) configSetVerifiedRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	g.VerifiedRole = i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID

	err = c.guilds.WriteOne(g)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Verified role set!")
}

func (c *Checkpoint) configAddGradeRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	grade := i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue()

	if grade < 1 || grade > 12 {
		warningRespond(s, i, "Grade must be between 1 and 12.")

		return
	}

	g.GradeRoles[grade-1] = i.ApplicationCommandData().Options[0].Options[0].Options[1].RoleValue(s, i.GuildID).ID

	err = c.guilds.WriteOne(g)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Grade role added!")
}

func (c *Checkpoint) configAddPronounRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	roleID := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID

	if contains(g.PronounRoles, roleID) {
		warningRespond(s, i, "Pronoun role already exists.")

		return
	}

	g.PronounRoles = append(g.PronounRoles, roleID)

	err = c.guilds.WriteOne(g)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Pronoun role added!")
}

func (c *Checkpoint) configRemoveGradeRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	grade := i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue()

	if grade < 1 || grade > 12 {
		warningRespond(s, i, "Grade must be between 1 and 12.")

		return
	}

	g.GradeRoles[grade-1] = ""

	err = c.guilds.WriteOne(g)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Grade role removed!")
}

func (c *Checkpoint) configRemovePronounRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	roleIndex := int(i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue()) - 1

	if roleIndex < 0 || roleIndex >= len(g.PronounRoles) {
		warningRespond(s, i, "Role index out of bounds.")

		return
	}

	g.PronounRoles = append(g.PronounRoles[:roleIndex], g.PronounRoles[roleIndex+1:]...)

	err = c.guilds.WriteOne(g)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Pronoun role removed!")
}
