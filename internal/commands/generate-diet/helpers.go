package generatediet

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"

	"diet_bot/internal/entity"
)

// GenerateJSONSchema создает JSON Schema на основе структуры GeneratedDiet
func GenerateJSONSchema() (string, error) {
	// Создаем рефлектор для генерации схемы
	reflector := jsonschema.Reflector{
		// Настройки рефлектора
		ExpandedStruct: true,
		DoNotReference: true,
	}

	// Получаем схему из структуры
	schema := reflector.Reflect(&entity.GeneratedDiet{})

	// Преобразуем схему в JSON
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка при маршалинге схемы: %w", err)
	}

	return string(schemaJSON), nil
}
