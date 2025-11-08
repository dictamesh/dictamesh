# Loki Query Runbook

**Document Version:** 1.0
**Last Updated:** 2025-11-08
**Audience:** Developers, DevOps, SREs

## Overview

This runbook provides practical LogQL queries for common debugging and troubleshooting scenarios in DictaMesh. All queries assume you're using the Grafana Explore interface or the Loki API directly.

## Quick Reference

### Access Points

- **Grafana Explore:** http://localhost:3000/explore (select Loki datasource)
- **Loki API:** http://localhost:3100
- **Direct Query:** `curl -G -s "http://localhost:3100/loki/api/v1/query" --data-urlencode 'query={job="dictamesh-logs"}'`

### Common Labels

- `job`: Always "dictamesh-logs" for DictaMesh services
- `service`: Service name (e.g., "metadata-catalog", "postgres", "redis")
- `container_name`: Full container name (e.g., "dictamesh-metadata-catalog")
- `level`: Log level (debug, info, warn, error, fatal)
- `logger`: Logger component name

## Common Query Patterns

### 1. Basic Log Viewing

#### View all logs from a specific service
```logql
{service="metadata-catalog"}
```

#### View logs from multiple services
```logql
{service=~"metadata-catalog|graphql-gateway"}
```

#### View all DictaMesh logs
```logql
{job="dictamesh-logs"}
```

### 2. Filtering by Log Level

#### Show only errors
```logql
{service="metadata-catalog"} | json | level="error"
```

#### Show errors and warnings
```logql
{service="metadata-catalog"} | json | level=~"error|warn"
```

#### Show everything except debug logs
```logql
{service="metadata-catalog"} | json | level!="debug"
```

### 3. Text Search

#### Search for specific text in logs
```logql
{service="metadata-catalog"} |= "database connection"
```

#### Case-insensitive search
```logql
{service="metadata-catalog"} |~ "(?i)database"
```

#### Exclude lines containing specific text
```logql
{service="metadata-catalog"} != "health check"
```

#### Multiple text filters
```logql
{service="metadata-catalog"} |= "error" |= "timeout" != "expected"
```

### 4. JSON Field Extraction

#### Filter by HTTP status code
```logql
{service="metadata-catalog"} | json | status >= 500
```

#### Filter by duration
```logql
{service="metadata-catalog"} | json | duration_ms > 1000
```

#### Filter by user
```logql
{service="metadata-catalog"} | json | user_id="user-123"
```

#### Filter by HTTP method and path
```logql
{service="metadata-catalog"} | json | method="POST" | path=~"/api/.*"
```

## Troubleshooting Scenarios

### Scenario 1: Investigating a 500 Error Spike

**Query:** Find all 5xx errors in the last hour
```logql
{job="dictamesh-logs"} | json | status >= 500 and status < 600
```

**Query:** Group errors by service
```logql
sum by (service) (
  count_over_time({job="dictamesh-logs"} | json | status >= 500 [1h])
)
```

**Query:** Show error details with context
```logql
{job="dictamesh-logs"} | json | status >= 500
  | line_format "{{.ts}} [{{.service}}] {{.method}} {{.path}} -> {{.status}} ({{.error}})"
```

### Scenario 2: Finding Slow Requests

**Query:** Requests taking more than 1 second
```logql
{job="dictamesh-logs"} | json | duration_ms > 1000
```

**Query:** Top 10 slowest requests
```logql
topk(10,
  max_over_time({job="dictamesh-logs"} | json | unwrap duration_ms [1h])
)
```

**Query:** Average duration by endpoint
```logql
avg by (path) (
  rate({job="dictamesh-logs"} | json | unwrap duration_ms [5m])
)
```

### Scenario 3: Tracking a Distributed Transaction

**Query:** Find all logs for a specific trace ID
```logql
{job="dictamesh-logs"} | json | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```

