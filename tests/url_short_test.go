package tests

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/hisnameisivan/demo_url_short/internal/config"
	"github.com/hisnameisivan/demo_url_short/internal/http-server/handlers/url/save"
	api "github.com/hisnameisivan/demo_url_short/internal/lib/api/response"
	"github.com/stretchr/testify/require"
)

func TestUrlShort_SaveRedirect(t *testing.T) {
	cwd, _ := os.Getwd()
	os.Setenv("CONFIG_PATH", path.Join(cwd, "../config/local.yaml")) // как установить env другим способом?

	cfg := config.MustLoad()

	baseUrl := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	testCases := []struct {
		name  string
		url   string
		alias string
		err   string
	}{
		{
			name:  "Valid Url",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid Url",
			url:   "invalid",
			alias: gofakeit.Word(),
			err:   "field Url is not valid Url",
		},
		{
			name: "Empty alias",
			url:  gofakeit.URL(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &url.URL{
				Scheme: "http",
				Host:   baseUrl,
			}
			e := httpexpect.Default(t, u.String())

			// save

			resp := e.POST("/url").
				WithJSON(save.Request{
					Url:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth(cfg.User, cfg.Password).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.err != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.err)

				return
			}

			var alias string
			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
				alias = tc.alias
			} else {
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			// redirect

			fullUrl, _ := url.JoinPath(u.String(), alias)
			redirectedToUrl, err := api.GetRedirect(fullUrl)
			require.NoError(t, err)

			require.Equal(t, tc.url, redirectedToUrl)
		})
	}

}
