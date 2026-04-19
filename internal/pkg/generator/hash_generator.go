package generator

import (
	"context"
	"log/slog"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/encoder"
)

type Generator interface {
	GenerateBatch(context.Context, int64) ([]string, error)
	GenerateBatchAndSave(c context.Context, quantity int64) ([]string, error)
}
type generator struct {
	rep     store.Store
	encoder encoder.Encoder
}

func NewGenerator(
	rep store.Store,
	encoder encoder.Encoder,
) Generator {
	return &generator{
		rep:     rep,
		encoder: encoder,
	}
}

func (g *generator) GenerateBatch(c context.Context, quantity int64) ([]string, error) {
	numbers, err := g.rep.GetUniqueNumbers(c, quantity)
	if err != nil {
		slog.Error("could not get unique numbers by quantity", "quantity", quantity, "error", err)
		return nil, err
	}

	hashes, err := g.encoder.Encode(numbers)
	if err != nil {
		return nil, err
	}

	return hashes, nil
}

func (g *generator) GenerateBatchAndSave(c context.Context, quantity int64) ([]string, error) {
	var hashes []string
	var err error

	hashes, err = g.GenerateBatch(c, quantity)
	if err != nil {
		return nil, err
	}

	if len(hashes) != 0 {
		_, err = g.rep.SaveAllHashes(c, hashes)
		if err != nil {
			slog.Error("could not save hashes", "error", err)
			return nil, err
		}
	}

	return hashes, nil
}
