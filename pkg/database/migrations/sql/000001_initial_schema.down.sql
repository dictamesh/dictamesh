-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS cache_status CASCADE;
DROP TABLE IF EXISTS data_lineage CASCADE;
DROP TABLE IF EXISTS event_log CASCADE;
DROP TABLE IF EXISTS schemas CASCADE;
DROP TABLE IF EXISTS entity_relationships CASCADE;
DROP TABLE IF EXISTS entity_catalog CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Note: We don't drop extensions as they might be used by other databases/schemas
