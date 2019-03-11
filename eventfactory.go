// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ycq

import (
	"fmt"
)

// EventMaker is the interface that an event factory should implement.
//
// An event factory returns instances of an event given the event type
// as a string.
// An event factory is required during deserialisation of events by the
// domain or Repository depending on your implementation.
//
// The domain will return a string describing the event type. To unmarshal
// the contents of the persisted event which will typically be in some serialised
// format such as JSON an instance of the event type will need to be created.
type EventMaker interface {
	MakeEvent(string) interface{}
}

// DelegateEventFactory uses delegate functions to instantiate event instances
// given the name of the event type as a string.
type DelegateEventFactory struct {
	eventFactories map[string]func() interface{}
}

// NewDelegateEventFactory constructs a new DelegateEventFactory
func NewDelegateEventFactory() *DelegateEventFactory {
	return &DelegateEventFactory{
		eventFactories: make(map[string]func() interface{}),
	}
}

// RegisterDelegate registers a delegate that will return an event instance given
// an event type name as a string.
//
// If an attempt is made to register multiple delegates for an event type, an error
// is returned.
func (t *DelegateEventFactory) RegisterDelegate(event interface{}, delegate func() interface{}) error {
	typeName := typeOf(event)
	if _, ok := t.eventFactories[typeName]; ok {
		return fmt.Errorf("Factory delegate already registered for type: \"%s\"", typeName)
	}
	t.eventFactories[typeName] = delegate
	return nil
}

// MakeEvent returns an event instance given an event type as a string.
//
// An appropriate delegate must be registered for the event type.
// If an appropriate delegate is not registered, the method will return nil.
func (t *DelegateEventFactory) MakeEvent(typeName string) interface{} {
	if f, ok := t.eventFactories[typeName]; ok {
		return f()
	}
	return nil
}
