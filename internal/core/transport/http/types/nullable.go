package core_http_types

import (
	"encoding/json"

	"github.com/fffeehh/practice_backend_todoapp/internal/core/domain"
)

// делаем встраивание структуры Nullable из пакета domian в Nullable уровня transport,
// чтобы описать для нее метод декодирования входящего JSON
type Nullable[T any] struct {
	domain.Nullable[T]
}

// переопределяем метод UnmarshalJSON
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	// если данный метод вызван, значит в JSON точно было что-то передано
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set: n.Set,
	}
}
