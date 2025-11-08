-- SPDX-License-Identifier: AGPL-3.0-or-later
-- Copyright (C) 2025 Controle Digital Ltda

-- ClickHouse Initialization for Sentry
-- This script sets up the necessary tables for Sentry event storage

-- Create the sentry database if it doesn't exist
CREATE DATABASE IF NOT EXISTS sentry;

-- Use the sentry database
USE sentry;

-- Note: Sentry will automatically create the required tables when it starts
-- This file is provided for reference and can be extended with custom configurations

-- Optional: Create materialized views or custom indexes here if needed
