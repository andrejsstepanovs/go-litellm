package response

import "strings"

type ToolResponses []ToolResponse

type ToolResponse struct {
	Type        string `json:"type"` // "text" for text responses, "image_url" for image URLs, etc.
	Text        string `json:"text"`
	Annotations any    `json:"annotations"`
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
