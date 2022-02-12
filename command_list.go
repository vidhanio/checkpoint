package checkpoint

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (c *Checkpoint) listGrades(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}
	if len(g.GradeRoles) == 0 {
		warningRespond(s, i, "Checkpoint is not configured. Please ask someone with `Manage Roles` permissions to use `/config set verified`.")

		return
	}

	var roles []string
	for _, role := range g.GradeRoles {
		if role != "" {
			roles = append(roles, role)
		}
	}

	if len(roles) == 0 {
		warningRespond(s, i, "There are no grade roles set up. Please ask someone with `Manage Roles` permissions to use `/config set grade_roles`.")

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Grade Roles",
		Color: green,
	}

	for i, role := range g.GradeRoles {
		field := &discordgo.MessageEmbedField{
			Name: "Grade " + strconv.Itoa(i+1),
		}

		if role != "" {
			field.Value = "<@&" + role + ">"
		} else {
			field.Value = "None"
		}

		embed.Fields = append(embed.Fields, field)
	}

	respond(s, i, ephemeralify(embedResponse(embed)))
}

func (c *Checkpoint) listPronouns(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}
	if len(g.PronounRoles) == 0 {
		warningRespond(s, i, "Pronoun roles are not set up. Please ask someone with `Manage Roles` permissions to use `/config add pronoun`.")

		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Pronoun Roles",
		Color: green,
	}

	guildRoles, err := c.session.GuildRoles(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	for i, role := range g.PronounRoles {
		for _, guildRole := range guildRoles {
			if guildRole.ID == role {
				field := &discordgo.MessageEmbedField{
					Name: "Pronoun " + strconv.Itoa(i+1),
				}

				field.Value = "<@&" + role + ">"

				embed.Fields = append(embed.Fields, field)
			}
		}
	}

	respond(s, i, ephemeralify(embedResponse(embed)))
}
