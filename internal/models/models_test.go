package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectPool(t *testing.T) {
	pool := NewObjectPool()
	assert.NotNil(t, pool)

	t.Run("RDAPResponse Pool", func(t *testing.T) {
		resp := pool.GetResponse()
		assert.NotNil(t, resp)
		assert.Empty(t, resp.ObjectClassName)
		assert.Empty(t, resp.Events)
		assert.Empty(t, resp.Notices)
		assert.Empty(t, resp.Links)

		// Add some data
		resp.ObjectClassName = "test"
		resp.Events = append(resp.Events, &Event{Action: "test"})
		resp.Notices = append(resp.Notices, &Notice{Title: "test"})
		resp.Links = append(resp.Links, &Link{Value: "test"})

		// Return to pool
		pool.PutResponse(resp)

		// Get another one
		resp2 := pool.GetResponse()
		assert.NotNil(t, resp2)
		assert.Empty(t, resp2.ObjectClassName)
		assert.Empty(t, resp2.Events)
		assert.Empty(t, resp2.Notices)
		assert.Empty(t, resp2.Links)
	})

	t.Run("Event Pool", func(t *testing.T) {
		event := pool.GetEvent()
		assert.NotNil(t, event)
		assert.Empty(t, event.Action)
		assert.Empty(t, event.Actor)
		assert.Empty(t, event.Date)

		event.Action = "test"
		event.Actor = "test"
		event.Date = "2023-01-01"

		pool.PutEvent(event)

		event2 := pool.GetEvent()
		assert.NotNil(t, event2)
		assert.Empty(t, event2.Action)
		assert.Empty(t, event2.Actor)
		assert.Empty(t, event2.Date)
	})

	t.Run("Notice Pool", func(t *testing.T) {
		notice := pool.GetNotice()
		assert.NotNil(t, notice)
		assert.Empty(t, notice.Title)
		assert.Empty(t, notice.Description)
		assert.Empty(t, notice.Links)

		notice.Title = "test"
		notice.Description = append(notice.Description, "test")
		notice.Links = append(notice.Links, &Link{Value: "test"})

		pool.PutNotice(notice)

		notice2 := pool.GetNotice()
		assert.NotNil(t, notice2)
		assert.Empty(t, notice2.Title)
		assert.Empty(t, notice2.Description)
		assert.Empty(t, notice2.Links)
	})

	t.Run("Link Pool", func(t *testing.T) {
		link := pool.GetLink()
		assert.NotNil(t, link)
		assert.Empty(t, link.Value)
		assert.Empty(t, link.Rel)
		assert.Empty(t, link.Href)
		assert.Empty(t, link.Type)

		link.Value = "test"
		link.Rel = "test"
		link.Href = "test"
		link.Type = "test"

		pool.PutLink(link)

		link2 := pool.GetLink()
		assert.NotNil(t, link2)
		assert.Empty(t, link2.Value)
		assert.Empty(t, link2.Rel)
		assert.Empty(t, link2.Href)
		assert.Empty(t, link2.Type)
	})
}

func BenchmarkObjectPool(b *testing.B) {
	pool := NewObjectPool()

	b.Run("Response Pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp := pool.GetResponse()
				resp.ObjectClassName = "test"
				resp.Events = append(resp.Events, &Event{Action: "test"})
				pool.PutResponse(resp)
			}
		})
	})

	b.Run("Event Pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				event := pool.GetEvent()
				event.Action = "test"
				pool.PutEvent(event)
			}
		})
	})

	b.Run("Notice Pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				notice := pool.GetNotice()
				notice.Title = "test"
				pool.PutNotice(notice)
			}
		})
	})

	b.Run("Link Pool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				link := pool.GetLink()
				link.Value = "test"
				pool.PutLink(link)
			}
		})
	})
}
