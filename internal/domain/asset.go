package domain

import "strings"

// AssetID uniquely identifies an asset record.
type AssetID string

// Asset represents a material or item that can be assigned under custody.
type Asset struct {
	id     AssetID
	name   string
	active bool
}

// NewAsset creates an active Asset with validated required fields.
func NewAsset(id AssetID, name string) (Asset, error) {
	return buildAsset(id, name, true)
}

// ReconstituteAsset rebuilds an Asset from persisted state.
func ReconstituteAsset(id AssetID, name string, active bool) (Asset, error) {
	return buildAsset(id, name, active)
}

func buildAsset(id AssetID, name string, active bool) (Asset, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Asset{}, ErrEmptyAssetID
	}

	name = NormalizeAssetName(name)
	if name == "" {
		return Asset{}, ErrEmptyAssetName
	}

	return Asset{
		id:     id,
		name:   name,
		active: active,
	}, nil
}

// NormalizeAssetName normalizes asset names for consistent identity.
func NormalizeAssetName(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// ID returns the asset identifier.
func (a Asset) ID() AssetID {
	return a.id
}

// Name returns the asset name.
func (a Asset) Name() string {
	return a.name
}

// Active reports whether the asset record is active.
func (a Asset) Active() bool {
	return a.active
}
