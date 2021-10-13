package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/bwmarrin/discordgo"
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

func includes(s string, a *[]string) bool {
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
	student.StudentNumber = studentNumber
	student.TeacherInitial = teacherInitial

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

func verifyStudent(student *Student, students *Students) (bool, error) {

	firstInitialPosition := []rune(student.Initials[0])[0] - 65
	if firstInitialPosition < 0 || firstInitialPosition > 26 {
		return false, errors.New("first initial not a letter")
	}

	lastInitialPosition := []rune(student.Initials[1])[0] - 65
	if lastInitialPosition < 0 || lastInitialPosition > 26 {
		return false, errors.New("last initial not a letter")
	}

	initialsArr := students.Students[firstInitialPosition][lastInitialPosition]

	for i := 0; i < len(initialsArr); i++ {
		if compareStudents(student, &initialsArr[i]) {
			return true, nil
		}
	}

	return false, nil
}

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
					Description: "Add a configuration option to a group.",
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
									Description: "Index of the pronoun you want to remove. [0-<number of pronouns - 1>]",
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
			var response *discordgo.InteractionResponse
			var messageContent string

			if len(guild.ID) != 0 {
				for p, pronoun := range guild.PronounRoles {
					if includes(strconv.Itoa(p), &selectedPronouns) {
						_ = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, pronoun)
					} else {
						_ = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, pronoun)
					}
				}

				var pronounsStringArr []string

				for _, p := range selectedPronouns {
					pi, err := strconv.Atoi(p)

					if err != nil {
						fmt.Println(err.Error())
					}

					pronounsStringArr = append(pronounsStringArr, fmt.Sprintf("<@&%s>", guild.PronounRoles[pi]))
				}

				messageContent = "Set pronouns: " + strings.Join(pronounsStringArr, ", ")

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: messageContent,
						Flags:   1 << 6,
					},
				}
			} else {
				messageContent = "Please ask an administrator to use `/config add pronoun_role`"
			}

			err := s.InteractionRespond(i.Interaction, response)

			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"verify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			firstName := strings.Title(i.ApplicationCommandData().Options[0].StringValue())
			lastName := strings.Title(i.ApplicationCommandData().Options[1].StringValue())
			grade := i.ApplicationCommandData().Options[2].IntValue()
			teacherName := i.ApplicationCommandData().Options[3].StringValue()
			studentNumber := i.ApplicationCommandData().Options[4].IntValue()

			guild, _ := getGuildByID(i.GuildID)
			var response *discordgo.InteractionResponse
			var messageContent string

			student := NewStudent(firstName, lastName, int(grade), teacherName, int(studentNumber))

			studentVerification, err := verifyStudent(student, students)

			if err != nil {
				messageContent = "Error: " + err.Error()
			} else {
				if studentVerification {

					if len(guild.ID) != 0 {

						for _, gradeRole := range guild.GradeRoles {
							_ = s.GuildMemberRoleRemove(guild.ID, i.Member.User.ID, gradeRole)
						}
						_ = s.GuildMemberRoleAdd(guild.ID, i.Member.User.ID, guild.VerifiedRole)
						_ = s.GuildMemberRoleAdd(guild.ID, i.Member.User.ID, guild.GradeRoles[student.Grade-7])
						_ = s.GuildMemberNickname(guild.ID, i.Member.User.ID, firstName+" "+string(lastName[0])+".")

						messageContent = "You are verified."
					} else {
						messageContent = "Please ask an administrator to use `/config set verified_role`."
					}
				} else {
					messageContent = "Sorry, your information is invalid."
				}
			}

			response = &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: messageContent,
					Flags:   1 << 6,
				},
			}

			err = s.InteractionRespond(i.Interaction, response)

			if err != nil {
				fmt.Println(err.Error())
			}
		},
		"set": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			guild, _ := getGuildByID(i.GuildID)
			memberRoles := i.Member.Roles
			var response *discordgo.InteractionResponse
			var messageContent string
			var components []discordgo.MessageComponent

			var pronounOptions []discordgo.SelectMenuOption

			switch i.ApplicationCommandData().Options[0].Name {
			case "pronouns":
				if len(guild.PronounRoles) != 0 {
					for pi, p := range guild.PronounRoles {
						pronounRole, _ := s.State.Role(guild.ID, p)

						pronounOption := discordgo.SelectMenuOption{
							Label:   pronounRole.Name,
							Value:   strconv.Itoa(pi),
							Default: false,
						}

						if includes(p, &memberRoles) {
							pronounOption.Default = true
						}
						pronounOptions = append(pronounOptions, pronounOption)
					}

					messageContent = "Set your pronouns with the dropdown below."

					components = []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "pronouns_dropdown",
									Placeholder: "Pronouns",
									MinValues:   1,
									MaxValues:   len(pronounOptions),
									Options:     pronounOptions,
								},
							},
						},
					}
				} else {
					messageContent = "Please ask an administrator to use `/config add pronoun_role`"
				}
			}

			response = &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:    messageContent,
					Flags:      1 << 6,
					Components: components,
				},
			}

			err := s.InteractionRespond(i.Interaction, response)

			if err != nil {
				fmt.Println(err.Error())
			}
		},
		"config": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			guild, guildIndex := getGuildByID(i.GuildID)
			var response *discordgo.InteractionResponse
			var messageContent string

			if isAdmin(i.Member) {

				switch i.ApplicationCommandData().Options[0].Name {
				case "set":
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

						messageContent = fmt.Sprintf("Set verified role: <@&%s>", roleID)

						err := writeToGuilds(guilds)

						if err != nil {
							fmt.Println(err.Error())
						}
					}
				case "add":
					if len(guild.ID) != 0 {
						switch i.ApplicationCommandData().Options[0].Options[0].Name {
						case "grade":
							grade := i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue()
							roleID := i.ApplicationCommandData().Options[0].Options[0].Options[1].RoleValue(s, "").ID

							if 1 <= grade && grade <= 12 {
								guilds.Guilds[guildIndex].GradeRoles[grade-1] = roleID
								messageContent = fmt.Sprintf("Added grade %d role: <@&%s>", grade, roleID)

								err := writeToGuilds(guilds)

								if err != nil {
									fmt.Println(err.Error())
								}
							} else {
								messageContent = "Grade must be in range: [1-12]"
							}

						case "pronoun":
							pronounRole := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, "").ID

							guilds.Guilds[guildIndex].PronounRoles = append(guilds.Guilds[guildIndex].PronounRoles, pronounRole)

							err := writeToGuilds(guilds)

							if err != nil {
								fmt.Println(err.Error())
							}

							messageContent = fmt.Sprintf("Added pronouns role: <@&%s>", pronounRole)
						}
					} else {
						messageContent = "Please run `/config set verified_role` to initialize your server first."
					}
				case "remove":
					if len(guild.ID) != 0 {
						switch i.ApplicationCommandData().Options[0].Options[0].Name {
						case "grade":
							grade := i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue()

							if 1 <= grade && grade <= 12 {
								guilds.Guilds[guildIndex].GradeRoles[grade-1] = ""

								err := writeToGuilds(guilds)

								if err != nil {
									fmt.Println(err.Error())
								}

								messageContent = fmt.Sprintf("Remove grade %d role.", grade)
							} else {
								messageContent = "Grade must be in range: [1-12]"
							}

						case "pronoun":
							pronounIndex := int(i.ApplicationCommandData().Options[0].Options[0].Options[0].IntValue())

							if 0 <= pronounIndex && pronounIndex <= len(guild.PronounRoles)-1 {
								pronounRoles := guilds.Guilds[guildIndex].PronounRoles
								guilds.Guilds[guildIndex].PronounRoles = append(pronounRoles[:pronounIndex], pronounRoles[pronounIndex+1:]...)

								err := writeToGuilds(guilds)

								if err != nil {
									fmt.Println(err.Error())
								}

								messageContent = "Removed pronoun role."
							} else {
								messageContent = fmt.Sprintf("Pronoun index must be in range: [0-%d]", len(guild.PronounRoles)-1)
							}
						}
					} else {
						messageContent = "Please run `/config set verified_role` to initialize your server first."
					}
				}

				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: messageContent,
						Flags:   1 << 6,
					},
				}

				err := s.InteractionRespond(i.Interaction, response)

				if err != nil {
					fmt.Println(err.Error())
				}
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
		},
	)
	s.AddHandler(
		func(s *discordgo.Session, c *discordgo.Connect) {
			err := s.UpdateListeningStatus("/verify")

			if err != nil {
				panic(err)
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

	for _, c := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", c)
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
