package ctx_test

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/http/ctx"
	"github.com/luminosita/honeycomb/pkg/server/utils"
	"github.com/luminosita/honeycomb/pkg/util"
	"github.com/luminosita/honeycomb/pkg/validators/adapters"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

type Mock struct {
	t    *testing.T
	ctrl *gomock.Controller
	app  *fiber.App
}

func newMock(t *testing.T) (m *Mock) {
	m = &Mock{}
	m.t = t
	m.ctrl = gomock.NewController(t)

	return
}

type a struct {
	A string `validate:"required"`
	B int    `validate:"gte=0,lte=130"`
}

func setupTest(m *Mock) func() {
	if m == nil {
		panic("Mock not initialized")
	}

	m.app = fiber.New()

	m.app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	_ = utils.SetupRoute(m.app, "", &http.Route{
		Type: http.GET,
		Path: "/hello",
		Handler: func(c *ctx.Ctx) error {
			return c.SendString("Hello, World!")
		},
	})
	_ = utils.SetupRoute(m.app, "", &http.Route{
		Type: http.GET,
		Path: "/error",
		Handler: func(c *ctx.Ctx) error {
			return errors.New("Test Error")
		},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/good",
		Handler: func(c *ctx.Ctx) error {
			err := adapters.NewValidatorAdapter().Validate(&a{
				A: "v1",
				B: 123,
			})
			if err != nil {
				return err
			}
			return c.SendString("Hello, World!")
		},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/bad",
		Handler: func(c *ctx.Ctx) error {
			return adapters.NewValidatorAdapter().Validate(&a{})
		},
	})
	_ = utils.SetupRoute(m.app, "/validation", &http.Route{
		Type: http.GET,
		Path: "/nil",
		Handler: func(c *ctx.Ctx) error {
			return adapters.NewValidatorAdapter().Validate(nil)
		},
	})

	return func() {
		defer util.AssertPanic(m.t)
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
