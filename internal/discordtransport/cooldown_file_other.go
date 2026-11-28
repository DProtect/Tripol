//go:build !linux && !windows

package discordtransport

import (
	"errors"
	"os"
)

func lockCooldownFile(string) (*os.File, error) {
	return nil, errors.New("Discord cooldown persistence requires Linux or Windows")
}
