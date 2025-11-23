package request

import "github.com/andrejsstepanovs/go-litellm/models"

type Speech struct {
	Model models.ModelID `json:"model"`
	Input string         `json:"input"`
	Voice string         `json:"voice"`

	// Optional instructions for voice synthesis
	Instructions string `json:"instructions,omitempty"`

	// Optional Defaults to mp3
	// The format to audio in.
	// Supported formats are mp3, opus, aac, flac, wav, and pcm.
	// !!! Do not set for gemini tts
	ResponseFormat string `json:"response_format,omitempty"`

	// Optional Defaults to 1.0
	// The speed of the generated audio.
	// Select a value from 0.25 to 4.0. 1.0 is the default.
	Speed float64 `json:"speed,omitempty"`

	// Optional Defaults to audio
	// The format to stream the audio in.
	// Supported formats are sse and audio. sse is not supported for tts-1 or tts-1-hd.
	StreamFormat string `json:"stream_format,omitempty"`
}
