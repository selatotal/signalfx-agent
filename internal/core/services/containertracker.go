package services

import (
	"sync"

	"github.com/pkg/errors"
)

type EndpointsByID map[ID]Endpoint

func (e EndpointsByID) First() (Endpoint, error) {
	for id := range e {
		return e[id], nil
	}
	return nil, errors.New("no endpoints present")
}

func (e EndpointsByID) AddEndpoint(endpoint Endpoint) EndpointsByID {
	id := endpoint.Core().ID
	if e == nil {
		return EndpointsByID{
			id: endpoint,
		}
	}
	e[id] = endpoint
	return e
}

func (e EndpointsByID) RemoveEndpoint(endpoint Endpoint) {
	if e == nil {
		return
	}
	delete(e, endpoint.Core().ID)
}

// EndpointTracker is used to maintain the relationship between an endpoint's IP
// address (host) and the endpoint(s) that pertain to it.
type EndpointTracker struct {
	sync.Mutex
	hostToEndpoints map[string]EndpointsByID
}

func NewEndpointTracker() *EndpointTracker {
	return &EndpointTracker{
		hostToEndpoints: make(map[string]EndpointsByID),
	}
}

func (ct *EndpointTracker) EndpointAdded(endpoint Endpoint) {
	ct.Lock()
	defer ct.Unlock()

	host := endpoint.Core().Host
	if host == "" {
		return
	}

	ct.hostToEndpoints[host] = ct.hostToEndpoints[host].AddEndpoint(endpoint)
}

func (ct *EndpointTracker) EndpointRemoved(endpoint Endpoint) {
	ct.Lock()
	defer ct.Unlock()

	host := endpoint.Core().Host
	if host == "" {
		return
	}

	ct.hostToEndpoints[host].RemoveEndpoint(endpoint)
}
