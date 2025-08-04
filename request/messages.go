package request

import (
	"fmt"
	"log"
	"strings"

	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/response"
)

type Messages []Message

type Message struct {
	Role       MessageRole      `json:"role"`
	Contents   MessageContents  `json:"content"`
	Name       string           `json:"name,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
	ToolCalls  common.ToolCalls `json:"tool_calls,omitempty"`
}

func (m *Messages) RemoveEmpty() {
	var filtered Messages
	for _, msg := range *m {
		if msg.Role == "" {
			log.Println("ALERT: Message with empty role found, skipping")
			continue
		}
		contents := make(MessageContents, 0)
		for _, c := range msg.Contents {
			if c.Type == "text" && c.Text == "" && c.ImageUrl == nil {
				log.Println("ALERT: Message content with empty text found, skipping")
				continue
			}
			contents = append(contents, c)
		}
		if len(contents) == 0 {
			log.Println("ALERT: Message with no valid contents found, skipping")
			continue
		}
		filtered = append(filtered, Message{
			Role:       msg.Role,
			Contents:   contents,
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
			ToolCalls:  msg.ToolCalls,
		})
	}
	*m = filtered
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s", m.Role, m.Contents)
}

func MessageImage(url string) ImageUrl {
	return ImageUrl{
		URL: url,
	}
}

func UserMessageImage(content string, imageUrl ImageUrl) Message {
	msg := make(MessageContents, 0, 2)
	if content != "" {
		msg = append(msg, MessageContent{
			Type: "text",
			Text: content,
		})
	}

	msg = append(msg, MessageContent{
		Type:     "image_url",
		ImageUrl: &imageUrl,
	})

	return UserMessage(msg)
}

func SystemMessageSimple(content string) Message {
	return SystemMessage(MessageContents{{
		Type: "text",
		Text: content,
	}})
}

func UserMessageSimple(content string) Message {
	return UserMessage(MessageContents{{
		Type: "text",
		Text: content,
	}})
}

func SystemMessage(content MessageContents) Message {
	return Message{
		Role:     ROLE_SYSTEM,
		Contents: content,
	}
}

func UserMessage(content MessageContents) Message {
	return Message{
		Role:     ROLE_USER,
		Contents: content,
	}
}

func ToolCallMessage(toolCall common.ToolCall, toolResponse response.ToolResponse) Message {
	txt := toolResponse.Text
	if txt == "" {
		txt = "-"
	}

	msg := Message{
		ToolCallID: toolCall.ID,
		Role:       ROLE_TOOL,
		Name:       toolCall.Function.Name,
		Contents: MessageContents{
			{
				Type: "text",
				Text: txt, // important, content text must be present
			},
		},
	}

	return msg
}

func AIMessage(msg response.ResponseMessage) Message {
	resp := Message{
		Role: MessageRole(msg.Role),
	}

	if len(msg.ToolCalls) > 0 {
		resp.ToolCalls = msg.ToolCalls.SortASC()
	}

	txt := msg.Content
	if txt == "" {
		txt = "-"
	}
	resp.Contents = MessageContents{
		{
			Type: "text",
			Text: txt,
		},
	}

	return resp
}

func AssistantMessageSimple(contents string) Message {
	if contents == "" {
		contents = "-"
	}

	return Message{
		Role: ROLE_ASSISTANT,
		Contents: MessageContents{
			{
				Type: "text",
				Text: contents,
			},
		},
	}
}

func AssistantMessage(contents MessageContents) Message {
	return Message{
		Role:     ROLE_ASSISTANT,
		Contents: contents,
	}
}

func (m *Messages) AddMessagePair(user, agent Message) {
	*m = append(*m, user, agent)
}

func (m *Messages) AddMessage(message Message) {
	*m = append(*m, message)
}

func (m *Messages) String() string {
	resp := make([]string, 0)
	for _, msg := range *m {
		resp = append(resp, msg.String())
	}
	return strings.Join(resp, "\n")
}
