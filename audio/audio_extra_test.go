package audio_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/audio"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestTranscribeAudio_WithExtraBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		assert.NoError(t, err)

		assert.Equal(t, "nova-2", r.FormValue("model"))
		assert.Equal(t, "true", r.FormValue("smart_format"))
		assert.Equal(t, "true", r.FormValue("punctuate"))
		assert.Equal(t, "some_value", r.FormValue("custom_field"))
		assert.Equal(t, "42", r.FormValue("number_field"))
		assert.Equal(t, "false", r.FormValue("bool_field"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"text":"hello world"}`))
	}))
	defer server.Close()

	audioFile := "testdata/file_174.oga"
	extraBody := map[string]any{
		"smart_format": true,
		"punctuate":    true,
		"custom_field": "some_value",
		"number_field": 42,
		"bool_field":   false,
	}

	resp, err := audio.TranscribeAudio(server.URL, "test-token", audioFile, "nova-2", extraBody, nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"text":"hello world"}`, string(body))
}

func TestTranscribeAudio_WithExtraHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "YourAppName", r.Header.Get("X-App-Name"))
		assert.Equal(t, "user_123", r.Header.Get("X-User-Id"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Empty(t, r.Header.Get("X-Empty-Value"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"text":"hello world"}`))
	}))
	defer server.Close()

	audioFile := "testdata/file_174.oga"
	extraHeaders := map[string]string{
		"  X-App-Name  ": "  YourAppName  ",
		"X-User-Id":       "user_123",
		"X-Empty-Value":   "   ",
	}

	resp, err := audio.TranscribeAudio(server.URL, "test-token", audioFile, "whisper-1", nil, extraHeaders)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSpeech_WithExtraHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "YourAppName", r.Header.Get("X-App-Name"))
		assert.Equal(t, "user_123", r.Header.Get("X-User-Id"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte("fake-audio-bytes"))
	}))
	defer server.Close()

	extraHeaders := map[string]string{
		"X-App-Name": "YourAppName",
		"X-User-Id":  "user_123",
	}

	resp, err := audio.Speech(server.URL, "test-token", request.Speech{
		Model: "tts-1",
		Input: "hello",
		Voice: "alloy",
	}, extraHeaders)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSpeech_WithoutExtraHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("X-App-Name"))
		assert.Empty(t, r.Header.Get("X-User-Id"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte("fake-audio-bytes"))
	}))
	defer server.Close()

	resp, err := audio.Speech(server.URL, "test-token", request.Speech{
		Model: "tts-1",
		Input: "hello",
		Voice: "alloy",
	}, nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
func TestTranscribeAudio_WithoutExtraBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		assert.NoError(t, err)

		assert.Equal(t, "whisper-1", r.FormValue("model"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"text":"hello world"}`))
	}))
	defer server.Close()

	audioFile := "testdata/file_174.oga"

	resp, err := audio.TranscribeAudio(server.URL, "test-token", audioFile, "whisper-1", nil, nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"text":"hello world"}`, string(body))
}
