package checkpoint

import (
	"github.com/bwmarrin/discordgo"
)

func (c *Checkpoint) setPronouns(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}
	if len(g.PronounRoles) == 0 {
		warningRespond(s, i, "Checkpoint is not configured. Please ask someone with `Manage Roles` permissions to use `/config add pronoun`.")

		return
	}

	dropdown := &discordgo.SelectMenu{
		CustomID:    "pronouns_select",
		Placeholder: "Select your pronouns",
		MinValues:   0,
		MaxValues:   len(g.PronounRoles),
	}

	pronounRoles := make([]discordgo.SelectMenuOption, len(g.PronounRoles))

	guildRoles, err := c.session.GuildRoles(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	for ri, roleID := range g.PronounRoles {
		for _, role := range guildRoles {
			if role.ID == roleID {
				pronounRoles[ri] = discordgo.SelectMenuOption{
					Label:   role.Name,
					Value:   role.ID,
					Default: contains(i.Member.Roles, role.ID),
				}
				break
			}
		}
	}

	dropdown.Options = pronounRoles

	resp := embedResponse(successEmbed("Select your pronouns"))

	resp.Data.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				dropdown,
			},
		},
	}

	respond(s, i, ephemeralify(resp))
}

func (c *Checkpoint) pronounsSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	for _, pr := range g.PronounRoles {
		if contains(i.MessageComponentData().Values, pr) && !contains(i.Member.Roles, pr) {
			err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, pr)
			if err != nil {
				errorRespond(s, i, err)

				return
			}
		} else if !contains(i.MessageComponentData().Values, pr) && contains(i.Member.Roles, pr) {
			err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, pr)
			if err != nil {
				errorRespond(s, i, err)

				return
			}
		}
	}

	successRespond(s, i, "Pronouns updated!")
}
