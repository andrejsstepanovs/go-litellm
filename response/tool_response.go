package response

import (
	"encoding/json"
	"strings"
)

type ToolResponses []ToolResponse

type ToolResponse struct {
	Type        string `json:"type"` // "text" for text responses, "image_url" for image URLs, etc.
	Text        string `json:"text"`
	Annotations any    `json:"annotations"`
}

type ToolResponseWrapper struct {
	Meta              any               `json:"meta"`
	Content           []ToolContentItem `json:"content"`
	IsError           bool              `json:"isError"`
	StructuredContent any               `json:"structuredContent"`
}

type ToolContentItem struct {
	Meta        any    `json:"meta"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Annotations any    `json:"annotations"`
}

func (t *ToolResponses) UnmarshalJSON(data []byte) error {
	var wrapper ToolResponseWrapper
	if err := json.Unmarshal(data, &wrapper); err == nil {
		if wrapper.Content != nil {
			for _, item := range wrapper.Content {
				*t = append(*t, ToolResponse{
					Type:        item.Type,
					Text:        item.Text,
					Annotations: item.Annotations,
				})
			}
			return nil
		}
	}

	var tools []ToolResponse
	if err := json.Unmarshal(data, &tools); err == nil {
		*t = tools
		return nil
	}

	var tool ToolResponse
	if err := json.Unmarshal(data, &tool); err == nil {
		*t = []ToolResponse{tool}
		return nil
	}

	return json.Unmarshal(data, (*[]ToolResponse)(t))
}

func (i *ToolContentItem) String() string {
	return i.Text
}

func (t *ToolResponse) String() string {
	return t.Text
}

func (t *ToolResponses) String() string {
	var ret []string
	for _, r := range *t {
		ret = append(ret, r.String())
	}
	return strings.Join(ret, "\n")
}
