package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Build Variables
var (
	BuildVersion = "dev"
	BuildTime    = "unknown"
	GOOS         = runtime.GOOS
	GOARCH       = runtime.GOARCH
)

type Students struct {
	Students [26][26][]Student `json:"students"`
}

type Student struct {
	Initials       [2]string `json:"initials"`
	Grade          int       `json:"grade"`
	TeacherInitial string    `json:"teacher_initial"`
	StudentNumber  int       `json:"student_number"`
}

type Guilds struct {
	Guilds []Guild `json:"guilds"`
}

type Guild struct {
	ID           string     `json:"id"`
	VerifiedRole string     `json:"verified_role"`
	GradeRoles   [12]string `json:"grade_roles"`
	PronounRoles []string   `json:"pronoun_roles"`
}

// Initialize session/bot
var s *discordgo.Session

func init() {
	BotToken := loadEnvVariable("BOT_TOKEN")
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func loadEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func writeToGuilds(guilds *Guilds) error {
	file, err := json.Marshal(guilds)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("guilds.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func contains(s string, a *[]string) bool {
	for _, i := range *a {
		if i == s {
			return true
		}
	}
	return false
}

func getGuildByID(id string) (Guild, int) {
	for gi, g := range guilds.Guilds {
		if g.ID == id {
			return g, gi
		}
	}
	return Guild{}, 0
}

func initEmbed(title string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title: title,
		Color: failureColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made with ❤️ by vidhan#0001",
			IconURL: "https://www.gravatar.com/avatar/30f4d3ebd5b0d0462ea90f7364a6afc4.png",
		},
	}
}

func isAdmin(m *discordgo.Member) bool {
	return m.Permissions&(1<<3) != 0
}

func NewStudent(firstName string, lastName string, grade int, teacherName string, studentNumber int) *Student {
	firstInitial := string(strings.Title(firstName)[0])
	lastInitial := string(strings.Title(lastName)[0])
	teacherFields := strings.Fields(strings.Title(teacherName))
	teacherInitial := string(teacherFields[len(teacherFields)-1][0])

	student := new(Student)

	student.Initials = [2]string{firstInitial, lastInitial}
	student.Grade = grade
	student.TeacherInitial = teacherInitial
	student.StudentNumber = studentNumber

	return student
}

func compareStudents(studentOne *Student, studentTwo *Student) bool {
	if studentOne.Initials[0] == studentTwo.Initials[0] &&
		studentOne.Initials[1] == studentTwo.Initials[1] &&
		studentOne.Grade == studentTwo.Grade &&
		studentOne.TeacherInitial == studentTwo.TeacherInitial {
		return true
	}

	return false
}
func verifyStudent(student *Student, students *Students) bool {
	firstInitialPosition := []rune(student.Initials[0])[0] - 65
	lastInitialPosition := []rune(student.Initials[1])[0] - 65

	if firstInitialPosition < 0 || firstInitialPosition > 26 {
		return false
	}
	if lastInitialPosition < 0 || lastInitialPosition > 26 {
		return false
	}

	initialsArr := students.Students[firstInitialPosition][lastInitialPosition]
	for i := 0; i < len(initialsArr); i++ {
		if compareStudents(student, &initialsArr[i]) {
			return true
		}
	}

	return false
}

// Define constants for the embed colors
const successColor = 0x57F287
const failureColor = 0xED4245

