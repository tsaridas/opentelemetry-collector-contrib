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

	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		for i := 0; i < metric.Gauge().DataPoints().Len(); i++ {
			dp := metric.Gauge().DataPoints().At(i)
			result = append(result, e.createNumberDataPointMetric(metric.Name(), dp))
		}

	case pmetric.MetricTypeSum:
		for i := 0; i < metric.Sum().DataPoints().Len(); i++ {
			dp := metric.Sum().DataPoints().At(i)
			result = append(result, e.createNumberDataPointMetric(metric.Name(), dp))
		}

	case pmetric.MetricTypeHistogram:
		for i := 0; i < metric.Histogram().DataPoints().Len(); i++ {
			dp := metric.Histogram().DataPoints().At(i)
			result = append(result, e.createHistogramMetric(metric.Name(), dp))
		}

	case pmetric.MetricTypeExponentialHistogram:
		for i := 0; i < metric.ExponentialHistogram().DataPoints().Len(); i++ {
			dp := metric.ExponentialHistogram().DataPoints().At(i)
			result = append(result, e.createExponentialHistogramMetric(metric.Name(), dp))
		}

	case pmetric.MetricTypeSummary:
		for i := 0; i < metric.Summary().DataPoints().Len(); i++ {
			dp := metric.Summary().DataPoints().At(i)
			result = append(result, e.createSummaryMetric(metric.Name(), dp))
		}
	}

	return result
}

func (e *customMetricEncodingExtension) createNumberDataPointMetric(metricName string, dp pmetric.NumberDataPoint) TransformedMetric {
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

	// Get timestamp (only current timestamp, not start timestamp)
	ts := fmt.Sprintf("%d", dp.Timestamp().AsTime().UnixNano())

	// Build the path
	path := fmt.Sprintf("/%s", metricName)

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

func (e *customMetricEncodingExtension) createHistogramMetric(metricName string, dp pmetric.HistogramDataPoint) TransformedMetric {
	// For histograms, we'll use the count as the value
	value := dp.Count()

	// Get timestamp (only current timestamp, not start timestamp)
	ts := fmt.Sprintf("%d", dp.Timestamp().AsTime().UnixNano())

	// Build the path
	path := fmt.Sprintf("/%s", metricName)

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

func (e *customMetricEncodingExtension) createExponentialHistogramMetric(metricName string, dp pmetric.ExponentialHistogramDataPoint) TransformedMetric {
	// For exponential histograms, we'll use the count as the value
	value := dp.Count()

	// Get timestamp (only current timestamp, not start timestamp)
	ts := fmt.Sprintf("%d", dp.Timestamp().AsTime().UnixNano())

	// Build the path
	path := fmt.Sprintf("/%s", metricName)

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

func (e *customMetricEncodingExtension) createSummaryMetric(metricName string, dp pmetric.SummaryDataPoint) TransformedMetric {
	// For summaries, we'll use the count as the value
	value := dp.Count()

	// Get timestamp (only current timestamp, not start timestamp)
	ts := fmt.Sprintf("%d", dp.Timestamp().AsTime().UnixNano())

	// Build the path
	path := fmt.Sprintf("/%s", metricName)

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
