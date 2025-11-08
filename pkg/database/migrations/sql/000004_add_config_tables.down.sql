-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

--
-- Migration Rollback: Remove Configuration Management Tables
-- Description: Drops all tables, triggers, functions, and views created by 000004_add_config_tables.up.sql
--

-- Drop view
DROP VIEW IF EXISTS dictamesh_active_configs;

-- Drop functions
DROP FUNCTION IF EXISTS dictamesh_cleanup_stale_watchers(INTEGER);
DROP FUNCTION IF EXISTS dictamesh_cleanup_old_audit_logs(INTEGER);
DROP FUNCTION IF EXISTS dictamesh_get_config(VARCHAR, VARCHAR, VARCHAR, VARCHAR);
DROP FUNCTION IF EXISTS dictamesh_create_initial_config_version();
DROP FUNCTION IF EXISTS dictamesh_increment_config_version();

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS dictamesh_config_watchers;
DROP TABLE IF EXISTS dictamesh_encryption_keys;
DROP TABLE IF EXISTS dictamesh_config_audit_logs;
DROP TABLE IF EXISTS dictamesh_config_versions;
DROP TABLE IF EXISTS dictamesh_configurations;