// Define the command formats
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "verify",
			Description: "Verify yourself for access to the server.",
			Options: []*discordgo.ApplicationCommandOption{
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
			Name:        "info",
			Description: "Get information about a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Get info about a user.",
					Required:    true,
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
							Name:        "verified_role",
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
									Description: "Grade. [1-12]",
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
									Description: "Grade. [1-12]",
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
									Description: "Index of the pronoun you want to remove. [1-(# of pronouns)]",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"pronouns_dropdown": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := getGuildByID(i.GuildID)
			selectedPronouns := i.MessageComponentData().Values
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Pronouns Not Set")

			if len(guild.ID) != 0 {
				for pi, pronoun := range guild.PronounRoles {
					if contains(strconv.Itoa(pi), &selectedPronouns) {
						_ = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, pronoun)
						embed.Fields = append(
							embed.Fields, &discordgo.MessageEmbedField{
								Name:  "Pronoun",
								Value: fmt.Sprintf("<@&%s>", pronoun),
							},
						)
					} else {
						_ = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, pronoun)
					}
				}

				embed.Title = "Pronouns Set"
				embed.Color = successColor
			} else {
				embed.Description = "Please ask an administrator to use `/config add pronoun`."
			}

			response.Data.Embeds = []*discordgo.MessageEmbed{&embed}
			err := s.InteractionRespond(i.Interaction, &response)

			if err != nil {
				log.Println(err)
			}
		},
	}

	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"verify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			firstName := strings.Title(i.ApplicationCommandData().Options[0].StringValue())
			lastName := strings.Title(i.ApplicationCommandData().Options[1].StringValue())
			grade := int(i.ApplicationCommandData().Options[2].IntValue())
			teacherName := strings.Title(i.ApplicationCommandData().Options[3].StringValue())
			studentNumber := int(i.ApplicationCommandData().Options[4].IntValue())

			guild, _ := getGuildByID(i.GuildID)
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Not Verified")

			student := NewStudent(firstName, lastName, grade, teacherName, studentNumber)

			studentVerification := verifyStudent(student, students)

			if studentVerification {
				if len(guild.ID) != 0 {
					studentRoles := i.Member.Roles
					for _, r := range guild.GradeRoles {
						if contains(r, &studentRoles) {
							s.GuildMemberRoleRemove(guild.ID, i.Member.User.ID, r)
						}
					}

					s.GuildMemberRoleAdd(guild.ID, i.Member.User.ID, guild.VerifiedRole)
					s.GuildMemberRoleAdd(guild.ID, i.Member.User.ID, guild.GradeRoles[student.Grade-1])
					s.GuildMemberNickname(guild.ID, i.Member.User.ID, firstName+" "+string(lastName[0])+".")

					teacherFields := strings.Fields(teacherName)

					embed.Title = "Verified"
					embed.Color = successColor
					embed.Fields = []*discordgo.MessageEmbedField{
						{
							Name:  "Name",
							Value: strings.Join([]string{firstName, lastName}, " "),
						},
						{
							Name:  "Grade",
							Value: strconv.Itoa(grade),
						},
						{
							Name:  "Teacher Name",
							Value: teacherFields[len(teacherFields)-1],
						},
						{
							Name:  "Student Number",
							Value: strconv.Itoa(studentNumber),
						},
					}
				} else {
					embed.Description = "Please ask an administrator to use `/config set verified_role`."
				}
			}

			response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

			err := s.InteractionRespond(i.Interaction, &response)

			if err != nil {
				log.Println(err)
			}
		},

		"set": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := getGuildByID(i.GuildID)
			memberRoles := i.Member.Roles
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Information Not Set")

			switch i.ApplicationCommandData().Options[0].Name {
			case "pronouns":
				embed.Title = "Pronouns Not Set"

				var dropdown discordgo.MessageComponent

				var pronounOptions []discordgo.SelectMenuOption

				if len(guild.PronounRoles) != 0 {
					for pi, p := range guild.PronounRoles {
						pronounRole, _ := s.State.Role(guild.ID, p)
						pronounOption := discordgo.SelectMenuOption{
							Label:   pronounRole.Name,
							Value:   strconv.Itoa(pi),
							Default: false,
						}

						if contains(p, &memberRoles) {
							pronounOption.Default = true
						}
						pronounOptions = append(pronounOptions, pronounOption)
					}

					dropdown = discordgo.SelectMenu{
						CustomID:    "pronouns_dropdown",
						Placeholder: "Pronouns",
						MinValues:   1,
						MaxValues:   len(pronounOptions),
						Options:     pronounOptions,
					}

					embed.Title = "Set Pronouns"
					embed.Color = successColor

					response.Data.Components = []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{dropdown},
						},
					}
				} else {
					embed.Description = "Please ask an administrator to use `/config add pronoun`."
				}
				response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

				err := s.InteractionRespond(i.Interaction, &response)

				if err != nil {
					log.Println(err)
				}
			}
		},

		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, _ := getGuildByID(i.GuildID)
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Information Not Listed")

			switch i.ApplicationCommandData().Options[0].Name {
			case "pronouns":
				embed.Title = "Pronouns Not Listed"

				if len(guild.PronounRoles) != 0 {
					for _, p := range guild.PronounRoles {
						pronounrole, _ := s.State.Role(guild.ID, p)

						embed.Fields = append(embed.Fields,
							&discordgo.MessageEmbedField{
								Name:  pronounrole.Name,
								Value: fmt.Sprintf("<@&%s>", p),
							},
						)
					}

					embed.Title = "Pronouns List"
					embed.Color = successColor

				} else {
					embed.Description = "Please ask an administrator to use `/config add pronoun`."
				}

				response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

				err := s.InteractionRespond(i.Interaction, &response)

				if err != nil {
					log.Println(err)
				}

			case "grades":
				embed.Title = "Grades Not Listed"

				if len(guild.GradeRoles) != 0 {
					for gi, g := range guild.GradeRoles {

						if len(g) != 0 {
							embed.Fields = append(embed.Fields,
								&discordgo.MessageEmbedField{
									Name:  fmt.Sprintf("Grade %d", gi+1),
									Value: fmt.Sprintf("<@&%s>", g),
								},
							)
						}
					}

					embed.Title = "Grade List"
					embed.Color = successColor

				} else {
					embed.Description = "Please ask an administrator to use `/config add grade`."
				}

				response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

				err := s.InteractionRespond(i.Interaction, &response)

				if err != nil {
					log.Println(err)
				}
			}
		},

		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(s)

			guild, _ := getGuildByID(i.GuildID)
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Information not obtained")

			member, err := s.GuildMember(guild.ID, user.ID)

			if err != nil {
				embed.Description = "Error while getting member."
			} else {
				studentRoles := member.Roles

				for ri, r := range guild.GradeRoles {
					if contains(r, &studentRoles) {
						embed.Fields = append(embed.Fields,
							&discordgo.MessageEmbedField{
								Name:  "Grade",
								Value: strconv.Itoa(ri + 1),
							},
						)
					}
				}

				for _, p := range guild.PronounRoles {
					if contains(p, &studentRoles) {
						embed.Fields = append(embed.Fields,
							&discordgo.MessageEmbedField{
								Name:  "Pronoun",
								Value: fmt.Sprintf("<@&%s>", p),
							},
						)
					}
				}

				embed.Title = fmt.Sprintf("Information for %s", user.Username)
				embed.Color = successColor
			}

			response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

			err = s.InteractionRespond(i.Interaction, &response)

			if err != nil {
				log.Println(err)
			}
		},

		"config": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			guild, guildIndex := getGuildByID(i.GuildID)
			response := discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			}
			embed := initEmbed("Configuration Not Modified")

			if isAdmin(i.Member) {
				switch i.ApplicationCommandData().Options[0].Name {
				case "set":
					embed.Title = "Option Not Set"

					switch i.ApplicationCommandData().Options[0].Options[0].Name {
					case "verified_role":
						roleID := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, "").ID
						guild := Guild{
							ID:           i.GuildID,
							VerifiedRole: roleID,
						}

						var newGuilds Guilds

						for _, g := range guilds.Guilds {
							if g.ID != guild.ID {
								newGuilds.Guilds = append(newGuilds.Guilds, g)
							}
						}

						newGuilds.Guilds = append(newGuilds.Guilds, guild)
						guilds = &newGuilds
						err := writeToGuilds(guilds)

						if err != nil {
							embed.Description = err.Error()
						} else {
							embed.Title = "Verified Role Set"
							embed.Fields = []*discordgo.MessageEmbedField{
								{
									Name:  "Verified Role",
									Value: fmt.Sprintf("<@&%s>", roleID),
								},
							}
							embed.Color = successColor
						}
					}

				case "add":
					embed.Title = "Option Not Added"
					if len(guild.ID) != 0 {
						switch i.ApplicationCommandData().Options[0].Options[0].Name {
						case "grade":
							grade := int(i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue())
							roleID := i.ApplicationCommandData().Options[0].Options[0].Options[1].RoleValue(s, "").ID

							if 1 <= grade && grade <= 12 {
								guilds.Guilds[guildIndex].GradeRoles[grade-1] = roleID
								err := writeToGuilds(guilds)

								if err != nil {
									embed.Description = err.Error()
								} else {
									embed.Title = "Grade Added"
									embed.Fields = []*discordgo.MessageEmbedField{
										{
											Name:  "Grade",
											Value: strconv.Itoa(grade),
										},
										{
											Name:  "Grade Role",
											Value: fmt.Sprintf("<@&%s>", roleID),
										},
									}
									embed.Color = successColor
								}
							} else {
								embed.Description = "Grade must be in range: [1-12]"
							}

						case "pronoun":
							roleID := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, "").ID
							guilds.Guilds[guildIndex].PronounRoles = append(guilds.Guilds[guildIndex].PronounRoles, roleID)
							err := writeToGuilds(guilds)

							if err != nil {
								embed.Description = err.Error()
							} else {
								embed.Title = "Pronoun Added"
								embed.Fields = []*discordgo.MessageEmbedField{
									{
										Name:  "Pronoun Role",
										Value: fmt.Sprintf("<@&%s>", roleID),
									},
								}
								embed.Color = successColor
							}
						}
					} else {
						embed.Description = "Please run `/config set verified_role` to initialize this server first."
					}

				case "remove":
					embed.Title = "Option Not Removed"
					if len(guild.ID) != 0 {
						switch i.ApplicationCommandData().Options[0].Options[0].Name {
						case "grade":
							grade := int(i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue())

							if 1 <= grade && grade <= 12 {
								guilds.Guilds[guildIndex].GradeRoles[grade-1] = ""
								err := writeToGuilds(guilds)

								if err != nil {
									log.Println(err)
								}

								if err != nil {
									embed.Description = err.Error()
								} else {
									embed.Title = "Grade Removed"
									embed.Fields = []*discordgo.MessageEmbedField{
										{
											Name:  "Grade",
											Value: strconv.Itoa(grade + 1),
										},
									}
									embed.Color = successColor
								}
							} else {
								embed.Description = "Grade must be in range: [1-12]"
							}

						case "pronoun":
							pronounIndex := int(i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue())

							if 1 <= pronounIndex && pronounIndex <= len(guild.PronounRoles) {
								pronounRoles := guilds.Guilds[guildIndex].PronounRoles
								guilds.Guilds[guildIndex].PronounRoles = append(pronounRoles[:pronounIndex-1], pronounRoles[pronounIndex:]...)
								err := writeToGuilds(guilds)

								if err != nil {
									log.Println(err)
								}

								if err != nil {
									embed.Description = err.Error()
								} else {
									embed.Title = "Pronoun Removed"
									embed.Color = successColor
								}
							} else {
								embed.Description = fmt.Sprintf("Pronoun index must be in range: [1-%d]", len(guild.PronounRoles))
							}
						}
					} else {
						embed.Description = "Please run `/config set verified_role` to initialize your server first."
					}
				}
			} else {
				embed.Description = "You must have administrator permissions to use `/config`."
			}

			response.Data.Embeds = []*discordgo.MessageEmbed{&embed}

			err := s.InteractionRespond(i.Interaction, &response)

			if err != nil {
				log.Println(err)
			}
		},
	}
)
var students *Students
var guilds *Guilds

