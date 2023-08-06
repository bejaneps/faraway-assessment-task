package log

import (
	"context"
	"sync"
)

// Attributes to store the context log attributes
type Attributes map[string]interface{}

type contextKey struct{}

type contextualAttributes struct {
	guard sync.Mutex
	attrs Attributes
}

// FromContext Returns the logger with contextual attributes if any
func FromContext(ctx context.Context) Logger {
	if ctxAttrs, ok := fromContext(ctx); ok {
		fields := make([]Field, 0, len(ctxAttrs.attrs))

		for attr, val := range ctxAttrs.attrs {
			fields = append(fields, Any(attr, val))
		}

		return With(fields...)
	}

	return Base()
}

// fromContext Returns the attributes stored in the context
func fromContext(ctx context.Context) (*contextualAttributes, bool) {
	if attrs, ok := ctx.Value(contextKey{}).(*contextualAttributes); ok {
		return attrs, true
	}
	return nil, false
}

// contextWithAttributes Returns the context with the given attributes
func contextWithAttributes(ctx context.Context, attrs Attributes) context.Context {
	return context.WithValue(ctx, contextKey{}, &contextualAttributes{attrs: attrs})
}

// ContextWithAttributes Upsert operation for the list of key-value pairs in the context and returning it
func ContextWithAttributes(ctx context.Context, attrs Attributes) context.Context {
	if ctxAttrs, ok := fromContext(ctx); ok {
		ctxAttrs.guard.Lock()
		defer ctxAttrs.guard.Unlock()

		// Updating/inserting the attribute in the context using its pointer
		for k, v := range attrs {
			if v == nil {
				continue
			}
			ctxAttrs.attrs[k] = v
		}

		return ctx
	}

	// In order to not share the same pointer, lets convert the attributes to new map
	var newAttrs = make(Attributes, len(attrs))
	for k, v := range attrs {
		if v == nil {
			continue
		}
		newAttrs[k] = v
	}

	return contextWithAttributes(ctx, newAttrs)
}
