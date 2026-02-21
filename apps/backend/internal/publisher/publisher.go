package publisher

import (
	"context"
	"fmt"

	"github.com/pulkitbhatt/ikiru/internal/handler/dto"
)

type FakePublisher struct {
}

func (f *FakePublisher) Publish(ctx context.Context, payload dto.DueMonitor) {
	fmt.Printf("publishing job: %v\n\n", payload)
}
