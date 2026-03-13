package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/YukiAminaka/cycle-route-backend/config"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
)

type KratosMiddleware struct {
	ory *ory.APIClient
	conf *config.Config
}

func NewMiddleware(conf *config.Config) *KratosMiddleware {
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: conf.Server.KratosPublicUrl, // Kratos Public API
		},
	}
	return &KratosMiddleware{
		ory: ory.NewAPIClient(configuration),
		conf: conf,
	}
}

func (k *KratosMiddleware) Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := k.validateSession(c.Request)
		redirectURL := fmt.Sprintf("%s/login", k.conf.Server.FrontendOrigin)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, redirectURL)
			return
		}
		if !*session.Active {
			c.Redirect(http.StatusMovedPermanently, redirectURL)
			return
		}

		// セッション情報をコンテキストに保存
		c.Set("session", session)
		// KratosのユーザーIDをコンテキストに保存
		c.Set("kratos_id", session.Identity.Id)

		c.Next()
	}
}

func (k *KratosMiddleware) validateSession(r *http.Request) (*ory.Session, error) {
	cookie, err := r.Cookie("ory_kratos_session")
	if err != nil {
		return nil, err
	}
	if cookie == nil {
		return nil, errors.New("no session found in cookie")
	}
	decoded, err := url.QueryUnescape(cookie.Value)// クッキーに=や;などの特殊記号が含まれている場合に備えてデコード。
	if err != nil {
		return nil, err
	}
	decoded = "ory_kratos_session=" + decoded
	resp, _, err := k.ory.FrontendAPI.ToSession(context.Background()).Cookie(decoded).Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return resp, nil
}
