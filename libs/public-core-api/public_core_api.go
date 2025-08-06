// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package public_core_api

import (
	"context"
	"errors"
	"github.com/staringfun/millsmess/libs/types"
	"net/http"
	"net/url"
)

type PublicCoreAPIEngine interface {
	Fetch(path string, payload any, reqHeaders map[string]string, res any, ctx context.Context) (
		status int,
		resHeaders map[string]string,
		err error,
	)
}

type PublicCoreAPI struct {
	Engine  PublicCoreAPIEngine
	BaseURL string
}

type FetchMeArgs struct {
	Token  string
	Cookie string
}

type FetchMeResponse struct {
	Data types.User `json:"data"`
}

var FetchMeError = errors.New("fetch me")
var UnauthorizedError = errors.New("unauthorized")

func (c *PublicCoreAPI) FetchMeHeaders(args FetchMeArgs, headers map[string]string, ctx context.Context) (types.User, map[string]string, error) {
	if headers == nil {
		headers = map[string]string{}
	}
	if args.Token != "" {
		headers["Authorization"] = "Bearer " + args.Token
	}
	if args.Cookie != "" {
		headers["Cookie"] = args.Cookie
	}
	p, _ := url.JoinPath(c.BaseURL, "v1/users/me")
	var res FetchMeResponse
	status, resHeaders, err := c.Engine.Fetch(p, nil, headers, &res, ctx)
	if err != nil {
		return nil, nil, err
	}
	switch {
	case status == http.StatusUnauthorized:
		return nil, resHeaders, UnauthorizedError
	case status != http.StatusOK || res.Data == nil:
		return nil, resHeaders, FetchMeError
	default:
		return res.Data, resHeaders, nil
	}
}

func (c *PublicCoreAPI) FetchMe(args FetchMeArgs, ctx context.Context) (types.User, error) {
	res, _, err := c.FetchMeHeaders(args, nil, ctx)
	return res, err
}
