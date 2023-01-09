package ctx_test

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/server/middleware"
	"github.com/luminosita/honeycomb/pkg/utils"
	"github.com/luminosita/honeycomb/pkg/validators/adapters"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

type Mock struct {
	t    *testing.T
	ctrl *gomock.Controller
	app  *fiber.App

	token  string
	secret string
}

func newMock(t *testing.T) (m *Mock) {
	m = &Mock{}
	m.t = t
	m.ctrl = gomock.NewController(t)

	m.secret = "hellomoto"
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "12345"
	claims["name"] = "Laza Lazic"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	m.token, _ = token.SignedString([]byte(m.secret))

	return
}

type a struct {
	A string `validate:"required"`
	B int    `validate:"gte=0,lte=130"`
}

type Handle = func(c *ctx.Ctx) error
type TestHandler struct {
	h Handle
}

func (t *TestHandler) Handle(c *ctx.Ctx) error {
	return t.h(c)
}

func setupTest(m *Mock, restricted ...bool) func() {
	if m == nil {
		panic("Mock not initialized")
	}

	m.app = fiber.New()

	if len(restricted) > 0 && restricted[0] {
		m.app.Use("/api", middleware.Protected(m.secret))
	}

	_ = utils.SetupRoute(m.app, "", &http.Route{
		Type: http.GET,
		Path: "/hello",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			return c.SendString("Hello, World!")
		}},
	})
	_ = utils.SetupRoute(m.app, "/api", &http.Route{
		Type: http.GET,
		Path: "/test",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			return c.SendString(fmt.Sprintf("Hello, %s!", c.Token.Claims.(jwt.MapClaims)["name"]))
		}},
	})
	_ = utils.SetupRoute(m.app, "", &http.Route{
		Type: http.GET,
		Path: "/error",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			return errors.New("Test Error")
		}},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/good",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			err := adapters.NewValidatorAdapter().Validate(&a{
				A: "v1",
				B: 123,
			})
			if err != nil {
				return err
			}
			return c.SendString("Hello, World!")
		}},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/bad",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			return adapters.NewValidatorAdapter().Validate(&a{})
		}},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/nil",
		Handler: &TestHandler{h: func(c *ctx.Ctx) error {
			return adapters.NewValidatorAdapter().Validate(nil)
		}},
	})

	return func() {
		defer utils.AssertPanic(m.t)
	}
}

func TestGood(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/hello", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 200, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"body\":\"Hello, World!\"}", string(bytes))
}

func TestError(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/error", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 500, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"error\":\"Test Error\"}", string(bytes))
}

func TestValidationGood(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/validation/good", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 200, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"body\":\"Hello, World!\"}", string(bytes))
}

func TestValidationBad(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/validation/bad", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 400, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"error\":[{\"failedField\":\"a.A\",\"tag\":\"required\",\"value\":\"\"}]}", string(bytes))
}

func TestValidationNil(t *testing.T) {
	m := newMock(t)
	defer setupTest(m)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/validation/nil", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 500, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"error\":\"Bad validation request: \\u003cnil\\u003e\"}", string(bytes))
}

func TestJwtPublic(t *testing.T) {
	m := newMock(t)
	defer setupTest(m, true)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/hello", nil)
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 200, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"body\":\"Hello, World!\"}", string(bytes))
}

func TestJwtProtectedGood(t *testing.T) {
	m := newMock(t)
	defer setupTest(m, true)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", m.token))
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 200, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"body\":\"Hello, Laza Lazic!\"}", string(bytes))
}

func TestJwtProtectedBad(t *testing.T) {
	m := newMock(t)
	defer setupTest(m, true)()

	// Create route with GET method for test
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", "dasdsadasdasda"))
	resp, _ := m.app.Test(req, 1)
	defer func() { _ = req.Body.Close() }()

	assert.Equal(t, 401, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"data\":null,\"message\":\"Invalid or expired JWT\",\"status\":\"error\"}", string(bytes))
}
