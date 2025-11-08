# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2025 Controle Digital Ltda

.PHONY: help dev-up dev-down dev-logs dev-ps clean-volumes dev-reset kafka-topics kafka-ui redis-cli postgres-cli

# Default target
help: ## Show this help message
	@echo 'DictaMesh Framework Infrastructure Management'
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development environment commands
dev-up: ## Start development infrastructure
	@echo "🚀 Starting DictaMesh development infrastructure..."
	docker-compose -f docker-compose/docker-compose.dev.yml up -d
	@echo ""
	@echo "✅ Infrastructure is starting up. Services:"
	@echo "   📊 Redpanda Console: http://localhost:8080"
	@echo "   📈 Grafana:          http://localhost:3000 (admin/admin)"
	@echo "   📉 Prometheus:       http://localhost:9090"
	@echo "   🔍 Jaeger UI:        http://localhost:16686"
	@echo "   🐘 PostgreSQL:       localhost:5432 (dictamesh/dictamesh_dev_password)"
	@echo "   🔴 Redis:            localhost:6379"
	@echo "   🎯 Redpanda Kafka:   localhost:19092"
	@echo ""
	@echo "Run 'make dev-logs' to see container logs"

dev-down: ## Stop development infrastructure
	@echo "🛑 Stopping DictaMesh development infrastructure..."
	docker-compose -f docker-compose/docker-compose.dev.yml down
	@echo "✅ Infrastructure stopped"

dev-logs: ## Follow logs from all containers
	docker-compose -f docker-compose/docker-compose.dev.yml logs -f

dev-ps: ## Show running containers
	docker-compose -f docker-compose/docker-compose.dev.yml ps

dev-reset: dev-down clean-volumes dev-up ## Reset development environment (removes all data!)
	@echo "⚠️  Development environment has been reset with fresh data"

clean-volumes: ## Remove all data volumes (DESTRUCTIVE!)
	@echo "⚠️  WARNING: This will delete all data volumes!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker volume rm -f dictamesh-redpanda-data dictamesh-postgres-data dictamesh-redis-data dictamesh-prometheus-data dictamesh-grafana-data; \
		echo "✅ Volumes removed"; \
	else \
		echo "❌ Cancelled"; \
	fi

# Individual service utilities
kafka-topics: ## List all Kafka topics
	@echo "📋 Redpanda Topics:"
	docker exec -it dictamesh-redpanda rpk topic list

kafka-ui: ## Open Redpanda Console in browser
	@echo "Opening Redpanda Console..."
	@command -v xdg-open > /dev/null && xdg-open http://localhost:8080 || open http://localhost:8080 || echo "Please open http://localhost:8080 in your browser"

redis-cli: ## Connect to Redis CLI
	docker exec -it dictamesh-redis redis-cli

postgres-cli: ## Connect to PostgreSQL CLI
	docker exec -it dictamesh-postgres psql -U dictamesh -d metadata_catalog

# Monitoring
grafana: ## Open Grafana in browser
	@echo "Opening Grafana..."
	@command -v xdg-open > /dev/null && xdg-open http://localhost:3000 || open http://localhost:3000 || echo "Please open http://localhost:3000 in your browser"

prometheus: ## Open Prometheus in browser
	@echo "Opening Prometheus..."
	@command -v xdg-open > /dev/null && xdg-open http://localhost:9090 || open http://localhost:9090 || echo "Please open http://localhost:9090 in your browser"

jaeger: ## Open Jaeger UI in browser
	@echo "Opening Jaeger..."
	@command -v xdg-open > /dev/null && xdg-open http://localhost:16686 || open http://localhost:16686 || echo "Please open http://localhost:16686 in your browser"

# Health checks
health: ## Check health of all services
	@echo "🏥 Checking infrastructure health..."
	@echo ""
	@echo "Redpanda:"
	@docker exec dictamesh-redpanda rpk cluster health || echo "  ❌ Unhealthy"
	@echo ""
	@echo "PostgreSQL:"
	@docker exec dictamesh-postgres pg_isready -U dictamesh || echo "  ❌ Unhealthy"
	@echo ""
	@echo "Redis:"
	@docker exec dictamesh-redis redis-cli ping || echo "  ❌ Unhealthy"
	@echo ""
	@echo "All services:"
	@docker-compose -f docker-compose/docker-compose.dev.yml ps
