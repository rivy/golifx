package common

import (
	"time"

	"github.com/satori/go.uuid"
)

const subscriptionChanSize = 16

// SubscriptionTarget defines the interface between a subscription and it's
// target object
type SubscriptionTarget interface {
	NewSubscription() (*Subscription, error)
	CloseSubscription(*Subscription) error
}

// Subscription exposes an event channel for consumers, and attaches to a
// SubscriptionTarget, that will feed it with events
type Subscription struct {
	events chan interface{}
	id     uuid.UUID
	target SubscriptionTarget
}

// ID returns the unique ID for this subscription
func (s *Subscription) ID() string {
	return s.id.String()
}

// Events returns a chan reader for reading events published to this
// subscription
func (s *Subscription) Events() <-chan interface{} {
	return s.events
}

// Write pushes an event onto the events channel
func (s *Subscription) Write(event interface{}) error {
	timeout := time.After(DefaultTimeout)
	select {
	case s.events <- event:
		return nil
	case <-timeout:
		return ErrTimeout
	}
}

// Close cleans up resources and notifies the target that the subscription
// should no longer be used.  It is important to close subscriptions when you
// are done with them to avoid blocking operations.
func (s *Subscription) Close() {
	close(s.events)
	s.target.CloseSubscription(s)
}

// NewSubscription returns a *Subscription attached to the specified target
func NewSubscription(target SubscriptionTarget) *Subscription {
	return &Subscription{
		events: make(chan interface{}, subscriptionChanSize),
		id:     uuid.NewV4(),
		target: target,
	}
}