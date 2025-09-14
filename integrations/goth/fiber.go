package goth

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const ProviderParamKey key = iota

var (
	SessionStore  *session.Store
	ErrSessionNil = errors.New("goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store")
)

type Params struct {
	ctx *fiber.Ctx
}

func (p *Params) Get(key string) string {
	return p.ctx.Query(key)
}

type key int

func init() {
	config := session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s", gothic.SessionName),
		CookieHTTPOnly: true,
	}

	SessionStore = session.New(config)
}

func BeginAuthHandler(ctx *fiber.Ctx) error {
	url, err := GetAuthURL(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func SetState(ctx *fiber.Ctx) string {
	state := ctx.Query("state")
	if len(state) > 0 {
		return state
	}

	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func GetState(ctx *fiber.Ctx) string {
	return ctx.Query("state")
}

func GetAuthURL(ctx *fiber.Ctx) (string, error) {
	if SessionStore == nil {
		return "", ErrSessionNil
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)
	if err != nil {
		return "", err
	}

	return url, err
}

type CompleteUserAuthOptions struct {
	ShouldLogout bool
}

func CompleteUserAuth(ctx *fiber.Ctx, options ...CompleteUserAuthOptions) (goth.User, error) {
	if SessionStore == nil {
		return goth.User{}, ErrSessionNil
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := GetFromSession(providerName, ctx)
	if err != nil {
		return goth.User{}, err
	}

	shouldLogout := true
	if len(options) > 0 && !options[0].ShouldLogout {
		shouldLogout = false
	}

	if shouldLogout {
		defer Logout(ctx)
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	err = validateState(ctx, sess)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		return user, err
	}

	_, err = sess.Authorize(provider, &Params{ctx: ctx})
	if err != nil {
		return goth.User{}, err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)

	if err != nil {
		return goth.User{}, err
	}

	gu, err := provider.FetchUser(sess)
	return gu, err
}

func validateState(ctx *fiber.Ctx, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != ctx.Query("state")) {
		return errors.New("state token mismatch")
	}
	return nil
}

func Logout(ctx *fiber.Ctx) error {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return err
	}

	if err := session.Destroy(); err != nil {
		return err
	}

	return nil
}

func GetProviderName(ctx *fiber.Ctx) (string, error) {
	if p := ctx.Query("provider"); p != "" {
		return p, nil
	}

	if p := ctx.Params("provider"); p != "" {
		return p, nil
	}

	if p := ctx.Get("provider", ""); p != "" {
		return p, nil
	}

	if p := ctx.Get(fmt.Sprint(ProviderParamKey), ""); p != "" {
		return p, nil
	}

	providers := goth.GetProviders()
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return "", err
	}

	for _, provider := range providers {
		p := provider.Name()
		value := session.Get(p)
		if _, ok := value.(string); ok {
			return p, nil
		}
	}

	return "", errors.New("you must select a provider")
}

func GetContextWithProvider(ctx *fiber.Ctx, provider string) *fiber.Ctx {
	ctx.Set(fmt.Sprint(ProviderParamKey), provider)
	return ctx
}

func StoreInSession(key string, value string, ctx *fiber.Ctx) error {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return err
	}

	if err := updateSessionValue(session, key, value); err != nil {
		return err
	}

	session.Save()
	return nil
}

func GetFromSession(key string, ctx *fiber.Ctx) (string, error) {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return "", err
	}

	value, err := getSessionValue(session, key)
	if err != nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value, nil
}

func getSessionValue(store *session.Session, key string) (string, error) {
	value := store.Get(key)
	if value == nil {
		return "", errors.New("could not find a matching session for this request")
	}

	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func updateSessionValue(session *session.Session, key, value string) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	session.Set(key, b.String())

	return nil
}
