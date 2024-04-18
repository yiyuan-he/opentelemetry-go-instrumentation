// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/auto/internal/pkg/instrumentation/probe"
)

type EBPFSourceIDGenerator struct{}

type eBPFEventKey struct{}

func NewEBPFSourceIDGenerator() *EBPFSourceIDGenerator {
	return &EBPFSourceIDGenerator{}
}

// ContextWithEBPFEvent returns a copy of parent in which event is stored.
func ContextWithEBPFEvent(parent context.Context, event probe.SpanEvent) context.Context {
	return context.WithValue(parent, eBPFEventKey{}, event)
}

// EventFromContext returns the event within ctx if one exists.
func EventFromContext(ctx context.Context) *probe.SpanEvent {
	val := ctx.Value(eBPFEventKey{})
	if val == nil {
		return nil
	}

	event, ok := val.(probe.SpanEvent)
	if !ok {
		return nil
	}

	return &event
}

func (e *EBPFSourceIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	event := EventFromContext(ctx)
	if event == nil || event.SpanContext == nil {
		return trace.TraceID{}, trace.SpanID{}
	}

	return event.SpanContext.TraceID(), event.SpanContext.SpanID()
}

func (e *EBPFSourceIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	event := EventFromContext(ctx)
	if event == nil {
		return trace.SpanID{}
	}

	return event.SpanContext.SpanID()
}
