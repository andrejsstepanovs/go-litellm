package client_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andrejsstepanovs/go-litellm/pkg/client"
	"github.com/andrejsstepanovs/go-litellm/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSpeechToText(t *testing.T) {
	for _, audioModelName := range []models.ModelID{testSTTOne, testSTTTwo} {
		t.Run(fmt.Sprintf("success %s", audioModelName), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := r.ParseMultipartForm(10 << 20)
				assert.NoError(t, err)
				file, _, err := r.FormFile("file")
				assert.NoError(t, err)
				defer func() {
					if err := file.Close(); err != nil {
						panic(err)
					}
				}()
				_, err = io.ReadAll(file)
				assert.NoError(t, err)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, err = w.Write([]byte(`{"text":"hello world", "words":[{"word": "me","start": 0.48,"end": 0.98},{"word": "a","start": 1.36,"end": 1.76}]}`))
				assert.NoError(t, err)
			}))
			defer server.Close()

			testUrl, err := url.Parse(server.URL)
			assert.NoError(t, err)

			c := getConfig()
			conn := getConn()
			conn.URL = *testUrl
			clientInstance := client.Litellm{Config: c, Connection: conn}

			file := "testdata/file_174.oga"
			res, err := clientInstance.SpeechToText(
				context.Background(),
				models.ModelMeta{ModelId: models.ModelID(audioModelName)},
				file,
			)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", res.Text)
			assert.Equal(t, 2, len(res.Words))
		})
	}

	t.Run("http error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad request", http.StatusBadRequest)
		}))
		defer server.Close()

		testUrl, err := url.Parse(server.URL)
		assert.NoError(t, err)

		c := getConfig()
		conn := getConn()
		conn.URL = *testUrl
		clientInstance := client.Litellm{Config: c, Connection: conn}

		file := "testdata/file_174.oga"
		_, err = clientInstance.SpeechToText(context.Background(), models.ModelMeta{ModelId: "whisper-1"}, file)
		assert.Error(t, err)
	})
}

func TestSpeechToText_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping functional test in short mode")
		return
	}

	testCases := []struct {
		name              string
		modelName         string
		expectedText      string
		expectedWordCount int
	}{
		{
			name:              "whisper-1",
			modelName:         string(testSTTOne),
			expectedText:      "Make me a story about Strocki the ostrich who met in friendly Puma.",
			expectedWordCount: 0,
		},
		{
			name:              "deepgram",
			modelName:         string(testSTTTwo),
			expectedText:      "me a story about storchy the ostrich who met in friendly poom",
			expectedWordCount: 12,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}

			file := "testdata/file_174.oga"
			res, err := clientInstance.SpeechToText(
				context.Background(),
				models.ModelMeta{ModelId: models.ModelID(tc.modelName)},
				file,
			)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedText, res.Text)
			assert.Equal(t, tc.expectedWordCount, len(res.Words))
		})
	}
}
