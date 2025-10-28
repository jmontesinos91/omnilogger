package worker

import (
	"github.com/jmontesinos91/oevents"
)

// IRoutingStrategy interface
type IRoutingStrategy interface {
	Apply(event oevents.OmniViewEvent) error
}
