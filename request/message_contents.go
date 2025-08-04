package request

import (
	"fmt"
	"strings"
)

type MessageContent struct {
	Type     string    `json:"type"`                // e.g. "text", "image_url", etc.
	Text     string    `json:"text,omitempty"`      // Text content
	ImageUrl *ImageUrl `json:"image_url,omitempty"` // Image content, if applicable
}

type MessageContents []MessageContent

func (mc MessageContents) String() string {
	resp := make([]string, 0, len(mc))
	for _, content := range mc {
		if content.Type == "text" {
			resp = append(resp, content.Text)
		} else if content.Type == "image_url" && content.ImageUrl != nil {
			resp = append(resp, fmt.Sprintf("[Image: %s]", content.ImageUrl.URL))
		}
	}
	return strings.Join(resp, " ")
}
