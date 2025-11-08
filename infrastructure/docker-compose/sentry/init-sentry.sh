#!/bin/bash
# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2025 Controle Digital Ltda

#
# Sentry Initialization Script
# This script initializes the Sentry database and creates a default superuser
#

set -e

echo "🔧 Initializing Sentry for DictaMesh..."

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL..."
until PGPASSWORD=$SENTRY_DB_PASSWORD psql -h "$SENTRY_POSTGRES_HOST" -U "$SENTRY_DB_USER" -d "$SENTRY_DB_NAME" -c '\q' 2>/dev/null; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done
echo "✅ PostgreSQL is ready"

# Wait for Redis to be ready
echo "⏳ Waiting for Redis..."
until redis-cli -h "$SENTRY_REDIS_HOST" -p "$SENTRY_REDIS_PORT" ping 2>/dev/null; do
  echo "Redis is unavailable - sleeping"
  sleep 2
done
echo "✅ Redis is ready"

# Run database migrations
echo "🔄 Running database migrations..."
sentry upgrade --noinput

# Create default superuser if it doesn't exist
echo "👤 Creating default superuser..."
sentry createuser \
  --email "${SENTRY_ADMIN_EMAIL:-admin@dictamesh.local}" \
  --password "${SENTRY_ADMIN_PASSWORD:-admin}" \
  --superuser \
  --no-input || echo "Superuser already exists"

# Create default organization and team
echo "🏢 Setting up default organization..."
python3 << END
from sentry.runner import configure
configure()

from sentry.models import Organization, Team, OrganizationMember, User

# Get or create the default organization
org, created = Organization.objects.get_or_create(
    slug='dictamesh',
    defaults={
        'name': 'DictaMesh Framework',
    }
)

if created:
    print(f"✅ Created organization: {org.name}")
else:
    print(f"ℹ️  Organization already exists: {org.name}")

# Get or create the default team
team, created = Team.objects.get_or_create(
    organization=org,
    slug='framework-team',
    defaults={
        'name': 'Framework Team',
    }
)

if created:
    print(f"✅ Created team: {team.name}")
else:
    print(f"ℹ️  Team already exists: {team.name}")

# Add admin user to the organization
try:
    admin_user = User.objects.get(email='${SENTRY_ADMIN_EMAIL:-admin@dictamesh.local}')
    member, created = OrganizationMember.objects.get_or_create(
        organization=org,
        user=admin_user,
        defaults={
            'role': 'owner',
        }
    )
    if created:
        print(f"✅ Added {admin_user.email} to organization")
    else:
        print(f"ℹ️  {admin_user.email} already in organization")
except User.DoesNotExist:
    print("⚠️  Admin user not found")
END

echo ""
echo "✅ Sentry initialization complete!"
echo ""
echo "You can now access Sentry at: http://localhost:9000"
echo "Default credentials:"
echo "  Email:    ${SENTRY_ADMIN_EMAIL:-admin@dictamesh.local}"
echo "  Password: ${SENTRY_ADMIN_PASSWORD:-admin}"
echo ""
