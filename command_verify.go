package checkpoint

import (
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

	school := i.ApplicationCommandData().Options[0].StringValue()
	firstName := strings.Title(i.ApplicationCommandData().Options[1].StringValue())
	lastName := strings.Title(i.ApplicationCommandData().Options[2].StringValue())
	grade := int(i.ApplicationCommandData().Options[3].IntValue())
	teacherName := strings.Title(i.ApplicationCommandData().Options[4].StringValue())
	StudentNumber := int(i.ApplicationCommandData().Options[5].IntValue())

	student := student.Student{
		School:         school,
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

	successRespond(s, i, "You have been verified!")

	err = c.session.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, g.VerifiedRole)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	if g.GradeRoles[grade-1] == "" {
		return
	}

	err = c.session.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, g.GradeRoles[grade-1])
	if err != nil {
		errorRespond(s, i, err)

		return
	}
}
