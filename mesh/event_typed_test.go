// Tideland Go Cells - Mesh - Typed Events Tests
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"
)

//--------------------
// TESTS
//--------------------

// TestNewEventTyped verifies type-safe event creation.
func TestNewEventTyped(t *testing.T) {
	type UserData struct {
		ID    string
		Email string
		Age   int
	}

	userData := UserData{
		ID:    "user-123",
		Email: "test@example.com",
		Age:   30,
	}

	evt, err := NewEventTyped("user-created", userData,
		WithTraceID("trace-abc"),
		WithCorrelation("session-xyz"),
		WithCustom("source", "api"),
	)

	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "user-created")
	verify.True(t, evt.HasPayload())

	// Verify metadata
	meta := evt.Metadata()
	verify.NotEmpty(t, meta.ID)
	verify.Equal(t, meta.TraceID, "trace-abc")
	verify.Equal(t, meta.Correlation, "session-xyz")
	verify.Equal(t, meta.Custom["source"], "api")
}

// TestPayloadAs verifies type-safe payload extraction.
func TestPayloadAs(t *testing.T) {
	type Product struct {
		Name  string
		Price float64
		Stock int
	}

	product := Product{
		Name:  "Widget",
		Price: 19.99,
		Stock: 100,
	}

	evt, err := NewEventTyped("product-added", product)
	verify.NoError(t, err)

	// Extract with type safety
	extracted, err := PayloadAs[Product](evt)
	verify.NoError(t, err)
	verify.Equal(t, extracted.Name, "Widget")
	verify.Equal(t, extracted.Price, 19.99)
	verify.Equal(t, extracted.Stock, 100)
}

// TestPayloadAsWrongType verifies type extraction with incompatible data.
func TestPayloadAsWrongType(t *testing.T) {
	// Create event with string payload
	evt, err := NewEventTyped("test", "this is a string")
	verify.NoError(t, err)

	// Try to extract as int - JSON will fail to unmarshal string to int
	type NumberData struct {
		Value int
	}
	_, err = PayloadAs[NumberData](evt)
	verify.Error(t, err)
	verify.Match(t, err.Error(), ".*cannot unmarshal.*")
}

// TestNewEventWithMetadata verifies non-generic metadata support.
func TestNewEventWithMetadata(t *testing.T) {
	evt, err := NewEventWithMetadata("test-event", map[string]string{
		"key": "value",
	},
		WithTraceID("trace-123"),
		WithCorrelation("corr-456"),
		WithCustom("env", "testing"),
	)

	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "test-event")

	meta := evt.Metadata()
	verify.NotEmpty(t, meta.ID)
	verify.Equal(t, meta.TraceID, "trace-123")
	verify.Equal(t, meta.Correlation, "corr-456")
	verify.Equal(t, meta.Custom["env"], "testing")
}

// TestMetadataInJSON verifies metadata is properly marshaled.
func TestMetadataInJSON(t *testing.T) {
	type Data struct {
		Value int
	}

	evt, err := NewEventTyped("test", Data{Value: 42},
		WithTraceID("trace-abc"),
		WithCorrelation("corr-xyz"),
	)
	verify.NoError(t, err)

	// Marshal to JSON
	jsonData, err := evt.MarshalJSON()
	verify.NoError(t, err)
	jsonStr := string(jsonData)
	verify.Substring(t, "metadata", jsonStr)
	verify.Substring(t, "trace-abc", jsonStr)
	verify.Substring(t, "corr-xyz", jsonStr)

	// Unmarshal from JSON
	var evt2 Event
	err = evt2.UnmarshalJSON(jsonData)
	verify.NoError(t, err)

	meta := evt2.Metadata()
	verify.Equal(t, meta.TraceID, "trace-abc")
	verify.Equal(t, meta.Correlation, "corr-xyz")
}

// TestBackwardCompatibility verifies old NewEvent still works.
func TestBackwardCompatibility(t *testing.T) {
	// Old style event creation should still work
	evt, err := NewEvent("old-style", 42, "test")
	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "old-style")
	verify.True(t, evt.HasPayload())

	// Should have auto-generated ID
	meta := evt.Metadata()
	verify.NotEmpty(t, meta.ID)
}

// TestMetadataGeneratesUniqueIDs verifies unique ID generation.
func TestMetadataGeneratesUniqueIDs(t *testing.T) {
	evt1, err := NewEventTyped("test", "data1")
	verify.NoError(t, err)

	evt2, err := NewEventTyped("test", "data2")
	verify.NoError(t, err)

	id1 := evt1.Metadata().ID
	id2 := evt2.Metadata().ID

	verify.NotEmpty(t, id1)
	verify.NotEmpty(t, id2)
	verify.Different(t, id1, id2)
}

// TestEventStringWithMetadata verifies String() includes metadata.
func TestEventStringWithMetadata(t *testing.T) {
	evt, err := NewEventTyped("test", "payload",
		WithTraceID("trace-123"),
		WithCorrelation("corr-456"),
	)
	verify.NoError(t, err)

	str := evt.String()
	verify.NotEmpty(t, str)
	// Event.String() includes TraceID and Correlation in the output
	verify.Substring(t, "trace-123", str)
	verify.Substring(t, "corr-456", str)
	verify.Substring(t, evt.Metadata().ID, str)
}

// TestCustomMetadata verifies custom metadata fields.
func TestCustomMetadata(t *testing.T) {
	evt, err := NewEventTyped("test", "data",
		WithCustom("requestID", "req-123"),
		WithCustom("userAgent", "Mozilla/5.0"),
		WithCustom("ipAddress", "192.168.1.1"),
	)
	verify.NoError(t, err)

	meta := evt.Metadata()
	verify.Equal(t, meta.Custom["requestID"], "req-123")
	verify.Equal(t, meta.Custom["userAgent"], "Mozilla/5.0")
	verify.Equal(t, meta.Custom["ipAddress"], "192.168.1.1")
}

// TestPayloadAsWithComplexTypes verifies complex type handling.
func TestPayloadAsWithComplexTypes(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		ZipCode string
	}

	type User struct {
		ID        string
		Name      string
		Email     string
		Addresses []Address
		CreatedAt time.Time
	}

	user := User{
		ID:    "user-123",
		Name:  "John Doe",
		Email: "john@example.com",
		Addresses: []Address{
			{Street: "123 Main St", City: "Springfield", ZipCode: "12345"},
			{Street: "456 Oak Ave", City: "Portland", ZipCode: "67890"},
		},
		CreatedAt: time.Now().UTC(),
	}

	evt, err := NewEventTyped("user-registered", user,
		WithTraceID("trace-user-reg"),
	)
	verify.NoError(t, err)

	extracted, err := PayloadAs[User](evt)
	verify.NoError(t, err)
	verify.Equal(t, extracted.ID, user.ID)
	verify.Equal(t, extracted.Name, user.Name)
	verify.Equal(t, extracted.Email, user.Email)
	verify.Length(t, extracted.Addresses, 2)
	verify.Equal(t, extracted.Addresses[0].City, "Springfield")
	verify.Equal(t, extracted.Addresses[1].City, "Portland")
}

// EOF
