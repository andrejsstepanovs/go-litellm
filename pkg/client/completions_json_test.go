package client_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/andrejsstepanovs/go-litellm/pkg/client"
	"github.com/andrejsstepanovs/go-litellm/pkg/json_schema"
	"github.com/andrejsstepanovs/go-litellm/pkg/request"
	"github.com/stretchr/testify/assert"
)

func TestNewCompletionRequest_WithJSONSchema_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		ctx := context.Background()

		messages := request.Messages{}
		messages.AddMessage(request.UserMessageSimple("List largest 3 cities!"))

		type City struct {
			CityName        string `json:"city_name"`
			PopulationCount int    `json:"population_count"`
		}
		type listOfCities struct {
			Cities []City `json:"cities"`
		}

		schema := json_schema.JSONSchema{
			Name: "list_of_cities",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cities": map[string]interface{}{
						"type":  "array",
						"title": "Cities",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"city_name": map[string]interface{}{
									"type":        "string",
									"description": "A city name",
								},
								"population_count": map[string]interface{}{
									"type":        "integer",
									"description": "Population count of the city",
								},
							},
							"required": []string{"city_name", "population_count"},
						},
					},
				},
				"required": []string{"cities"},
			},
			Strict: true,
		}

		modelMeta, err := clientInstance.Model(ctx, testModelGood)
		assert.NoError(t, err)

		var temp float32 = 0.2
		req := request.NewCompletionRequest(modelMeta, messages, request.LLMCallTools{}, &temp, 0.2)
		req.SetJSONSchema(schema)
		resp, err := clientInstance.Completion(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp)

		responseData := listOfCities{}
		err = json.Unmarshal(resp.Bytes(), &responseData)
		assert.NoError(t, err)

		assert.NotEmpty(t, responseData.Cities)
		assert.Equal(t, 3, len(responseData.Cities))
		for _, city := range responseData.Cities {
			assert.NotEmpty(t, city.CityName)
			assert.Greater(t, city.PopulationCount, 0)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		// Simulate a response missing the required "population_count"
		raw := []byte(`{"cities":[{"city_name":"London"}]}`)
		type City struct {
			CityName        string `json:"city_name"`
			PopulationCount int    `json:"population_count"`
		}
		type listOfCities struct {
			Cities []City `json:"cities"`
		}
		var responseData listOfCities
		err := json.Unmarshal(raw, &responseData)
		assert.NoError(t, err)
		assert.Equal(t, "London", responseData.Cities[0].CityName)
		assert.Equal(t, 0, responseData.Cities[0].PopulationCount) // zero value
	})

	t.Run("empty cities array", func(t *testing.T) {
		raw := []byte(`{"cities":[]}`)
		type City struct {
			CityName        string `json:"city_name"`
			PopulationCount int    `json:"population_count"`
		}
		type listOfCities struct {
			Cities []City `json:"cities"`
		}
		var responseData listOfCities
		err := json.Unmarshal(raw, &responseData)
		assert.NoError(t, err)
		assert.Empty(t, responseData.Cities)
	})

	t.Run("malformed json", func(t *testing.T) {
		raw := []byte(`{"cities":[{"city_name":"Paris","population_count":"not_a_number"}]}`)
		type City struct {
			CityName        string `json:"city_name"`
			PopulationCount int    `json:"population_count"`
		}
		type listOfCities struct {
			Cities []City `json:"cities"`
		}
		var responseData listOfCities
		err := json.Unmarshal(raw, &responseData)
		assert.Error(t, err)
	})
}
