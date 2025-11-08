# Sentry Configuration for DictaMesh

This directory contains configuration files for the self-hosted Sentry instance used by the DictaMesh framework.

## Directory Structure

```
sentry/
├── config/
│   ├── sentry.conf.py    # Main Sentry configuration (Python)
│   └── config.yml        # Additional configuration (YAML)
├── clickhouse/
│   └── 01-init-sentry.sql # ClickHouse initialization script
├── init-sentry.sh        # Sentry initialization script
└── README.md            # This file
```

## Configuration Files

### sentry.conf.py

The main Python configuration file for Sentry. This file configures:
- Database connections (PostgreSQL)
- Redis caching and rate limiting
- Kafka/Redpanda event streaming
- ClickHouse event storage
- File storage
- Security settings
- Performance monitoring

**Important**: Change `SENTRY_SECRET_KEY` in production!

### config.yml

Additional YAML-based configuration for:
- Email settings
- System configuration
- Authentication settings
- Integration configurations (GitHub, Slack, etc.)

### init-sentry.sh

Initialization script that:
1. Waits for PostgreSQL and Redis to be ready
2. Runs database migrations
3. Creates a default superuser
4. Sets up a default organization and team

**Usage**:
```bash
docker exec -it dictamesh-sentry-web /etc/sentry/init-sentry.sh
```

Or run it as part of the Docker Compose startup.

## Environment Variables

The following environment variables can be configured in `docker-compose.dev.yml`:

### Required
- `SENTRY_SECRET_KEY`: Secret key for cryptographic signing (change in production!)
- `SENTRY_POSTGRES_HOST`: PostgreSQL host
- `SENTRY_DB_NAME`: PostgreSQL database name
- `SENTRY_DB_USER`: PostgreSQL user
- `SENTRY_DB_PASSWORD`: PostgreSQL password
- `SENTRY_REDIS_HOST`: Redis host
- `SENTRY_KAFKA_HOSTS`: Kafka/Redpanda broker addresses

### Optional
- `SENTRY_ADMIN_EMAIL`: Admin user email (default: admin@dictamesh.local)
- `SENTRY_ADMIN_PASSWORD`: Admin user password (default: admin)
- `SENTRY_URL_PREFIX`: Sentry URL prefix (default: http://localhost:9000)
- `SENTRY_SINGLE_ORGANIZATION`: Single organization mode (default: 0)
- `SENTRY_EVENT_RETENTION_DAYS`: Event retention period (default: 90)
- `SENTRY_METRICS_SAMPLE_RATE`: Metrics sample rate (default: 1.0)
- `SENTRY_PROFILES_SAMPLE_RATE`: Profiles sample rate (default: 1.0)

## First-Time Setup

1. Start the infrastructure:
   ```bash
   make dev-up
   ```

2. Wait for all services to be healthy (this may take a few minutes)

3. Initialize Sentry:
   ```bash
   make sentry-init
   ```

4. Access Sentry at http://localhost:9000

5. Log in with default credentials:
   - Email: `admin@dictamesh.local`
   - Password: `admin`

6. **Important**: Change the admin password immediately!

## Production Considerations

### Security

1. **Change the secret key**:
   ```bash
   # Generate a new secret key
   python3 -c "import secrets; print(secrets.token_urlsafe(50))"
   ```

2. **Use strong passwords**: Never use default passwords in production

3. **Enable HTTPS**: Configure TLS/SSL certificates

4. **Configure email**: Set up SMTP for notifications

### Performance

1. **Scale workers**: Increase the number of Sentry worker containers

2. **Database tuning**: Optimize PostgreSQL for production workload

3. **Redis clustering**: Use Redis Cluster or Sentinel for HA

4. **ClickHouse optimization**: Configure ClickHouse for high-volume events

### Monitoring

1. **Health checks**: Monitor all Sentry services

2. **Resource usage**: Track CPU, memory, and disk usage

3. **Event ingestion**: Monitor event processing rates

4. **Error rates**: Track Sentry's own error rates

## Integration with DictaMesh

### Framework Components

DictaMesh framework components should use Sentry for:
- Error tracking and monitoring
- Performance monitoring (APM)
- Release tracking
- User feedback collection

### SDK Configuration

Add the Sentry SDK to your framework components:

**Go**:
```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: "http://your-dsn@localhost:9000/1",
    Environment: "development",
    Release: "dictamesh@1.0.0",
})
```

**Node.js**:
```javascript
const Sentry = require('@sentry/node');

Sentry.init({
  dsn: 'http://your-dsn@localhost:9000/1',
  environment: 'development',
  release: 'dictamesh@1.0.0',
});
```

**Python**:
```python
import sentry_sdk

sentry_sdk.init(
    dsn="http://your-dsn@localhost:9000/1",
    environment="development",
    release="dictamesh@1.0.0",
)
```

### Getting Your DSN

1. Log in to Sentry (http://localhost:9000)
2. Create a new project for your framework component
3. Copy the DSN from the project settings
4. Configure your application with the DSN

## Troubleshooting

### Sentry won't start

Check the logs:
```bash
docker logs dictamesh-sentry-web
docker logs dictamesh-sentry-worker
```

Verify dependencies are healthy:
```bash
make health
```

### Database migrations fail

Reset the Sentry database:
```bash
docker exec -it dictamesh-sentry-postgres psql -U sentry -d sentry -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
make sentry-init
```

### ClickHouse connection issues

Check ClickHouse logs:
```bash
docker logs dictamesh-clickhouse
```

Test ClickHouse connection:
```bash
docker exec -it dictamesh-clickhouse clickhouse-client --query "SELECT 1"
```

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
