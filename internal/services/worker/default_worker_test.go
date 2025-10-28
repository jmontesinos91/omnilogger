package worker

import (
	"context"
	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/ologs/logger"
	"testing"
)

func TestDefaultWorker_Handle(t *testing.T) {

	contextLogger := logger.NewContextLogger("OMNILOGGER", "test", logger.TextFormat)

	type args struct {
		event oevents.OmniViewEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Happy Path",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewDefaultWorker(contextLogger)
			err := w.Handle(context.Background(), tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultWorker.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
