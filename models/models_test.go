package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/models"
)

func Test_Get_Unit(t *testing.T) {
	model1 := models.Model{
		ID:      "model-1",
		Object:  "object-1",
		OwnedBy: "owner-1",
	}
	model2 := models.Model{
		ID:      "model-2",
		Object:  "object-2",
		OwnedBy: "owner-2",
	}

	tests := []struct {
		name          string
		models        models.Models
		id            models.ModelID
		expectedModel models.Model
		expectedFound bool
	}{
		{
			name:          "empty models slice",
			models:        models.Models{},
			id:            "model-1",
			expectedModel: models.Model{},
			expectedFound: false,
		},
		{
			name:          "model found in non-empty slice",
			models:        models.Models{model1, model2},
			id:            "model-2",
			expectedModel: model2,
			expectedFound: true,
		},
		{
			name:          "model not found in non-empty slice",
			models:        models.Models{model1, model2},
			id:            "model-3",
			expectedModel: models.Model{},
			expectedFound: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualModel, actualFound := tc.models.Get(tc.id)
			assert.Equal(t, tc.expectedModel, actualModel)
			assert.Equal(t, tc.expectedFound, actualFound)
		})
	}
}
