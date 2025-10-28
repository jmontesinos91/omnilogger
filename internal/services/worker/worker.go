package worker

import (
	"context"
	"github.com/jmontesinos91/oevents"
)

// IWorker interface
type IWorker interface {
	Handle(ctx context.Context, event oevents.OmniViewEvent) error
}
