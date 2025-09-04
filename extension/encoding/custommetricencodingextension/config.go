// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package custommetricencodingextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/custommetricencodingextension"

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the custom metric encoding extension.
type Config struct {
	// No configuration needed - host identifier is derived from resource attributes
}

// Validate checks if the extension configuration is valid.
func (cfg *Config) Validate() error {
	// No validation needed
	return nil
}

// Ensure Config implements the component.Config interface.
var _ component.Config = (*Config)(nil)
