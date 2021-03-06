package checkpoint

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/checkpoint/student"
)

func (c *Checkpoint) verify(s *discordgo.Session, i *discordgo.InteractionCreate) {
	g, err := c.guilds.Guild(i.GuildID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}
	if g.VerifiedRole == "" {
		warningRespond(s, i, "Checkpoint is not configured. Please ask someone with `Manage Roles` permissions to use `/config set verified_role`.")

		return
	}

	firstName := strings.Title(i.ApplicationCommandData().Options[0].StringValue())
	lastName := strings.Title(i.ApplicationCommandData().Options[1].StringValue())
	grade := int(i.ApplicationCommandData().Options[2].IntValue())
	teacherName := strings.Title(i.ApplicationCommandData().Options[3].StringValue())
	StudentNumber := int(i.ApplicationCommandData().Options[4].IntValue())

	student := student.Student{
		Initials:       [2]rune{rune(firstName[0]), rune(lastName[0])},
		Grade:          grade,
		TeacherInitial: rune(teacherName[0]),
		StudentNumber:  StudentNumber,
	}

	v := c.students.Verify(student)

	if !v {
		warningRespond(s, i, "You have not been verified.")

		return
	}

	err = c.session.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, g.VerifiedRole)
	if err != nil {
		errorRespond(s, i, err)
	}

	successRespond(s, i, "You have been verified!")

	err = c.session.GuildMemberNickname(i.GuildID, i.Member.User.ID, fmt.Sprintf("%s %c.", firstName, lastName[0]))
	if err != nil {
		errorRespond(s, i, err)
	}

	if g.GradeRoles[grade-1] == "" {
		return
	}

	err = c.session.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, g.GradeRoles[grade-1])
	if err != nil {
		errorRespond(s, i, err)
	}
}
