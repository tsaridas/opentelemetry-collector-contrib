// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package custommetricencodingextension provides a custom metric encoding extension
// for the OpenTelemetry Collector that transforms metrics into a simplified JSON format.
//
// This extension removes hostname information and uses only the current timestamp,
// producing output in the format:
//
//	[
//	  {
//	    "path": "/metric.name/attribute.value",
//	    "value": 123,
//	    "ts": "1756975296124471296"
//	  }
//	]
//
// The path is constructed from the metric name and any attributes present on the data points.
package custommetricencodingextension
