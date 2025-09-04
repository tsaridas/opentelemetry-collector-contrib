// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package custommetricencodingextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/custommetricencodingextension"

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding"
)

var (
	_ encoding.MetricsMarshalerExtension = (*customMetricEncodingExtension)(nil)
)

type customMetricEncodingExtension struct {
	config *Config
	logger *zap.Logger
}

type TransformedMetric struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
	Ts    string      `json:"ts"`
}

func newExtension(config *Config, logger *zap.Logger) (*customMetricEncodingExtension, error) {
	return &customMetricEncodingExtension{
		config: config,
		logger: logger,
	}, nil
}

func (e *customMetricEncodingExtension) Start(ctx context.Context, host component.Host) error {
	e.logger.Info("Custom metric encoding extension started")
	return nil
}

func (e *customMetricEncodingExtension) Shutdown(ctx context.Context) error {
	e.logger.Info("Custom metric encoding extension stopped")
	return nil
}

// MarshalMetrics transforms the OpenTelemetry metrics to your desired format
func (e *customMetricEncodingExtension) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	transformedMetrics := e.transformMetrics(md)
	return json.Marshal(transformedMetrics)
}

func (e *customMetricEncodingExtension) transformMetrics(md pmetric.Metrics) []TransformedMetric {
	var result []TransformedMetric

	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)

		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			sm := rm.ScopeMetrics().At(j)

			for k := 0; k < sm.Metrics().Len(); k++ {
				metric := sm.Metrics().At(k)

				// Transform each metric based on its type
				transformed := e.transformMetric(metric, rm.Resource())
				result = append(result, transformed...)
			}
		}
	}

	return result
}

func (e *customMetricEncodingExtension) transformMetric(metric pmetric.Metric, resource pcommon.Resource) []TransformedMetric {
	var result []TransformedMetric

	// Extract host identifier from resource attributes
	hostIdentifier := e.extractHostIdentifier(resource)

	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		for i := 0; i < metric.Gauge().DataPoints().Len(); i++ {
			dp := metric.Gauge().DataPoints().At(i)
			result = append(result, e.createNumberDataPointMetric(metric.Name(), dp, hostIdentifier))
		}

	case pmetric.MetricTypeSum:
		for i := 0; i < metric.Sum().DataPoints().Len(); i++ {
			dp := metric.Sum().DataPoints().At(i)
			result = append(result, e.createNumberDataPointMetric(metric.Name(), dp, hostIdentifier))
		}

	case pmetric.MetricTypeHistogram:
		for i := 0; i < metric.Histogram().DataPoints().Len(); i++ {
			dp := metric.Histogram().DataPoints().At(i)
			result = append(result, e.createHistogramMetric(metric.Name(), dp, hostIdentifier))
		}

	case pmetric.MetricTypeExponentialHistogram:
		for i := 0; i < metric.ExponentialHistogram().DataPoints().Len(); i++ {
			dp := metric.ExponentialHistogram().DataPoints().At(i)
			result = append(result, e.createExponentialHistogramMetric(metric.Name(), dp, hostIdentifier))
		}

	case pmetric.MetricTypeSummary:
		for i := 0; i < metric.Summary().DataPoints().Len(); i++ {
			dp := metric.Summary().DataPoints().At(i)
			result = append(result, e.createSummaryMetric(metric.Name(), dp, hostIdentifier))
		}
	}

	return result
}

// extractHostIdentifier extracts the host identifier from resource attributes
func (e *customMetricEncodingExtension) extractHostIdentifier(resource pcommon.Resource) string {
	// Try to get host.name first
	if hostName, exists := resource.Attributes().Get("host.name"); exists {
		return hostName.Str()
	}

	// Fallback to other common host identifiers
	if hostName, exists := resource.Attributes().Get("host"); exists {
		return hostName.Str()
	}

	if instanceID, exists := resource.Attributes().Get("instance.id"); exists {
		return instanceID.Str()
	}

	if serviceInstanceID, exists := resource.Attributes().Get("service.instance.id"); exists {
		return serviceInstanceID.Str()
	}

	// If no host identifier found, use a default
	return "unknown-host"
}

