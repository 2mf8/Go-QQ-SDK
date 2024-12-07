package v1

import (
	"context"

	"github.com/2mf8/Go-QQ-SDK/dto"
)

// Guild 拉取频道信息
func (o *openAPI) Guild(ctx context.Context, guildID string) (*dto.Guild, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Guild{}).
		SetPathParam("guild_id", guildID).
		Get(o.getURL(guildURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Guild), nil
}
