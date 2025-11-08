-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- DictaMesh Notifications Service - Rollback Migration
-- This migration removes all notifications service tables

-- Drop triggers
DROP TRIGGER IF EXISTS update_dictamesh_notification_templates_updated_at ON dictamesh_notification_templates;
DROP TRIGGER IF EXISTS update_dictamesh_notification_rules_updated_at ON dictamesh_notification_rules;
DROP TRIGGER IF EXISTS update_dictamesh_notifications_updated_at ON dictamesh_notifications;
DROP TRIGGER IF EXISTS update_dictamesh_notification_preferences_updated_at ON dictamesh_notification_preferences;
DROP TRIGGER IF EXISTS update_dictamesh_notification_rate_limits_updated_at ON dictamesh_notification_rate_limits;

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS dictamesh_notification_audit CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_rate_limits CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_batches CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_preferences CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_delivery CASCADE;
DROP TABLE IF EXISTS dictamesh_notifications CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_rules CASCADE;
DROP TABLE IF EXISTS dictamesh_notification_templates CASCADE;

-- Note: We don't drop the function dictamesh_update_updated_at_column()
-- as it might be used by other tables
