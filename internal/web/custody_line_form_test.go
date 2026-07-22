package web

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestParseCustodyLineCommandsFromRequestParsesMultipleLines(t *testing.T) {
	form := make(url.Values)
	form.Add("asset_id", "asset-1")
	form.Add("quantity", "1")
	form.Add("asset_id", "asset-2")
	form.Add("quantity", "2")

	request := &http.Request{
		PostForm: form,
	}

	lines, err := parseCustodyLineCommandsFromRequest(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	if lines[0].AssetID != "asset-1" {
		t.Fatalf("expected asset-1, got %s", lines[0].AssetID)
	}

	if lines[1].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", lines[1].Quantity)
	}
}

func TestParseCustodyLineCommandsFromRequestRejectsNoLines(t *testing.T) {
	request := &http.Request{
		PostForm: make(url.Values),
	}

	_, err := parseCustodyLineCommandsFromRequest(request)
	if !errors.Is(err, errNoCustodyLineSubmitted) {
		t.Fatalf("expected no line error, got %v", err)
	}
}

func TestDefaultCustodyLineFormRowsStartsQuantityAtOne(t *testing.T) {
	rows := defaultCustodyLineFormRows()

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	if rows[0].Quantity != "1" {
		t.Fatalf("expected quantity 1, got %q", rows[0].Quantity)
	}
}

func TestCustodyLineFormRowsFromRequestDefaultsEmptyQuantityToOne(t *testing.T) {
	form := make(url.Values)
	form.Add("asset_id", "asset-1")
	form.Add("quantity", "")

	request := &http.Request{
		PostForm: form,
	}

	rows := custodyLineFormRowsFromRequest(request)

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	if rows[0].Quantity != "1" {
		t.Fatalf("expected quantity 1, got %q", rows[0].Quantity)
	}
}
