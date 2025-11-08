// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // Main documentation sidebar
  docsSidebar: [
    {
      type: 'category',
      label: 'Getting Started',
      collapsed: false,
      items: [
        'getting-started/introduction',
        'getting-started/quickstart',
        'getting-started/installation',
        'getting-started/core-concepts',
      ],
    },
    {
      type: 'category',
      label: 'Guides',
      collapsed: true,
      items: [
        'guides/building-adapters',
        'guides/graphql-federation',
        'guides/event-streaming',
        'guides/deployment',
        'guides/testing',
      ],
    },
    {
      type: 'category',
      label: 'Architecture',
      collapsed: true,
      items: [
        'architecture/overview',
        'architecture/core-framework',
        'architecture/connectors',
        'architecture/adapters',
        'architecture/services',
        'architecture/event-driven-integration',
        'architecture/metadata-catalog',
      ],
    },
    {
      type: 'category',
      label: 'API Reference',
      collapsed: true,
      items: [
        'api-reference/rest-api',
        'api-reference/graphql-api',
        'api-reference/go-packages',
        'api-reference/event-schemas',
      ],
    },
    {
      type: 'category',
      label: 'Operations',
      collapsed: true,
      items: [
        'operations/installation',
        'operations/configuration',
        'operations/monitoring',
        'operations/scaling',
        'operations/backup-restore',
        'operations/troubleshooting',
      ],
    },
    {
      type: 'category',
      label: 'Contributing',
      collapsed: true,
      items: [
        'contributing/contributing',
        'contributing/code-of-conduct',
        'contributing/development-setup',
      ],
    },
  ],
};

module.exports = sidebars;
