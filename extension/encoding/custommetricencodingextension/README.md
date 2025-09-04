# Custom Metric Encoding Extension

This extension transforms OpenTelemetry metrics into a custom JSON format with `path`, `value`, and `ts` fields.

## Description

The Custom Metric Encoding Extension transforms OpenTelemetry metrics into a simplified JSON format that's easier to consume by downstream systems. It creates a path-based structure that includes the metric name, host identifier (derived from resource attributes), and attribute values.

## Configuration

```yaml
extensions:
  custommetricencoding: {}
```

**No configuration needed** - the extension automatically derives the host identifier from the resource attributes of each metric.

### Host Identifier Derivation

The extension automatically extracts the host identifier from resource attributes in this priority order:

1. `host.name` - Primary host identifier
2. `host` - Alternative host identifier  
3. `instance.id` - Instance identifier
4. `service.instance.id` - Service instance identifier
5. `unknown-host` - Default fallback if none found

## Output Format

The extension transforms OpenTelemetry metrics into a simplified JSON array:

```json
[
  {
    "path": "/go.memory.used/idubsskcsrv1101.test.ansp.skyguide.ch/stack",
    "value": 589824,
    "ts": "1753861810458287915 1753861860459767723"
  }
]
```

### Field Descriptions

- **`path`**: Hierarchical path containing metric name, host identifier, and attribute values
- **`value`**: The actual metric value (gauge value, sum value, histogram count, etc.)
- **`ts`**: Timestamps in format "start_timestamp end_timestamp" (Unix nanoseconds)

## Path Structure

The path follows this pattern:
```
/metric_name/host_identifier/attribute_value1/attribute_value2/...
```

**Examples:**
- `/go.memory.used/idubsskcsrv1101.test.ansp.skyguide.ch/stack`
- `/go.goroutine.count/idubsskcsrv1101.test.ansp.skyguide.ch`
- `/cpu.usage/idubsskcsrv1101.test.ansp.skyguide.ch/cpu0/user`

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
1. **Dynamic host identification** - Automatically extracts host from resource attributes
2. **Flexible attribute handling** - All attributes are included in the path
3. **Type-agnostic** - Works with all OpenTelemetry metric types
4. **Zero configuration** - No manual host identifier setup needed
5. **Resource-aware** - Each metric batch can have different host identifiers

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

1. **Multi-host environments** - Automatically identify metrics by host
2. **Containerized systems** - Extract instance identifiers from resource attributes
3. **Cloud-native applications** - Use cloud instance IDs as host identifiers
4. **IoT deployments** - Identify devices by their resource attributes
5. **Microservices** - Distinguish between different service instances
