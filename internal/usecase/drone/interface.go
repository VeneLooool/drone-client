package drone

import (
	"context"
	
	"github.com/VeneLooool/drone-client/internal/model"
)

type Publisher interface {
	Publish(ctx context.Context, event model.Event) (err error)
}
