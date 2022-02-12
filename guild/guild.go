package guild

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type GuildStore struct {
	filename string
	file     *os.File
}

type Guild struct {
	ID           string     `json:"id"`
	VerifiedRole string     `json:"verified_role"`
	GradeRoles   [12]string `json:"grade_roles"`
	PronounRoles []string   `json:"pronoun_roles"`
}

func NewGuilds(filename string) *GuildStore {
	return &GuildStore{
		filename: filename,
	}
}

func (g *GuildStore) Open() error {
	var err error
	g.file, err = os.OpenFile(g.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GuildStore) Close() error {
	return gs.file.Close()
}

func (gs *GuildStore) Write(guilds []Guild) error {
	b, err := json.Marshal(guilds)
	if err != nil {
		return err
	}

	err = gs.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = gs.file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = gs.file.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GuildStore) WriteOne(guild Guild) error {
	guilds, err := gs.Guilds()
	if err != nil {
		return err
	}

	n := false

	for i, g := range guilds {
		if g.ID == guild.ID {
			guilds[i] = guild
			n = true
			break
		}
	}

	if !n {
		guilds = append(guilds, guild)
	}

	return gs.Write(guilds)
}

func (gs *GuildStore) Guilds() ([]Guild, error) {
	_, err := gs.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(gs.file)
	if err != nil {
		return nil, err
	}

	var guilds []Guild
	err = json.Unmarshal(b, &guilds)

	return guilds, nil
}

func (gs *GuildStore) Guild(guildID string) (Guild, error) {
	guilds, err := gs.Guilds()
	if err != nil {
		return Guild{}, err
	}

	for _, g := range guilds {
		if g.ID == guildID {
			return g, nil
		}
	}

	return Guild{
		ID: guildID,
	}, nil
}
