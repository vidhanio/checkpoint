package checkpoint

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "verify",
		Description: "Verify yourself for access to the server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "school",
				Description: "Your school.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "first_name",
				Description: "Your first name.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "last_name",
				Description: "Your last name.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "grade",
				Description: "Your grade.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "teacher_name",
				Description: "The last name of your homeroom teacher (Week 1, Period 1)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "student_number",
				Description: "Your student number (6 digits).",
				Required:    true,
			},
		},
	},
	{
		Name:        "set",
		Description: "Set information about yourself.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "pronouns",
				Description: "Set your pronouns",
			},
		},
	},
	{
		Name:        "list",
		Description: "List information for options.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "pronouns",
				Description: "List all pronouns.",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "grades",
				Description: "List all grades.",
			},
		},
	},
	{
		Name:        "config",
		Description: "Configure Checkpoint.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "set",
				Description: "Set a configuration option.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "verified",
						Description: "Set the verified role.",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionRole,
								Name:        "role",
								Description: "Role to assign to verified users.",
								Required:    true,
							},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "add",
				Description: "Add a configuration option.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "grade",
						Description: "Add a grade role.",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "grade",
								Description: "Grade. [1 - 12]",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionRole,
								Name:        "role",
								Description: "Role to assign to users in this grade.",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "pronoun",
						Description: "Add a pronoun role.",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionRole,
								Name:        "role",
								Description: "Role to assign to users with these pronouns.",
								Required:    true,
							},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "remove",
				Description: "Remove a configuration option from a group.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "grade",
						Description: "Remove a grade role.",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "grade",
								Description: "Grade. [1 - 12]",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "pronoun",
						Description: "Remove a pronoun role.",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "pronoun_index",
								Description: "Index of the pronoun you want to remove. [1 - <# of pronouns>]",
								Required:    true,
							},
						},
					},
				},
			},
		},
	},
}

func (c *Checkpoint) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "verify":
			c.verify(s, i)
		case "set":
			switch i.ApplicationCommandData().Options[0].Name {
			case "pronouns":
				c.setPronouns(s, i)
			}
		case "list":
			switch i.ApplicationCommandData().Options[0].Name {
			case "grades":
				c.listGrades(s, i)
			case "pronouns":
				c.listPronouns(s, i)
			}
		case "config":
			if !isManageRoles(s, i) {
				warningRespond(s, i, "You do not have permission to configure Checkpoint.")

				return
			}
			switch i.ApplicationCommandData().Options[0].Name {
			case "set":
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "verified":
					c.configSetVerifiedRole(s, i)
				}
			case "add":
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "grade":
					c.configAddGradeRole(s, i)
				case "pronoun":
					c.configAddPronounRole(s, i)
				}
			case "remove":
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "grade":
					c.configRemoveGradeRole(s, i)
				case "pronoun":
					c.configRemovePronounRole(s, i)
				}
			}
		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "pronouns_select":
			c.pronounsSelect(s, i)
		}
	}
}
