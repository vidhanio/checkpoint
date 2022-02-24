package checkpoint

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const ephemeralFlag = 1 << 6

const (
	green  = 0x57F287
	yellow = 0xFEE75C
	black  = 0x000000
	red    = 0xED4245
)

func contains[T comparable](ts []T, t T) bool {
	for _, t2 := range ts {
		if t == t2 {
			return true
		}
	}

	return false
}

func isManageRoles(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user permissions")
		return false
	}

	return perms&discordgo.PermissionManageRoles != 0
}

func ephemeralify(r *discordgo.InteractionResponse) *discordgo.InteractionResponse {
	r.Data.Flags |= ephemeralFlag
	return r
}

func deferred(r *discordgo.InteractionResponse) *discordgo.InteractionResponse {
	r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	return r
}

func response() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	}
}

func contentResponse(c string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c,
		},
	}
}

func embedResponse(es ...*discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: es,
		},
	}
}

func embed(title string, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       black,
	}
}

func successEmbed(m string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success",
		Description: m,
		Color:       green,
	}
}

func warningEmbed(m string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Warning",
		Description: m,
		Color:       yellow,
	}
}

func errorEmbed(err error) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: err.Error(),
		Color:       red,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "⚠️ Please report issues at https://github.com/vidhanio/checkpoint/issues",
		},
	}
}

func contentMessage(c string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: c,
	}
}

func embedMessage(es ...*discordgo.MessageEmbed) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: es,
	}
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, r *discordgo.InteractionResponse) {
	err := s.InteractionRespond(i.Interaction, r)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to respond to interaction")
	}
}

func successRespond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	respond(s, i, ephemeralify(embedResponse(successEmbed(m))))
}

func warningRespond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	respond(s, i, ephemeralify(embedResponse(warningEmbed(m))))
}

func errorRespond(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	respond(s, i, ephemeralify(embedResponse(errorEmbed(err))))
}
