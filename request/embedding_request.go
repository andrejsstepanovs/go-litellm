package request

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EmbeddingMedia is a single multimodal embedding input element.
// At least one of the fields must be set.
// Image and audio references (file paths or URLs) are resolved by the LiteLLM server, not by the client.
type EmbeddingMedia struct {
	Text  string   `json:"text,omitempty"`
	Image []string `json:"image,omitempty"`
	Audio []string `json:"audio,omitempty"`
}

// Validate checks that the media element is usable.
func (m EmbeddingMedia) Validate() error {
	if strings.TrimSpace(m.Text) != "" {
		return nil
	}
	if hasNonEmptyEntry(m.Image) || hasNonEmptyEntry(m.Audio) {
		return nil
	}
	return fmt.Errorf("embedding media input must have at least one of text, image or audio set")
}

func hasNonEmptyEntry(entries []string) bool {
	for _, e := range entries {
		if strings.TrimSpace(e) != "" {
			return true
		}
	}
	return false
}

func nonEmptyEntries(entries []string) []string {
	filtered := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.TrimSpace(e) != "" {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// MarshalJSON implements json.Marshaler, dropping empty entries.
func (m EmbeddingMedia) MarshalJSON() ([]byte, error) {
	type alias struct {
		Text  string   `json:"text,omitempty"`
		Image []string `json:"image,omitempty"`
		Audio []string `json:"audio,omitempty"`
	}
	return json.Marshal(alias{
		Text:  strings.TrimSpace(m.Text),
		Image: nonEmptyEntries(m.Image),
		Audio: nonEmptyEntries(m.Audio),
	})
}

// EmbeddingInput is the `input` value of an embedding request.
// Exactly one variant must be set:
//   - Text: a single plain string
//   - Texts: a batch of plain strings
//   - Media: a single dict element, e.g. {"text": "...", "image": ["..."]}
//   - MediaList: a batch of dict elements
type EmbeddingInput struct {
	Text      string
	Texts     []string
	Media     EmbeddingMedia
	MediaList []EmbeddingMedia
}

// NewTextInput returns an input for a single plain string.
func NewTextInput(text string) EmbeddingInput {
	return EmbeddingInput{Text: text}
}

// NewTextsInput returns an input for a batch of plain strings.
func NewTextsInput(texts []string) EmbeddingInput {
	return EmbeddingInput{Texts: texts}
}

// NewMediaInput returns an input for a single dict element.
func NewMediaInput(media EmbeddingMedia) EmbeddingInput {
	return EmbeddingInput{Media: media}
}

// NewMediaListInput returns an input for a batch of dict elements.
func NewMediaListInput(media []EmbeddingMedia) EmbeddingInput {
	return EmbeddingInput{MediaList: media}
}

// Validate checks that exactly one input variant is set and it is non-empty.
func (i EmbeddingInput) Validate() error {
	variants := 0

	if strings.TrimSpace(i.Text) != "" {
		variants++
	}
	if len(i.Texts) > 0 {
		if !hasNonEmptyEntry(i.Texts) {
			return fmt.Errorf("embedding input texts must not be empty")
		}
		variants++
	}
	if i.Media.Text != "" || i.Media.Image != nil || i.Media.Audio != nil {
		if err := i.Media.Validate(); err != nil {
			return err
		}
		variants++
	}
	if len(i.MediaList) > 0 {
		for n, media := range i.MediaList {
			if err := media.Validate(); err != nil {
				return fmt.Errorf("embedding media input %d: %w", n, err)
			}
		}
		variants++
	}

	if variants == 0 {
		return fmt.Errorf("embedding input must not be empty")
	}
	if variants > 1 {
		return fmt.Errorf("embedding input must have exactly one variant set, got %d", variants)
	}
	return nil
}

// kind returns the marshaled JSON value for the input.
func (i EmbeddingInput) value() ([]byte, error) {
	switch {
	case strings.TrimSpace(i.Text) != "":
		return json.Marshal(i.Text)
	case len(i.Texts) > 0:
		return json.Marshal(nonEmptyEntries(i.Texts))
	case i.Media.Text != "" || i.Media.Image != nil || i.Media.Audio != nil:
		return json.Marshal(i.Media)
	case len(i.MediaList) > 0:
		return json.Marshal(i.MediaList)
	default:
		return nil, fmt.Errorf("embedding input must not be empty")
	}
}

// EmbeddingRequest is the payload of the /v1/embeddings endpoint.
type EmbeddingRequest struct {
	Model string `json:"model"`
	// Input is the embedding input. See EmbeddingInput for the supported shapes.
	Input EmbeddingInput `json:"-"`
	// Dimensions is the MRL truncation size. Known supported values: 128, 256, 512, 768.
	// When zero, the server default is used.
	Dimensions int `json:"dimensions,omitempty"`
	// PromptName selects a prompt template, e.g. "Retrieval-query", "Retrieval-document".
	// Only valid with plain string inputs (Text or Texts).
	PromptName string `json:"prompt_name,omitempty"`
}

// Validate checks the request before it is sent.
func (r EmbeddingRequest) Validate() error {
	if strings.TrimSpace(r.Model) == "" {
		return fmt.Errorf("embedding model must not be empty")
	}
	if err := r.Input.Validate(); err != nil {
		return err
	}
	if r.PromptName == "" {
		return nil
	}
	if r.Input.Text == "" && len(r.Input.Texts) == 0 {
		return fmt.Errorf("embedding prompt_name is only valid with plain string inputs")
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (r EmbeddingRequest) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	input, err := r.Input.value()
	if err != nil {
		return nil, err
	}

	type alias struct {
		Model      string          `json:"model"`
		Input      json.RawMessage `json:"input"`
		Dimensions int             `json:"dimensions,omitempty"`
		PromptName string          `json:"prompt_name,omitempty"`
	}
	return json.Marshal(alias{
		Model:      r.Model,
		Input:      input,
		Dimensions: r.Dimensions,
		PromptName: r.PromptName,
	})
}