**Query:** Trace flow across services (with formatted output)
```logql
{job="dictamesh-logs"} | json | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
  | line_format "{{.service}} | {{.msg}}"
```

**Query:** Find errors in a specific trace
```logql
{job="dictamesh-logs"} | json
  | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
  | level="error"
```

### Scenario 4: Database Issues

**Query:** Find all database-related errors
```logql
{job="dictamesh-logs"} |~ "(?i)(database|postgres|sql)" | json | level="error"
```

**Query:** Connection pool issues
```logql
{job="dictamesh-logs"} |~ "(?i)(connection pool|max connections|acquire timeout)"
```

**Query:** Slow queries
```logql
{service="postgres"} |= "duration:" | regexp "duration: (?P<duration>[0-9.]+) ms" | duration > 1000
```

### Scenario 5: Authentication Failures

**Query:** Failed login attempts
```logql
{job="dictamesh-logs"} |~ "(?i)(authentication failed|invalid credentials|unauthorized)"
```

**Query:** Count failed logins by user
```logql
sum by (user_id) (
  count_over_time({job="dictamesh-logs"} |= "authentication failed" | json [1h])
)
```

### Scenario 6: Service Health Issues

**Query:** Find service startup/shutdown events
```logql
{job="dictamesh-logs"} |~ "(?i)(starting|stopping|shutdown|terminated)"
```

**Query:** Out of memory errors
```logql
{job="dictamesh-logs"} |~ "(?i)(out of memory|oom|killed)"
```

**Query:** Panic/crash detection
```logql
{job="dictamesh-logs"} | json | level="fatal"
```

### Scenario 7: Rate Limiting Issues

**Query:** Find rate limit hits
```logql
{job="dictamesh-logs"} |~ "(?i)(rate limit|too many requests|429)"
```

**Query:** Count rate limits by service
```logql
sum by (service) (
  rate({job="dictamesh-logs"} |= "rate limit" [5m])
)
```

## Metrics from Logs

### Error Rate

```logql
sum by (service) (
  rate({job="dictamesh-logs"} | json | level="error" [5m])
)
```

### Request Rate

```logql
sum by (service, path) (
  rate({job="dictamesh-logs"} | json | status > 0 [5m])
)
```

### P95 Latency

```logql
histogram_quantile(0.95,
  sum by (le, service) (
    rate({job="dictamesh-logs"} | json | unwrap duration_ms | __error__="" [5m])
  )
)
```

### Success Rate

```logql
sum by (service) (
  rate({job="dictamesh-logs"} | json | status >= 200 and status < 300 [5m])
)
/
sum by (service) (
  rate({job="dictamesh-logs"} | json | status > 0 [5m])
)
```

## Log Formatting

### Basic Formatting

```logql
{service="metadata-catalog"} | json
  | line_format "{{.level | upper}} | {{.msg}}"
```

### Include Timestamp and Service

```logql
{job="dictamesh-logs"} | json
  | line_format "{{.ts}} [{{.service}}] {{.level | upper}}: {{.msg}}"
```

### HTTP Request Format

```logql
{job="dictamesh-logs"} | json | method != ""
  | line_format "{{.method}} {{.path}} → {{.status}} ({{.duration_ms}}ms)"
```

### Error Details

```logql
{job="dictamesh-logs"} | json | level="error"
  | line_format "ERROR in {{.service}}: {{.msg}}\n  File: {{.caller}}\n  Trace: {{.trace_id}}"
```

## Pattern Extraction

### Extract HTTP Status Codes

```logql
{job="dictamesh-logs"} | pattern `<_> status=<status>` | status >= 400
```

### Extract User IDs

```logql
{job="dictamesh-logs"} | pattern `<_> user_id=<user>` | user != ""
```

### Extract Error Types

```logql
{job="dictamesh-logs"} | pattern `<_> error=<error_type>:`
```

## Performance Optimization Tips

### 1. Use Labels Wisely
- Always start with label filters: `{service="foo"}` not `{} | json | service="foo"`
- Labels are indexed, JSON fields are not

