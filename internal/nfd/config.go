package nfd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
)

// NFDConfig contains only the fallback catalog source index image for NFD.
type NFDConfig struct {
	FallbackCatalogSourceIndexImage string `envconfig:"NFD_FALLBACK_CATALOGSOURCE_INDEX_IMAGE"`
}

// NewNFDConfig attempts to load NFDConfig from the environment.
// Logs at V(100) and returns (*NFDConfig, nil) on success, or (nil, error) on failure.
func NewNFDConfig() (*NFDConfig, error) {
	glog.V(100).Info("Creating new NFDConfig")

	cfg := &NFDConfig{}
	if err := envconfig.Process("NFD_", cfg); err != nil {
		return nil, fmt.Errorf("failed to process NFD_ env vars: %w", err)
	}

	glog.V(100).Info("NFDConfig created successfully")
	return cfg, nil
}
