package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

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
	ID           string `json:"id"`
	VerifiedRole string `json:"verified_role"`
}

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "verify",
			Description: "Verify yourself for access to the server.",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "first-name",
					Description: "Your first name.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "last-name",
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
					Name:        "teacher-name",
					Description: "The last name of your homeroom teacher (Week 1, Period 1)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "student-number",
					Description: "Your student number (6 digits).",
					Required:    true,
				},
			},
		},
		{
			Name:        "initialize",
			Description: "Initialize the server with Woodlands Checkpoint.",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "verified-role",
					Description: "The role to give to verified users.",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"verify": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			firstName := strings.Title(i.ApplicationCommandData().Options[0].StringValue())
			lastName := strings.Title(i.ApplicationCommandData().Options[1].StringValue())
			grade := i.ApplicationCommandData().Options[2].IntValue()
			teacherName := i.ApplicationCommandData().Options[3].StringValue()
			studentNumber := i.ApplicationCommandData().Options[4].IntValue()

			student := NewStudent(firstName, lastName, int(grade), teacherName, int(studentNumber))

			studentVerification, err := verifyStudent(student, students)

			var msg string

			if err != nil {
				msg = "Error: " + err.Error()
			} else {
				if studentVerification {
					var currentGuild Guild
					for _, guild := range guilds.Guilds {
						if guild.ID == i.GuildID {
							currentGuild = guild
						}
					}
					msg = "You are verified! Welcome!"
					if len(currentGuild.ID) > 0 {
						_ = s.GuildMemberRoleAdd(currentGuild.ID, i.Member.User.ID, currentGuild.VerifiedRole)
					}

				} else {
					msg = "Sorry, your information is invalid."
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
					Flags:   1 << 6,
				},
			})
		},
		"initialize": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			roleID := i.ApplicationCommandData().Options[0].RoleValue(s, "").ID
			guildID := i.GuildID

			var msg string

			if i.Member.Permissions&(1<<3) != 0 {
				guild := Guild{
					ID:           guildID,
					VerifiedRole: roleID,
				}
				guilds.Guilds = append(guilds.Guilds, guild)
				msg = fmt.Sprintf("Set role to <@&%s>", roleID)

				file, _ := json.Marshal(guilds)

				_ = ioutil.WriteFile("guilds.json", file, 0644)
			} else {
				msg = "You do not have sufficient permissions. You must be an administrator."
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
					Flags:   1 << 6,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
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
		return false, errors.New("first initial not an uppercase character")
	}

	lastInitialPosition := []rune(student.Initials[1])[0] - 65
	if lastInitialPosition < 0 || lastInitialPosition > 26 {
		return false, errors.New("last initial not an uppercase character")
	}

	initialsArr := students.Students[firstInitialPosition][lastInitialPosition]

	for i := 0; i < len(initialsArr); i++ {
		if compareStudents(student, &initialsArr[i]) {
			return true, nil
		}
	}

	return false, nil
}

var students *Students
var guilds *Guilds

func main() {
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

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}