### 2. Limit Time Range
- Shorter time ranges = faster queries
- Use `$__range` variable in Grafana dashboards

### 3. Use Sampling for High Volume
- Add line filters early: `{service="foo"} |= "error"` before `| json`
- Use `| json` only when needed

### 4. Avoid High Cardinality
- Don't add high-cardinality fields (user_id, request_id) as labels
- Keep these as JSON fields and filter after parsing

## Common Issues and Solutions

### Issue: "too many outstanding requests"
**Solution:** Reduce query time range or add more filters

### Issue: "maximum of series (20000) reached for a single query"
**Solution:** Add more specific label filters to reduce cardinality

### Issue: "Query timeout"
**Solution:**
- Reduce time range
- Add text filters before JSON parsing
- Break query into smaller time chunks

### Issue: "no logs found"
**Solutions:**
1. Check Promtail is running: `docker ps | grep promtail`
2. Check Loki is receiving logs: `curl http://localhost:3100/ready`
3. Verify label values: `curl http://localhost:3100/loki/api/v1/label/service/values`

## Useful Loki API Endpoints

### Health Check
```bash
curl http://localhost:3100/ready
curl http://localhost:3100/metrics
```

### Query Labels
```bash
# List all label names
curl http://localhost:3100/loki/api/v1/labels

# List values for a specific label
curl http://localhost:3100/loki/api/v1/label/service/values
```

### Query Logs
```bash
# Query range (time series)
curl -G -s "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={service="metadata-catalog"}' \
  --data-urlencode 'start=1h' \
  | jq .

# Instant query
curl -G -s "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query=sum(count_over_time({job="dictamesh-logs"}[1h]))' \
  | jq .
```

### Tail Logs (Live)
```bash
curl -G -s "http://localhost:3100/loki/api/v1/tail" \
  --data-urlencode 'query={service="metadata-catalog"}' \
  -H "Connection: keep-alive"
```

## LogQL Cheat Sheet

| Operation | Syntax | Example |
|-----------|--------|---------|
| Exact match | `{label="value"}` | `{service="postgres"}` |
| Regex match | `{label=~"regex"}` | `{service=~".*-adapter"}` |
| Not equal | `{label!="value"}` | `{level!="debug"}` |
| Regex not match | `{label!~"regex"}` | `{service!~"sentry.*"}` |
| Contains | `\|= "text"` | `\|= "error"` |
| Not contains | `!= "text"` | `!= "health check"` |
| Regex filter | `\|~ "regex"` | `\|~ "(?i)database"` |
| Parse JSON | `\| json` | `\| json \| level="error"` |
| Line format | `\| line_format "..."` | `\| line_format "{{.msg}}"` |
| Label format | `\| label_format new=old` | `\| label_format svc=service` |
| Unwrap | `\| unwrap field` | `\| unwrap duration_ms` |
| Rate | `rate({...}[5m])` | `rate({job="logs"}[5m])` |
| Count over time | `count_over_time({...}[1h])` | `count_over_time({level="error"}[1h])` |
| Sum | `sum(...)` | `sum(rate({...}[5m]))` |
| Sum by | `sum by (label)(...)` | `sum by (service)(...)` |

## Additional Resources

- [LogQL Documentation](https://grafana.com/docs/loki/latest/logql/)
- [Loki API Reference](https://grafana.com/docs/loki/latest/api/)
- [Grafana Loki Best Practices](https://grafana.com/docs/loki/latest/best-practices/)
- [DictaMesh Loki Integration Design](../architecture/LOKI-LOGGING-INTEGRATION.md)

## Support

For issues with Loki or logging:
1. Check Loki logs: `docker logs dictamesh-loki`
2. Check Promtail logs: `docker logs dictamesh-promtail`
3. Verify Grafana datasource configuration
4. Review this runbook for common issues

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Maintainer: DevOps Team
