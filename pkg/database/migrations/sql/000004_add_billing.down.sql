-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- Migration Down: Remove billing system tables
-- Drops all billing tables in reverse order of dependencies

DROP TABLE IF EXISTS dictamesh_billing_audit_log CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_credits CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_pricing_tiers CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_payments CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_invoice_line_items CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_invoices CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_usage_metrics CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_subscriptions CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_subscription_plans CASCADE;
DROP TABLE IF EXISTS dictamesh_billing_organizations CASCADE;
