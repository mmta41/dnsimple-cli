package api

import (
	"context"
	"errors"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"golang.org/x/oauth2"
	"strconv"
)

var ErrUnAuthorized = errors.New("authorization failed")
var ErrForbidden = errors.New("forbidden operation")

type Client struct {
	*dnsimple.Client

	token     string
	accountId int64
	isLogged  bool
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) SetToken(token string) {
	c.token = token
	c.accountId = 0
	c.isLogged = false
}

func (c *Client) AccountId() int64 {
	return c.accountId
}
func (c *Client) AccountIdStr() string {
	return strconv.FormatInt(c.accountId, 10)
}

func (c *Client) SetAccountId(accountId int64) {
	c.accountId = accountId
}

func (c *Client) IsLogged() bool {
	return c.isLogged
}

func (c *Client) Login() ([]dnsimple.Account, error) {
	data, err := c.Client.Accounts.ListAccounts(context.Background(), &dnsimple.ListOptions{})
	if err != nil {
		return nil, err
	}
	if len(data.Data) > 0 {
		c.SetAccountId(data.Data[0].ID)
	}
	return data.Data, nil
}

func InitClient(token string, debug bool) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	client := &Client{Client: dnsimple.NewClient(tc)}
	if debug {
		client.BaseURL = "https://api.sandbox.dnsimple.com"
	}
	return client
}

func (c *Client) ConvertError(err error) error {
	return ConvertError(err)
}

func ConvertError(err error) error {
	var errResponse *dnsimple.ErrorResponse
	if errors.As(err, &errResponse) {
		if errResponse.HTTPResponse.StatusCode == 401 {
			return ErrUnAuthorized
		}
		if errResponse.HTTPResponse.StatusCode == 403 {
			return ErrForbidden
		}
	}
	return err
}
