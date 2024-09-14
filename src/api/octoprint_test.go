package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marlinraker/src/api/notification"
	"marlinraker/src/config"
	"marlinraker/src/files"
	"marlinraker/src/marlinraker"
)

func testOctoPrintEndpoint(t *testing.T, endpoint string,
	f func(*testing.T, string)) {
	t.Run("GET"+endpoint, func(t *testing.T) {

		urlQuery, err := url.Parse(endpoint)
		if err != nil {
			t.Fatal(err)
		}

		request, err := http.NewRequestWithContext(context.Background(), "GET", urlQuery.String(), nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		HTTPHandler{}.ServeHTTP(recorder, request)

		f(t, recorder.Body.String())
	})
}

func TestOctoPrint(t *testing.T) {

	marlinraker.Config = config.DefaultConfig()
	marlinraker.Config.Misc.OctoprintCompat = true

	testOctoPrintEndpoint(t, "/api/version", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "server": "1.5.0",
            "api": "0.1",
            "text": "OctoPrint (Marlinraker)"
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/server", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "server": "1.5.0",
            "safemode": "settings"
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/login", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "_is_external_client": false,
            "_login_mechanism": "apikey",
            "name": "_api",
            "active": true,
            "user": true,
            "admin": true,
            "apikey": null,
            "permissions": [],
            "groups": ["admins", "users"]
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/settings", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "plugins": {
                "UltimakerFormatPackage": {
                    "align_inline_thumbnail": false,
                    "inline_thumbnail": false,
                    "inline_thumbnail_align_value": "left",
                    "inline_thumbnail_scale_value": "50",
                    "installed": true,
                    "installed_version": "0.2.2",
                    "scale_inline_thumbnail": false,
                    "state_panel_thumbnail": true
                }
            },
            "feature": {
                "sdSupport": false,
                "temperatureGraph": false
            },
            "webcam": {
                "flipH": false,
                "flipV": false,
                "rotate90": false,
                "streamUrl": "/webcam/?action=stream",
                "webcamEnabled": true
            }
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/job", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "job": {
                "file": {
                    "name": null
                },
                "estimatedPrintTime": null,
                "filament": {
                    "length": null
                },
                "user": null
            },
            "progress": {
                "completion": null,
                "filepos": null,
                "printTime": null,
                "printTimeLeft": null,
                "printTimeOrigin": null
            },
            "state": "Offline"
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/printer", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "temperature": {
                "tool0": {
                    "actual": 22.25,
                    "offset": 0,
                    "target": 0
                },
                "bed": {
                    "actual": 22.25,
                    "offset": 0,
                    "target": 0
                }
            },
            "state": {
                "text": "state",
                "flags": {
                    "operational": true,
                    "paused": false,
                    "printing": false,
                    "cancelling": false,
                    "pausing": false,
                    "error": false,
                    "ready": false,
                    "closedOrError": false
                }
            }
        }
		`, result)
	})

	testOctoPrintEndpoint(t, "/api/printerprofiles", func(t *testing.T, result string) {
		assert.JSONEq(t, `
		{
            "profiles": {
                "_default": {
                    "id": "_default",
                    "name": "Default",
                    "color": "default",
                    "model": "Default",
                    "default": true,
                    "current": true,
                    "heatedBed": true,
                    "heatedChamber": false
                }
            }
        }
		`, result)
	})

	// TODO: create MockPrinter and test this
	// testOctoPrintEndpoint(t, "POST", "/api/printer/command", []byte(`{"commands":["M115","G28"]}`), func(t *testing.T, result string) {
	//     assert.Equal(t, "{}", result)
	// })
}

func TestOctoPrintFileUpload(t *testing.T) {

	var err error
	files.DataDir, err = filepath.Abs("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	files.Fs = afero.NewCopyOnWriteFs(afero.NewReadOnlyFs(afero.NewOsFs()), afero.NewMemMapFs())
	err = files.Init()
	if err != nil {
		t.Fatal(err)
	}

	marlinraker.State = marlinraker.Ready
	marlinraker.Config = config.DefaultConfig()
	notification.Testing = true

	testFileUpload(t, "/api/files/local", map[string]string{
		"root": "config",
		"path": "test/path",
	}, "file.txt", "testdata/test_upload.txt",
		func(t *testing.T, _ *httptest.ResponseRecorder, _ *any, error *Error) {

			if error != nil {
				t.Fatal(error)
			}

			filePath, err := filepath.Abs("testdata/config/test/path/file.txt")
			if err != nil {
				t.Fatal(err)
			}

			contents, err := afero.ReadFile(files.Fs, filePath)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "this is a test\n", string(contents))
		})
}
