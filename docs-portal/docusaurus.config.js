// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'DictaMesh',
  tagline: 'Enterprise-Grade Reference Architecture: Integration of Federated Authority Sources with Event-Driven Coordination',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://docs.dictamesh.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  organizationName: 'dictamesh',
  projectName: 'dictamesh',

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  markdown: {
    mermaid: true,
  },

  themes: ['@docusaurus/theme-mermaid'],

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/dictamesh/dictamesh/tree/main/docs-portal/',
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
          versions: {
            current: {
              label: 'v1.0 (Current)',
              path: '',
            },
          },
        },
        blog: {
          showReadingTime: true,
          blogTitle: 'DictaMesh Blog',
          blogDescription: 'Release notes, updates, and technical insights',
          postsPerPage: 'ALL',
          blogSidebarTitle: 'All posts',
          blogSidebarCount: 'ALL',
          editUrl: 'https://github.com/dictamesh/dictamesh/tree/main/docs-portal/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        // Sitemap
        sitemap: {
          changefreq: 'weekly',
          priority: 0.5,
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/dictamesh-social-card.jpg',

      // Navbar configuration
      navbar: {
        title: 'DictaMesh',
        logo: {
          alt: 'DictaMesh Logo',
          src: 'img/logo.svg',
          srcDark: 'img/logo-dark.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docsSidebar',
            position: 'left',
            label: 'Docs',
          },
          {
            to: '/docs/api-reference/rest-api',
            label: 'API',
            position: 'left',
          },
          {
            to: '/blog',
            label: 'Blog',
            position: 'left',
          },
          {
            type: 'docsVersionDropdown',
            position: 'right',
            dropdownActiveClassDisabled: true,
          },
          {
            href: 'https://github.com/dictamesh/dictamesh',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },

      // Footer configuration
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Documentation',
            items: [
              {
                label: 'Getting Started',
                to: '/docs/getting-started/introduction',
              },
              {
                label: 'Architecture',
                to: '/docs/architecture/overview',
              },
              {
                label: 'API Reference',
                to: '/docs/api-reference/rest-api',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub Discussions',
                href: 'https://github.com/dictamesh/dictamesh/discussions',
              },
              {
                label: 'Issues',
                href: 'https://github.com/dictamesh/dictamesh/issues',
              },
              {
                label: 'Contributing',
                to: '/docs/contributing',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/dictamesh/dictamesh',
              },
              {
                label: 'License',
                href: 'https://github.com/dictamesh/dictamesh/blob/main/LICENSE',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Controle Digital Ltda. Licensed under AGPL-3.0-or-later.`,
      },

      // Syntax highlighting
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ['bash', 'yaml', 'json', 'go', 'graphql', 'typescript'],
      },

      // Algolia search (will be configured later)
      // algolia: {
      //   appId: 'YOUR_APP_ID',
      //   apiKey: 'YOUR_SEARCH_API_KEY',
      //   indexName: 'dictamesh',
      //   contextualSearch: true,
      // },

      // Metadata
      metadata: [
        {name: 'keywords', content: 'data mesh, framework, graphql, federation, event-driven, kafka, go, enterprise'},
        {name: 'twitter:card', content: 'summary_large_image'},
      ],

      // Announcement bar
      announcementBar: {
        id: 'v1_0_release',
        content:
          '⭐️ If you like DictaMesh, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/dictamesh/dictamesh">GitHub</a>! ⭐️',
        backgroundColor: '#0066cc',
        textColor: '#ffffff',
        isCloseable: true,
      },

      // Color mode
      colorMode: {
        defaultMode: 'light',
        disableSwitch: false,
        respectPrefersColorScheme: true,
      },
    }),

  plugins: [
    // Future plugins can be added here
    // - Analytics plugin
    // - Search plugin
    // - Custom plugins
  ],
};

module.exports = config;
