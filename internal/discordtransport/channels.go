package discordtransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
)

type overwrite struct {
	ID    string `json:"id"`
	Type  int    `json:"type"`
	Allow string `json:"allow"`
	Deny  string `json:"deny"`
}

type channelLock struct {
	before overwrite
	after  overwrite
}

type channelLocks struct {
	mu        sync.Mutex
	snapshots map[string]channelLock
}

func (s *channelLocks) change(ctx context.Context, client *Client, id, token, action string) error {
	if !snowflake(id) || !validToken(token) || (action != "lock" && action != "unlock") {
		return errors.New("invalid Discord channel lock request")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	snapshot, saved := s.snapshots[id]
	if action == "unlock" && !saved {
		return errors.New("no channel lock snapshot exists in this launcher")
	}
	path := "/api/v10/channels/" + id
	var channel struct {
		Guild      string `json:"guild_id"`
		Type       *int
		Overwrites []overwrite `json:"permission_overwrites"`
	}
	if err := client.readDiscord(ctx, path, token, &channel); err != nil {
		return err
	}
	if !snowflake(channel.Guild) || channel.Type == nil || channel.Overwrites == nil {
		return errors.New("Discord channel lacks a guild or permission overwrites")
	}
	switch *channel.Type {
	case 0, 2, 5, 13, 15, 16: // Guild channels that can accept send permissions.
	default:
		return errors.New("channel type does not support the TRIPOL lock")
	}
	var current overwrite
	for _, item := range channel.Overwrites {
		if item.ID == channel.Guild && item.Type == 0 {
			current = item
			break
		}
	}
	if saved && snapshot.after.ID != channel.Guild {
		return errors.New("Discord channel lock refers to a different guild")
	}
	if action == "unlock" {
		switch current {
		case snapshot.before:
			delete(s.snapshots, id)
			return nil
		case snapshot.after:
			method := http.MethodPut
			var body []byte
			if snapshot.before.ID == "" {
				method = http.MethodDelete
			} else {
				body = snapshot.before.payload()
			}
			if _, err := client.Request(ctx, method, path+"/permissions/"+channel.Guild, token, body); err != nil {
				return err
			}
			delete(s.snapshots, id)
			return nil
		default:
			return errors.New("channel overwrite changed since lock; no permissions were restored")
		}
	}
	if saved {
		if current == snapshot.after {
			return nil
		}
		if current != snapshot.before {
			return errors.New("channel overwrite changed since lock")
		}
	} else {
		var allow, deny uint64
		if current.ID != "" {
			var err error
			allow, err = strconv.ParseUint(current.Allow, 10, 64)
			if err != nil {
				return errors.New("invalid Discord allow overwrite")
			}
			deny, err = strconv.ParseUint(current.Deny, 10, 64)
			if err != nil {
				return errors.New("invalid Discord deny overwrite")
			}
		}
		const sendMessages = 1 << 11
		after := overwrite{ID: channel.Guild, Type: 0, Allow: strconv.FormatUint(allow&^sendMessages, 10), Deny: strconv.FormatUint(deny|sendMessages, 10)}
		if current == after {
			return nil // An existing third-party lock is not ours to restore.
		}
		snapshot = channelLock{before: current, after: after}
		if s.snapshots == nil {
			s.snapshots = make(map[string]channelLock)
		}
		// Keep the snapshot on an ambiguous HTTP failure so the next call can
		// check what Discord actually applied.
		s.snapshots[id] = snapshot
	}
	body := snapshot.after.payload()
	_, err := client.Request(ctx, http.MethodPut, path+"/permissions/"+channel.Guild, token, body)
	return err
}

func (item overwrite) payload() []byte {
	body, _ := json.Marshal(struct {
		Type  int    `json:"type"`
		Allow string `json:"allow"`
		Deny  string `json:"deny"`
	}{item.Type, item.Allow, item.Deny})
	return body
}
