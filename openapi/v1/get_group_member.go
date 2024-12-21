package v1

import (
	"context"

	"github.com/2mf8/Go-QQ-SDK/dto"
)

func (o *openAPI) GetGroupMembers(ctx context.Context, groupId string, limit, startIndex int) (*dto.GetGroupMembersResp, error) {
	req := &dto.GetGroupMembersReq{
		Limit:      limit,
		StartIndex: startIndex,
	}
	resp, err := o.request(ctx).
		SetResult(dto.GetGroupMembersResp{}).
		SetPathParam("group_openid", groupId).
		SetBody(req).
		Post(o.getURL(groupMembersGet))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.GetGroupMembersResp), nil
}
