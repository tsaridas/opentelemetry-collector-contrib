# Custom Metric Encoding Extension

This extension transforms OpenTelemetry metrics into a custom JSON format with `path`, `value`, and `ts` fields.

## Description

The Custom Metric Encoding Extension transforms OpenTelemetry metrics into a simplified JSON format that's easier to consume by downstream systems. It creates a path-based structure that includes the metric name and attribute values, without hostname information.

## Configuration

```yaml
extensions:
  custommetricencoding: {}
```

**No configuration needed** - the extension creates clean metric paths without hostname dependencies.

## Output Format

The extension transforms OpenTelemetry metrics into a simplified JSON array:

```json
[
  {
    "path": "/go.memory.used/stack",
    "value": 589824,
    "ts": "1756975296124471296"
  },
  {
    "path": "/go.memory.used/other",
    "value": 10378256,
    "ts": "1756975296124471296"
  }
]
```

### Field Descriptions

- **`path`**: Hierarchical path containing metric name and attribute values
- **`value`**: The actual metric value (gauge value, sum value, histogram count, etc.)
- **`ts`**: Current timestamp in Unix nanoseconds

## Path Structure

The path follows this pattern:
```
/metric_name/attribute_value1/attribute_value2/...
```

**Examples:**
- `/go.memory.used/stack`
- `/go.goroutine.count`
- `/cpu.usage/cpu0/user`

## Supported Metric Types

The extension handles all OpenTelemetry metric types:

- **Gauge**: Uses the gauge value directly
- **Sum**: Uses the sum value directly  
- **Histogram**: Uses the count value
- **Exponential Histogram**: Uses the count value
- **Summary**: Uses the count value

## Integration with MQTT Exporter

This extension works seamlessly with the MQTT exporter:

```yaml
exporters:
  mqtt:
    # ... MQTT configuration ...
    encoding_extension: custommetricencoding
```

**Key Benefits:**
1. **Clean metric paths** - No hostname clutter in metric paths
2. **Flexible attribute handling** - All attributes are included in the path
3. **Type-agnostic** - Works with all OpenTelemetry metric types
4. **Zero configuration** - No manual setup needed
5. **Simplified timestamps** - Single timestamp per metric

## Building and Testing

### Build the Extension

```bash
cd extension/encoding/custommetricencodingextension
./build.sh
```

### Build the Collector

```bash
make otelcontribcol
```

### Test Configuration

```bash
./bin/otelcontribcol_linux_arm64 validate --config extension/encoding/custommetricencodingextension/test-config.yaml
```

## Use Cases

1. **Simplified metric consumption** - Clean paths without hostname dependencies
2. **Containerized systems** - Focus on metric attributes rather than host identification
3. **Cloud-native applications** - Streamlined metric paths for better aggregation
4. **IoT deployments** - Device-agnostic metric collection
5. **Microservices** - Service-agnostic metric paths for better correlation