func init() {
	s.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
					h(s, i)
				}

			case discordgo.InteractionMessageComponent:
				if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
					h(s, i)
				}
			}
		},
	)
}

func main() {
	s.AddHandler(
		func(s *discordgo.Session, r *discordgo.Ready) {
			log.Println("Bot is up!")
			log.Println(fmt.Sprintf("Built on %s (%s) with %s/%s", BuildVersion, BuildTime, GOOS, GOARCH))
		},
	)

	s.AddHandler(
		func(s *discordgo.Session, c *discordgo.Connect) {
			err := s.UpdateListeningStatus("/verify")
			if err != nil {
				log.Println(err)
			}
		},
	)

	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	file, err := ioutil.ReadFile("students.json")

	if err != nil {
		log.Fatalf("Could not open students.json")
	}

	err = json.Unmarshal([]byte(file), &students)

	if err != nil {
		log.Fatalf("Could not unmarshal students.json")
	}

	file, err = ioutil.ReadFile("guilds.json")

	if err != nil {
		log.Fatalf("Could not open guilds.json")
	}

	err = json.Unmarshal([]byte(file), &guilds)

	if err != nil {
		log.Fatalf("Could not unmarshal guilds.json")
	}

	guildID := loadEnvVariable("GUILD_ID")

	for _, c := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, c)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", c.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}
