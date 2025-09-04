// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package custommetricencodingextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/custommetricencodingextension"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	// TypeStr is the type of the extension.
	TypeStr = "custommetricencoding"
)

// NewFactory creates a factory for the custom metric encoding extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		component.MustNewType(TypeStr),
		createDefaultConfig,
		createExtension,
		component.StabilityLevelAlpha,
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createExtension(
	_ context.Context,
	set extension.Settings,
	config component.Config,
) (extension.Extension, error) {
	cfg := config.(*Config)

	return newExtension(cfg, set.Logger)
}
