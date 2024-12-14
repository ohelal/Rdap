// models.go
package models

import (
    "sync"
)

// ObjectPool manages pools of frequently used objects
type ObjectPool struct {
    responses sync.Pool
    events    sync.Pool
    notices   sync.Pool
    links     sync.Pool
    entities  sync.Pool
}

// NewObjectPool creates a new object pool
func NewObjectPool() *ObjectPool {
    return &ObjectPool{
        responses: sync.Pool{
            New: func() interface{} {
                return &RDAPResponse{
                    Events:  make([]*Event, 0, 2),
                    Notices: make([]*Notice, 0, 2),
                    Links:   make([]*Link, 0, 2),
                    Entities: make([]Entity, 0, 4),
                }
            },
        },
        events: sync.Pool{
            New: func() interface{} {
                return &Event{}
            },
        },
        notices: sync.Pool{
            New: func() interface{} {
                return &Notice{
                    Description: make([]string, 0, 2),
                    Links:       make([]*Link, 0, 2),
                }
            },
        },
        links: sync.Pool{
            New: func() interface{} {
                return &Link{}
            },
        },
        entities: sync.Pool{
            New: func() interface{} {
                return &Entity{
                    VCardArray: make([]interface{}, 0, 8),
                    Roles:      make([]string, 0, 2),
                    Events:     make([]*Event, 0, 2),
                    Links:      make([]*Link, 0, 2),
                    Status:     make([]string, 0, 2),
                }
            },
        },
    }
}

// GetResponse gets a response object from the pool
func (p *ObjectPool) GetResponse() *RDAPResponse {
    return p.responses.Get().(*RDAPResponse)
}

// PutResponse returns a response object to the pool
func (p *ObjectPool) PutResponse(r *RDAPResponse) {
    if r == nil {
        return
    }
    // Clear slices without deallocating
    r.Events = r.Events[:0]
    r.Notices = r.Notices[:0]
    r.Links = r.Links[:0]
    r.Entities = r.Entities[:0]
    // Clear string fields
    r.ObjectClassName = ""
    r.Handle = ""
    r.StartAddress = ""
    r.EndAddress = ""
    r.IPVersion = ""
    r.Name = ""
    r.Type = ""
    r.Country = ""
    r.Status = nil
    r.Port43 = ""
    p.responses.Put(r)
}

// GetEvent gets an event object from the pool
func (p *ObjectPool) GetEvent() *Event {
    return p.events.Get().(*Event)
}

// PutEvent returns an event object to the pool
func (p *ObjectPool) PutEvent(e *Event) {
    if e == nil {
        return
    }
    e.Action = ""
    e.Actor = ""
    e.Date = ""
    p.events.Put(e)
}

// GetNotice gets a notice object from the pool
func (p *ObjectPool) GetNotice() *Notice {
    return p.notices.Get().(*Notice)
}

// PutNotice returns a notice object to the pool
func (p *ObjectPool) PutNotice(n *Notice) {
    if n == nil {
        return
    }
    n.Title = ""
    n.Description = n.Description[:0]
    n.Links = n.Links[:0]
    p.notices.Put(n)
}

// GetLink gets a link object from the pool
func (p *ObjectPool) GetLink() *Link {
    return p.links.Get().(*Link)
}

// PutLink returns a link object to the pool
func (p *ObjectPool) PutLink(l *Link) {
    if l == nil {
        return
    }
    l.Value = ""
    l.Rel = ""
    l.Href = ""
    l.Type = ""
    p.links.Put(l)
}

// GetEntity gets an entity object from the pool
func (p *ObjectPool) GetEntity() *Entity {
    return p.entities.Get().(*Entity)
}

// PutEntity returns an entity object to the pool
func (p *ObjectPool) PutEntity(e *Entity) {
    if e == nil {
        return
    }
    e.ObjectClassName = ""
    e.Handle = ""
    e.VCardArray = e.VCardArray[:0]
    e.Roles = e.Roles[:0]
    e.Events = e.Events[:0]
    e.Links = e.Links[:0]
    e.Status = e.Status[:0]
    p.entities.Put(e)
}

// RDAPResponse represents an RDAP response
type RDAPResponse struct {
    ObjectClassName string    `json:"objectClassName,omitempty"`
    Handle         string    `json:"handle,omitempty"`
    StartAddress   string    `json:"startAddress,omitempty"`
    EndAddress     string    `json:"endAddress,omitempty"`
    IPVersion      string    `json:"ipVersion,omitempty"`
    Name           string    `json:"name,omitempty"`
    Type           string    `json:"type,omitempty"`
    Country        string    `json:"country,omitempty"`
    Status         []string  `json:"status,omitempty"`
    Events         []*Event  `json:"events,omitempty"`
    Entities       []Entity  `json:"entities,omitempty"`
    Notices        []*Notice `json:"notices,omitempty"`
    Links          []*Link   `json:"links,omitempty"`
    Port43         string    `json:"port43,omitempty"`
    Remarks        []string  `json:"remarks,omitempty"`
}

// Event represents an RDAP event
type Event struct {
    Action string `json:"eventAction,omitempty"`
    Actor  string `json:"eventActor,omitempty"`
    Date   string `json:"eventDate,omitempty"`
}

// Notice represents an RDAP notice
type Notice struct {
    Title       string   `json:"title,omitempty"`
    Description []string `json:"description,omitempty"`
    Links       []*Link  `json:"links,omitempty"`
}

// Link represents an RDAP link
type Link struct {
    Value string `json:"value,omitempty"`
    Rel   string `json:"rel,omitempty"`
    Href  string `json:"href,omitempty"`
    Type  string `json:"type,omitempty"`
}

// Entity represents an RDAP entity
type Entity struct {
    ObjectClassName string     `json:"objectClassName"`
    Handle         string     `json:"handle,omitempty"`
    VCardArray     []interface{} `json:"vcardArray,omitempty"`
    Roles          []string   `json:"roles,omitempty"`
    Events         []*Event    `json:"events,omitempty"`
    Links          []*Link     `json:"links,omitempty"`
    Status         []string   `json:"status,omitempty"`
}