func (e *customMetricEncodingExtension) createNumberDataPointMetric(metricName string, dp pmetric.NumberDataPoint, hostIdentifier string) TransformedMetric {
	// Get the value
	var value interface{}
	switch dp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		value = dp.IntValue()
	case pmetric.NumberDataPointValueTypeDouble:
		value = dp.DoubleValue()
	default:
		value = 0
	}

	// Get timestamps
	startTs := dp.StartTimestamp().AsTime().UnixNano()
	endTs := dp.Timestamp().AsTime().UnixNano()
	ts := fmt.Sprintf("%d %d", startTs, endTs)

	// Build the path
	path := fmt.Sprintf("/%s/%s", metricName, hostIdentifier)

	// Add label values to path if they exist
	if dp.Attributes().Len() > 0 {
		for _, attrVal := range dp.Attributes().AsRaw() {
			// Add all attributes to the path
			path += "/" + fmt.Sprintf("%v", attrVal)
		}
	}

	return TransformedMetric{
		Path:  path,
		Value: value,
		Ts:    ts,
	}
}

func (e *customMetricEncodingExtension) createHistogramMetric(metricName string, dp pmetric.HistogramDataPoint, hostIdentifier string) TransformedMetric {
	// For histograms, we'll use the count as the value
	value := dp.Count()

	// Get timestamps
	startTs := dp.StartTimestamp().AsTime().UnixNano()
	endTs := dp.Timestamp().AsTime().UnixNano()
	ts := fmt.Sprintf("%d %d", startTs, endTs)

	// Build the path
	path := fmt.Sprintf("/%s/%s", metricName, hostIdentifier)

	// Add label values to path if they exist
	if dp.Attributes().Len() > 0 {
		for _, attrVal := range dp.Attributes().AsRaw() {
			path += "/" + fmt.Sprintf("%v", attrVal)
		}
	}

	return TransformedMetric{
		Path:  path,
		Value: value,
		Ts:    ts,
	}
}

func (e *customMetricEncodingExtension) createExponentialHistogramMetric(metricName string, dp pmetric.ExponentialHistogramDataPoint, hostIdentifier string) TransformedMetric {
	// For exponential histograms, we'll use the count as the value
	value := dp.Count()

	// Get timestamps
	startTs := dp.StartTimestamp().AsTime().UnixNano()
	endTs := dp.Timestamp().AsTime().UnixNano()
	ts := fmt.Sprintf("%d %d", startTs, endTs)

	// Build the path
	path := fmt.Sprintf("/%s/%s", metricName, hostIdentifier)

	// Add label values to path if they exist
	if dp.Attributes().Len() > 0 {
		for _, attrVal := range dp.Attributes().AsRaw() {
			path += "/" + fmt.Sprintf("%v", attrVal)
		}
	}

	return TransformedMetric{
		Path:  path,
		Value: value,
		Ts:    ts,
	}
}

func (e *customMetricEncodingExtension) createSummaryMetric(metricName string, dp pmetric.SummaryDataPoint, hostIdentifier string) TransformedMetric {
	// For summaries, we'll use the count as the value
	value := dp.Count()

	// Get timestamps
	startTs := dp.StartTimestamp().AsTime().UnixNano()
	endTs := dp.Timestamp().AsTime().UnixNano()
	ts := fmt.Sprintf("%d %d", startTs, endTs)

	// Build the path
	path := fmt.Sprintf("/%s/%s", metricName, hostIdentifier)

	// Add label values to path if they exist
	if dp.Attributes().Len() > 0 {
		for _, attrVal := range dp.Attributes().AsRaw() {
			path += "/" + fmt.Sprintf("%v", attrVal)
		}
	}

	return TransformedMetric{
		Path:  path,
		Value: value,
		Ts:    ts,
	}
}

