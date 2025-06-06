package usergrpc

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/proto/userpb"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(client userpb.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserEmailByID(ctx context.Context, id int64) (string, error) {
	resp, err := u.client.GetUserByID(ctx, &userpb.GetUserByIDRequest{Id: id})
	if err != nil {
		return "", err
	}
	return resp.Email, nil
}
