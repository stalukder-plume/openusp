# Operations & Observability (Stub)

Purpose: Practical runbook material for operating OpenUSP clusters.

## 1. Health Endpoints
| Path | Service | Purpose |
|------|---------|---------|
| /api/v1/health | API | Composite health |
| /health | Controller | Basic liveness (confirm) |
| /metrics | All (if enabled) | Prometheus metrics |

## 2. Key Metrics (Draft)
| Metric Group | Examples | Actionable Use |
|--------------|----------|----------------|
| Device Fleet | device_total, device_connected_ratio | Capacity planning |
| Protocol | usp_messages_per_sec, cwmp_sessions_active | Transport scaling |
| API | http_requests_total, http_request_duration_seconds | SLO tracking |
| Database | mongo_ops_total, mongo_conn_pool_in_use | Index/scale decisions |
| Cache | redis_hits_ratio | Tuning TTL / sizing |

## 3. Logging
- Structured JSON recommended
- Correlate by request ID / device ID
- Consider log sampling at scale

## 4. Tracing (Future)
- OpenTelemetry integration planned
- Trace spans: ingress → controller → broker publish → device response

## 5. Alerting Suggestions
| Condition | Threshold (example) | Response |
|-----------|---------------------|----------|
| Device connect fail rate | >5% 5m window | Investigate broker / auth |
| API p95 latency | >500ms 10m | Check DB indices / saturation |
| Broker unacked frames | > N threshold | Inspect consumer lag |
| MongoDB replication lag | > 30s | Check secondary health |

## 6. Backup & Recovery (Future Outline)
- MongoDB replica set snapshot strategy
- Parameter state export (tooling TBD)

## 7. Capacity Planning Inputs
| Dimension | Driver | Notes |
|-----------|--------|-------|
| Controller instances | Active device count, message rate | Horizontal scale |
| Broker cluster size | Concurrent connections, throughput | STOMP/MQTT differences |
| MongoDB IOPS | Parameter churn, event volume | Consider sharding |
| Redis memory | Session / ephemeral state size | Monitor fragmentation |

## 8. Troubleshooting Pointers
| Symptom | First Checks |
|---------|-------------|
| Devices not appearing | Broker connectivity, controller logs |
| High API latency | DB metrics, goroutine dumps |
| Frequent reconnects | Heartbeat mismatch, network stability |

## TODO
- Replace speculative paths with confirmed endpoints.
- Add sample Prometheus rules and Grafana dashboard JSON links.
- Include structured logging field dictionary.
