package discordtransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func (c *Client) executeWebhook(ctx context.Context, id, credential string, body []byte, multipart bool) ([]byte, error) {
	if !snowflake(id) || !validToken(credential) || credential == "." || credential == ".." || !json.Valid(body) {
		return nil, errors.New("invalid Discord webhook request")
	}
	path := "/api/v10/webhooks/" + id + "/" + credential + "?wait=true&with_components=true"
	if multipart {
		return c.upload(ctx, path, "", body, false, true)
	}
	return c.request(ctx, http.MethodPost, path, "", body, true)
}
