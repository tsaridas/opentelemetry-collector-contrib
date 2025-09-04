// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package custommetricencodingextension

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestNewExtension(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()

	ext, err := newExtension(config, logger)
	require.NoError(t, err)
	assert.NotNil(t, ext)
}

func TestExtensionLifecycle(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	host := componenttest.NewNopHost()

	// Test start
	err = ext.Start(context.Background(), host)
	assert.NoError(t, err)

	// Test shutdown
	err = ext.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestMarshalMetrics(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create test metrics
	md := createTestMetrics()

	// Marshal metrics
	result, err := ext.MarshalMetrics(md)
	require.NoError(t, err)

	// Parse the result
	var transformedMetrics []TransformedMetric
	err = json.Unmarshal(result, &transformedMetrics)
	require.NoError(t, err)

	// Verify the structure
	assert.Len(t, transformedMetrics, 2) // Two data points

	// Check first metric (go.memory.used with stack attribute)
	firstMetric := transformedMetrics[0]
	assert.Equal(t, "/go.memory.used/stack", firstMetric.Path)
	assert.Equal(t, float64(589824), firstMetric.Value)
	assert.NotEmpty(t, firstMetric.Ts)

	// Check second metric (go.memory.used with other attribute)
	secondMetric := transformedMetrics[1]
	assert.Equal(t, "/go.memory.used/other", secondMetric.Path)
	assert.Equal(t, float64(10378256), secondMetric.Value)
	assert.NotEmpty(t, secondMetric.Ts)
}

func TestTransformMetrics(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create test metrics
	md := createTestMetrics()

	// Transform metrics
	result := ext.transformMetrics(md)

	// Verify the result
	assert.Len(t, result, 2)

	// Check that paths don't contain hostname
	for _, metric := range result {
		assert.NotContains(t, metric.Path, "hostname2")
		assert.NotContains(t, metric.Path, "unknown-host")
		assert.True(t, len(metric.Ts) > 0)
	}
}

func TestCreateNumberDataPointMetric(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test data point
	dp := pmetric.NewNumberDataPoint()
	dp.SetIntValue(42)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// Add attributes
	dp.Attributes().PutStr("type", "test")

	// Create metric
	result := ext.createNumberDataPointMetric("test.metric", dp)

	// Verify the result
	assert.Equal(t, "/test.metric/test", result.Path)
	assert.Equal(t, int64(42), result.Value)
	assert.NotEmpty(t, result.Ts)
}

func TestCreateNumberDataPointMetricWithDoubleValue(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test data point with double value
	dp := pmetric.NewNumberDataPoint()
	dp.SetDoubleValue(42.5)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// Create metric
	result := ext.createNumberDataPointMetric("test.metric", dp)

	// Verify the result
	assert.Equal(t, "/test.metric", result.Path)
	assert.Equal(t, 42.5, result.Value)
	assert.NotEmpty(t, result.Ts)
}

func TestCreateHistogramMetric(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test histogram data point
	dp := pmetric.NewHistogramDataPoint()
	dp.SetCount(100)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// Add attributes
	dp.Attributes().PutStr("bucket", "0-1")

	// Create metric
	result := ext.createHistogramMetric("test.histogram", dp)

	// Verify the result
	assert.Equal(t, "/test.histogram/0-1", result.Path)
	assert.Equal(t, uint64(100), result.Value)
	assert.NotEmpty(t, result.Ts)
}

func TestCreateExponentialHistogramMetric(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test exponential histogram data point
	dp := pmetric.NewExponentialHistogramDataPoint()
	dp.SetCount(200)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// Add attributes
	dp.Attributes().PutStr("scale", "2")

	// Create metric
	result := ext.createExponentialHistogramMetric("test.exp_histogram", dp)

	// Verify the result
	assert.Equal(t, "/test.exp_histogram/2", result.Path)
	assert.Equal(t, uint64(200), result.Value)
	assert.NotEmpty(t, result.Ts)
}

func TestCreateSummaryMetric(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test summary data point
	dp := pmetric.NewSummaryDataPoint()
	dp.SetCount(50)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// Add attributes
	dp.Attributes().PutStr("quantile", "0.5")

	// Create metric
	result := ext.createSummaryMetric("test.summary", dp)

	// Verify the result
	assert.Equal(t, "/test.summary/0.5", result.Path)
	assert.Equal(t, uint64(50), result.Value)
	assert.NotEmpty(t, result.Ts)
}

func TestTimestampFormat(t *testing.T) {
	config := &Config{}
	logger := zap.NewNop()
	ext, err := newExtension(config, logger)
	require.NoError(t, err)

	// Create a test data point with known timestamp
	testTime := time.Unix(0, 1756975296124471296)
	dp := pmetric.NewNumberDataPoint()
	dp.SetIntValue(42)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(testTime))

	// Create metric
	result := ext.createNumberDataPointMetric("test.metric", dp)

	// Verify the timestamp format
	assert.Equal(t, "1756975296124471296", result.Ts)
}

// Helper function to create test metrics similar to the example
func createTestMetrics() pmetric.Metrics {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()

	// Add resource attributes (but they won't be used in the path anymore)
	rm.Resource().Attributes().PutStr("host.name", "hostname2")
	rm.Resource().Attributes().PutStr("os.name", "linux")
	rm.Resource().Attributes().PutStr("service.name", "tracksreceiver")
	rm.Resource().Attributes().PutStr("service.version", "devtest")

	sm := rm.ScopeMetrics().AppendEmpty()
	sm.Scope().SetName("go.opentelemetry.io/contrib/instrumentation/runtime")
	sm.Scope().SetVersion("0.62.0")

	metric := sm.Metrics().AppendEmpty()
	metric.SetName("go.memory.used")
	metric.SetDescription("Memory used by the Go runtime.")
	metric.SetUnit("By")

	// Create sum metric
	sum := metric.SetEmptySum()
	sum.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)

	// First data point - stack
	dp1 := sum.DataPoints().AppendEmpty()
	dp1.Attributes().PutStr("go.memory.type", "stack")
	dp1.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(-5 * time.Second)))
	dp1.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	dp1.SetIntValue(589824)

	// Second data point - other
	dp2 := sum.DataPoints().AppendEmpty()
	dp2.Attributes().PutStr("go.memory.type", "other")
	dp2.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now().Add(-5 * time.Second)))
	dp2.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	dp2.SetIntValue(10378256)

	return md
}
