package discordtransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
)

type moderationRequest struct {
	Guild      string
	Actor      string
	Target     string
	Permission int
}

type guildRole struct {
	ID          string
	Position    int
	Permissions string
}

type guildMember struct {
	Roles []string
}

func (c *Client) readDiscord(ctx context.Context, path, token string, result any) error {
	data, err := c.Request(ctx, http.MethodGet, path, token, nil)
	if err != nil {
		return err
	}
	if json.Unmarshal(data, result) != nil {
		return errors.New("invalid Discord response")
	}
	return nil
}

func (c *Client) authorizeModeration(ctx context.Context, token string, body []byte) (string, error) {
	var request moderationRequest
	if json.Unmarshal(body, &request) != nil || !validToken(token) ||
		!snowflake(request.Guild) || !snowflake(request.Actor) || !snowflake(request.Target) ||
		request.Permission != 1 && request.Permission != 2 && request.Permission != 28 && request.Permission != 40 {
		return "unavailable", nil
	}
	if request.Actor == request.Target {
		return "self", nil
	}
	prefix := "/api/v10/guilds/" + request.Guild
	var guild struct {
		Owner string `json:"owner_id"`
		Roles []guildRole
	}
	var actor, target guildMember
	var bot struct{ ID string }
	work := []struct {
		path  string
		value any
	}{
		{prefix, &guild},
		{prefix + "/members/" + request.Actor, &actor},
		{prefix + "/members/" + request.Target, &target},
		{"/api/v10/users/@me", &bot},
	}
	var wg sync.WaitGroup
	failures := make(chan error, len(work))
	for _, read := range work {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := c.readDiscord(ctx, read.path, token, read.value); err != nil {
				failures <- err
			}
		}()
	}
	wg.Wait()
	if len(failures) != 0 {
		return "unavailable", <-failures
	}
	if !snowflake(guild.Owner) || !snowflake(bot.ID) {
		return "unavailable", nil
	}
	if request.Target == guild.Owner {
		return "owner", nil
	}
	if request.Target == bot.ID {
		return "self", nil
	}

	var botMember guildMember
	if err := c.readDiscord(ctx, prefix+"/members/"+bot.ID, token, &botMember); err != nil {
		return "unavailable", err
	}
	roles := make(map[string]guildRole, len(guild.Roles))
	for _, role := range guild.Roles {
		if !snowflake(role.ID) || role.Position < 0 {
			return "unavailable", nil
		}
		roles[role.ID] = role
	}
	everyone, ok := roles[request.Guild]
	if !ok {
		return "unavailable", nil
	}
	memberRank := func(member guildMember) (guildRole, uint64, bool) {
		if member.Roles == nil {
			return guildRole{}, 0, false
		}
		rank := everyone
		permissions, err := strconv.ParseUint(everyone.Permissions, 10, 64)
		if err != nil {
			return guildRole{}, 0, false
		}
		for _, id := range member.Roles {
			role, ok := roles[id]
			bits, err := strconv.ParseUint(role.Permissions, 10, 64)
			if !ok || err != nil {
				return guildRole{}, 0, false
			}
			if higherRole(role, rank, everyone.ID) {
				rank = role
			}
			permissions |= bits
		}
		return rank, permissions, true
	}
	actorRole, actorPermissions, actorOK := memberRank(actor)
	botRole, botPermissions, botOK := memberRank(botMember)
	targetRole, targetPermissions, targetOK := memberRank(target)
	if !actorOK || !botOK || !targetOK {
		return "unavailable", nil
	}
	const administrator = 1 << 3
	needed := uint64(1)<<request.Permission | administrator
	if request.Actor != guild.Owner && actorPermissions&needed == 0 {
		return "actor_permission", nil
	}
	if botPermissions&needed == 0 {
		return "bot_permission", nil
	}
	if request.Permission == 40 && targetPermissions&administrator != 0 {
		return "target_admin", nil
	}
	if request.Actor != guild.Owner && !higherRole(actorRole, targetRole, everyone.ID) {
		return "actor_role", nil
	}
	if !higherRole(botRole, targetRole, everyone.ID) {
		return "bot_role", nil
	}
	return "1", nil
}

func higherRole(a, b guildRole, everyone string) bool {
	if a.ID == b.ID || a.ID == everyone {
		return false
	}
	if b.ID == everyone {
		return true
	}
	if a.Position != b.Position {
		return a.Position > b.Position
	}
	aID, e1 := strconv.ParseUint(a.ID, 10, 64)
	bID, e2 := strconv.ParseUint(b.ID, 10, 64)
	return e1 == nil && e2 == nil && aID < bID
}
