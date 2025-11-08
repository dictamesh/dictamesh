# Documentation Portal Planning

[← Previous: Contribution Guidelines](20-CONTRIBUTION-GUIDELINES.md)

---

## 🎯 Purpose

This document provides comprehensive planning and stack decisions for the DictaMesh documentation portal. This extends the documentation planning (14-DOCUMENTATION-PLANNING.md) by defining the infrastructure, framework, and tooling for a production-ready documentation website.

**Reading Time:** 25 minutes
**Prerequisites:** 14-DOCUMENTATION-PLANNING.md
**Outputs:** Complete documentation portal infrastructure, stack decisions, implementation guide

---

## 📊 Executive Summary

### Portal Objectives

1. **User-Friendly**: Beautiful, modern UI with excellent UX
2. **Developer-Focused**: Code examples, API references, interactive demos
3. **Searchable**: Fast, accurate search across all documentation
4. **Versioned**: Support multiple versions (current, next, legacy)
5. **Accessible**: WCAG 2.1 AA compliant
6. **Fast**: Sub-second page loads, optimized for performance
7. **Maintainable**: Easy to update, automated workflows

### Stack Decision Summary

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Framework** | Docusaurus 3 | Best-in-class versioning, React-based, extensive plugin ecosystem |
| **Search** | Algolia DocSearch (free tier) | Production-ready, fast, excellent UX |
| **API Docs** | Multi-tool approach | Swagger UI (REST), GraphQL Playground, godoc for Go packages |
| **Diagrams** | Mermaid.js + Excalidraw | Built-in support, version-controlled diagrams |
| **Hosting** | GitHub Pages / Vercel | Free, reliable, automatic deployments |
| **Analytics** | Plausible (self-hosted) | Privacy-friendly, GDPR compliant |
| **Feedback** | GitHub Discussions integration | Community-driven feedback loop |

---

## 🏗️ Architecture Decision Records (ADRs)

### ADR-001: Documentation Framework Selection

**Status**: Accepted

**Context**: Need a modern, maintainable documentation framework supporting versioning, search, and API documentation.

**Candidates Evaluated**:

1. **Docusaurus 3** (Meta/Facebook)
   - ✅ Excellent versioning support
   - ✅ Large ecosystem, active development
   - ✅ Built-in search integration
   - ✅ MDX support (React components in markdown)
   - ✅ Optimized for performance (React SSG)
   - ✅ Plugin system for extensibility
   - ❌ React learning curve for contributors

2. **VitePress** (Vue.js)
   - ✅ Extremely fast (Vite-powered)
   - ✅ Simple, lightweight
   - ✅ Vue 3 based
   - ❌ Less mature ecosystem
   - ❌ Limited versioning support
   - ❌ Fewer plugins

3. **Astro Starlight**
   - ✅ Ultra-fast (partial hydration)
   - ✅ Framework-agnostic
   - ✅ Modern, clean design
   - ❌ Newer, smaller community
   - ❌ Limited versioning
   - ❌ Fewer integrations

4. **MkDocs Material**
   - ✅ Beautiful design
   - ✅ Great search
   - ✅ Python-based (familiar for data engineers)
   - ❌ No versioning without manual work
   - ❌ Limited interactivity
   - ❌ Less modern architecture

**Decision**: **Docusaurus 3**

**Rationale**:
- Versioning is critical for a framework project (support v1.x, v2.x simultaneously)
- Large open-source projects (React, Jest, Redux) use Docusaurus successfully
- MDX enables interactive code examples and embedded components
- Strong plugin ecosystem for API docs, search, analytics
- Performance is excellent with SSG + partial hydration

---

### ADR-002: Search Solution

**Status**: Accepted

**Context**: Need fast, accurate search across documentation, API references, and code examples.

**Candidates**:

1. **Algolia DocSearch** (Free for Open Source)
   - ✅ Production-ready, fast
   - ✅ Excellent UI/UX
   - ✅ Free for OSS projects
   - ✅ Automatic indexing
   - ❌ External dependency

2. **Meilisearch** (Self-hosted)
   - ✅ Fast, typo-tolerant
   - ✅ Self-hosted, privacy-friendly
   - ✅ Good API
   - ❌ Requires hosting/maintenance
   - ❌ More complex setup

3. **Local Search Plugin**
   - ✅ No external dependencies
   - ✅ Works offline
   - ❌ Limited features
   - ❌ Slower for large docs
   - ❌ No analytics

**Decision**: **Algolia DocSearch** (primary) + Local Search (fallback)

**Rationale**:
- Algolia free tier perfect for OSS
- Best-in-class search UX
- Local search as fallback for offline use
- Can migrate to self-hosted later if needed

---

### ADR-003: API Documentation Strategy

**Status**: Accepted

**Context**: Need to document multiple API types: REST APIs, GraphQL APIs, Go packages.

**Solution**: Multi-tool approach based on API type

| API Type | Tool | Integration Method |
|----------|------|-------------------|
| **REST APIs** | Swagger UI / Redoc | OpenAPI 3.0 specs auto-generated from code |
| **GraphQL APIs** | GraphQL Playground / GraphiQL | Schema introspection, embedded in portal |
| **Go Packages** | pkgsite (godoc) | Hosted separately, linked from main portal |
| **Event Schemas** | Avro Schema Viewer | Custom React component in Docusaurus |

**Implementation**:
```
docs-portal/
├── docs/
│   └── api/
│       ├── rest/
│       │   └── openapi.yaml (auto-generated)
│       ├── graphql/
│       │   └── schema.graphql (introspected)
│       └── events/
│           └── schemas/ (Avro schemas)
├── src/
│   ├── components/
│   │   ├── ApiExplorer.tsx (REST)
│   │   ├── GraphQLExplorer.tsx
│   │   └── SchemaViewer.tsx (Avro)
│   └── pages/
│       └── api.tsx
```

---

### ADR-004: Diagram and Visualization Strategy

**Status**: Accepted

**Tools**:

1. **Mermaid.js** (Built-in to Docusaurus)
   - Use for: Sequence diagrams, flowcharts, architecture diagrams
   - Benefits: Version-controlled, renders as code
   - Example:
   ```mermaid
   graph TD
     A[Client] -->|GraphQL Query| B[Gateway]
     B -->|Batched Requests| C[Adapters]
     C -->|Event Stream| D[Kafka]
   ```

2. **Excalidraw** (Plugin available)
   - Use for: Hand-drawn style diagrams, wireframes
   - Benefits: Interactive, exportable to SVG

3. **D3.js** (Custom React components)
   - Use for: Interactive data visualizations
   - Example: Relationship graph explorer

**Decision**: Use all three based on use case

---

### ADR-005: Versioning Strategy

**Status**: Accepted

**Versions to Maintain**:
- `current`: Latest stable (v1.0, v1.1, etc.)
- `next`: Development version (main branch)
- `v0.x`: Legacy (if needed)

**URL Structure**:
```
docs.dictamesh.controle.digital/        → current version
docs.dictamesh.controle.digital/next/   → next version
docs.dictamesh.controle.digital/v0.x/   → legacy
```

**Docusaurus Config**:
```js
module.exports = {
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          versions: {
            current: {
              label: 'v1.0 (Current)',
              path: '',
            },
            next: {
              label: 'Next',
              path: 'next',
            },
          },
        },
      },
    ],
  ],
}
```

---

## 🎨 Design System

### Color Palette

**Light Mode**:
```css
--primary: #0066cc;      /* DictaMesh Blue */
--secondary: #6b46c1;    /* Purple accent */
--success: #10b981;      /* Green */
--warning: #f59e0b;      /* Amber */
--danger: #ef4444;       /* Red */
--text: #1a202c;         /* Dark gray */
--background: #ffffff;   /* White */
```

**Dark Mode**:
```css
--primary: #3b82f6;      /* Lighter blue */
--secondary: #8b5cf6;    /* Lighter purple */
--text: #e2e8f0;         /* Light gray */
--background: #0f172a;   /* Dark blue-gray */
```

### Typography

- **Headings**: Inter (system font fallback)
- **Body**: Inter
- **Code**: JetBrains Mono

### Components

1. **Hero Section** (Homepage)
   - Large, clear value proposition
   - Quick start CTA
   - Code example preview
   - Architecture diagram

2. **Navigation**
   - Top navbar: Logo, Docs, API, Community, GitHub
   - Sidebar: Nested documentation structure
   - Version dropdown
   - Search bar (prominent)
   - Theme toggle (dark/light)

3. **Code Blocks**
   - Syntax highlighting (Prism)
   - Copy button
   - Line numbers
   - Language badges
   - Runnable examples (CodeSandbox integration)

4. **API Explorer**
   - Interactive REST API explorer
   - GraphQL query builder
   - Response previews
   - Authentication helpers

5. **Feedback Widget**
   - "Was this page helpful?" (Yes/No)
   - GitHub edit link
   - Report issue link

---

## 📁 Portal Structure

```
dictamesh-docs/
├── docs/                       # Documentation content
│   ├── getting-started/
│   │   ├── introduction.md
│   │   ├── quickstart.md
│   │   ├── installation.md
│   │   └── core-concepts.md
│   ├── guides/
│   │   ├── building-adapters.md
│   │   ├── graphql-federation.md
│   │   ├── event-streaming.md
│   │   └── deployment.md
│   ├── api-reference/
│   │   ├── rest-api.md
│   │   ├── graphql-api.md
│   │   ├── go-packages.md
│   │   └── event-schemas.md
│   ├── architecture/
│   │   ├── overview.md
│   │   ├── core-framework.md
│   │   ├── connectors.md
│   │   ├── adapters.md
│   │   └── services.md
│   ├── operations/
│   │   ├── installation.md
│   │   ├── configuration.md
│   │   ├── monitoring.md
│   │   ├── scaling.md
│   │   └── troubleshooting.md
│   └── contributing/
│       ├── contributing.md
│       ├── code-of-conduct.md
│       └── development-setup.md
│
├── blog/                       # Release notes, updates
│   ├── 2025-01-15-v1.0-release.md
│   └── 2025-02-01-roadmap-2025.md
│
├── src/
│   ├── components/             # React components
│   │   ├── HomepageFeatures/
│   │   ├── ApiExplorer/
│   │   ├── GraphQLExplorer/
│   │   └── SchemaViewer/
│   ├── css/
│   │   ├── custom.css
│   │   └── dark-mode.css
│   └── pages/
│       ├── index.tsx           # Homepage
│       ├── api.tsx             # API explorer page
│       └── community.tsx       # Community page
│
├── static/                     # Static assets
│   ├── img/
│   │   ├── logo.svg
│   │   ├── architecture-diagram.svg
│   │   └── screenshots/
│   ├── schemas/                # Avro schemas
│   └── openapi/                # OpenAPI specs
│
├── docusaurus.config.js        # Docusaurus configuration
├── sidebars.js                 # Sidebar structure
├── package.json
├── tsconfig.json
└── README.md
```

---

## 🔧 Implementation Phases

### Phase 1: Foundation (Week 1)

**Goal**: Basic Docusaurus site with core documentation structure

**Tasks**:
1. ✅ Initialize Docusaurus project
2. ✅ Configure basic theme and branding
3. ✅ Set up GitHub Pages deployment
4. ✅ Migrate existing markdown docs
5. ✅ Configure sidebar navigation
6. ✅ Set up dark mode

**Deliverables**:
- Live documentation site at `docs.dictamesh.controle.digital`
- Basic navigation and structure
- Migrated content from `/docs` directory

---

### Phase 2: Content & Styling (Week 2)

**Goal**: Enhanced UI, complete content migration

**Tasks**:
1. ✅ Implement custom theme (colors, fonts)
2. ✅ Create homepage with hero section
3. ✅ Add code block enhancements (copy button, line numbers)
4. ✅ Migrate all planning docs to user-facing format
5. ✅ Add Mermaid.js diagram support
6. ✅ Create footer with links

**Deliverables**:
- Polished, branded documentation site
- All existing content migrated and formatted
- Interactive diagrams

---

### Phase 3: API Documentation (Week 3)

**Goal**: Integrated API reference and explorers

**Tasks**:
1. ✅ Generate OpenAPI specs from Go code (swag)
2. ✅ Integrate Swagger UI component
3. ✅ Set up GraphQL Playground
4. ✅ Create Avro schema viewer component
5. ✅ Link to godoc for Go packages
6. ✅ Add interactive code examples

**Deliverables**:
- Complete API reference
- Interactive API explorers
- Code examples that can be tested

---

### Phase 4: Search & Discovery (Week 4)

**Goal**: Production-ready search and navigation

**Tasks**:
1. ✅ Apply for Algolia DocSearch (free tier)
2. ✅ Configure Algolia crawler
3. ✅ Add local search fallback
4. ✅ Implement search analytics
5. ✅ Add "Suggest Edit" functionality
6. ✅ Set up feedback widget

**Deliverables**:
- Fast, accurate search
- Community contribution workflow
- User feedback mechanism

---

### Phase 5: Advanced Features (Week 5)

**Goal**: Versioning, analytics, and automation

**Tasks**:
1. ✅ Configure versioning (current, next)
2. ✅ Set up Plausible analytics
3. ✅ Automate OpenAPI spec generation (CI/CD)
4. ✅ Add auto-generated changelog from commits
5. ✅ Set up broken link checker
6. ✅ Configure automatic deployments

**Deliverables**:
- Multi-version support
- Analytics and insights
- Automated documentation updates

---

## 🚀 Deployment Strategy

### Hosting Options

**Primary: GitHub Pages** (Free)
- ✅ Free for public repos
- ✅ Automatic deployments with GitHub Actions
- ✅ Custom domain support
- ✅ HTTPS included
- ❌ Limited to static sites (perfect for Docusaurus)

**Alternative: Vercel** (Free tier)
- ✅ Excellent performance (global CDN)
- ✅ Preview deployments for PRs
- ✅ Analytics included
- ✅ Easy custom domains

**Decision**: **GitHub Pages** (primary), Vercel (backup)

### CI/CD Workflow

```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Documentation

on:
  push:
    branches: [main]
    paths:
      - 'docs-portal/**'
      - 'services/**/*.go' # Re-generate API docs

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: 20

      - name: Install dependencies
        working-directory: ./docs-portal
        run: npm ci

      - name: Generate OpenAPI specs
        run: |
          go install github.com/swaggo/swag/cmd/swag@latest
          cd services/metadata-catalog
          swag init -g cmd/server/main.go -o ../../docs-portal/static/openapi

      - name: Build documentation
        working-directory: ./docs-portal
        run: npm run build

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs-portal/build
          cname: docs.dictamesh.controle.digital
```

### Custom Domain Setup

1. Add `CNAME` file in `static/` directory:
   ```
   docs.dictamesh.controle.digital
   ```

2. Configure DNS:
   ```
   CNAME docs.dictamesh.controle.digital -> click2-run.github.io
   ```

3. Enable HTTPS in GitHub Pages settings

---

## 🔍 Search Configuration

### Algolia DocSearch Setup

1. **Apply for Free Tier**:
   - URL: https://docsearch.algolia.com/apply/
   - Requirements: Public repo, documentation site, open source

2. **Configuration File** (`.algolia/config.json`):
```json
{
  "index_name": "dictamesh",
  "start_urls": ["https://docs.dictamesh.controle.digital/"],
  "sitemap_urls": ["https://docs.dictamesh.controle.digital/sitemap.xml"],
  "selectors": {
    "lvl0": {
      "selector": ".menu__link--sublist.menu__link--active",
      "global": true
    },
    "lvl1": "header h1",
    "lvl2": "article h2",
    "lvl3": "article h3",
    "lvl4": "article h4",
    "text": "article p, article li"
  }
}
```

3. **Docusaurus Integration**:
```js
// docusaurus.config.js
module.exports = {
  themeConfig: {
    algolia: {
      appId: 'YOUR_APP_ID',
      apiKey: 'YOUR_SEARCH_API_KEY',
      indexName: 'dictamesh',
      contextualSearch: true,
      searchParameters: {},
    },
  },
}
```

---

## 📊 Analytics & Monitoring

### Plausible Analytics (Self-hosted)

**Why Plausible**:
- ✅ Privacy-friendly (no cookies, GDPR compliant)
- ✅ Lightweight (< 1KB script)
- ✅ Self-hosted option available
- ✅ Simple, clear metrics

**Metrics to Track**:
- Page views
- Most visited pages
- Search queries
- Referrers
- Geographic distribution
- Device/browser breakdown

**Integration**:
```html
<!-- In docusaurus.config.js -->
scripts: [
  {
    src: 'https://plausible.dictamesh.controle.digital/js/script.js',
    'data-domain': 'docs.dictamesh.controle.digital',
    defer: true,
  },
]
```

---

## 🧪 Quality Assurance

### Automated Checks

1. **Broken Link Checker**:
```yaml
# .github/workflows/link-check.yml
- name: Check links
  uses: gaurav-nelson/github-action-markdown-link-check@v1
  with:
    config-file: '.github/workflows/link-check-config.json'
```

2. **Accessibility Check** (Lighthouse CI):
```yaml
- name: Lighthouse CI
  uses: treosh/lighthouse-ci-action@v9
  with:
    urls: |
      https://docs.dictamesh.controle.digital
    uploadArtifacts: true
    temporaryPublicStorage: true
```

3. **Spell Check**:
```yaml
- name: Spell check
  uses: rojopolis/spellcheck-github-actions@v0
  with:
    config_path: .spellcheck.yml
```

### Manual Review Checklist

Before releasing new documentation version:
- [ ] All links work
- [ ] Code examples tested
- [ ] Screenshots up-to-date
- [ ] API references generated
- [ ] Version numbers correct
- [ ] Diagrams render correctly
- [ ] Dark mode works
- [ ] Mobile responsive
- [ ] Search returns relevant results
- [ ] Load time < 2s

---

## 🎯 Success Metrics

### KPIs

1. **Adoption Metrics**:
   - Unique visitors per month
   - Page views per month
   - Average time on site
   - Bounce rate < 60%

2. **Content Metrics**:
   - Top 10 most visited pages
   - Search queries (trending topics)
   - 404 errors (broken links)

3. **Engagement Metrics**:
   - GitHub edits suggested (community contributions)
   - Feedback widget responses
   - External links clicked

4. **Technical Metrics**:
   - Page load time (target: < 2s)
   - Lighthouse score (target: > 90)
   - Uptime (target: 99.9%)

### Quarterly Review

- Analyze top search queries → add missing content
- Review most visited pages → ensure up-to-date
- Check 404 errors → fix broken links
- Survey users → gather feedback

---

## 🗺️ Content Roadmap

### Phase 1: Core Documentation (Q1 2025)
- ✅ Getting Started
- ✅ Architecture Overview
- ✅ API Reference
- ✅ Deployment Guide

### Phase 2: Advanced Topics (Q2 2025)
- 🔲 Building Custom Connectors
- 🔲 Advanced GraphQL Patterns
- 🔲 Performance Tuning
- 🔲 Security Best Practices

### Phase 3: Ecosystem Documentation (Q3 2025)
- 🔲 Example Adapters (CMS, E-commerce, ERP)
- 🔲 Integration Patterns
- 🔲 Case Studies
- 🔲 Video Tutorials

### Phase 4: Community Content (Q4 2025)
- 🔲 Community Adapters Showcase
- 🔲 Blog Posts from Contributors
- 🔲 Conference Talks
- 🔲 Workshops and Training Materials

---

## 🛠️ Maintenance Plan

### Weekly Tasks
- [ ] Review and merge community contributions
- [ ] Update changelog
- [ ] Monitor analytics for trends
- [ ] Check for broken links

### Monthly Tasks
- [ ] Generate and publish API docs
- [ ] Update screenshots if UI changed
- [ ] Review feedback widget responses
- [ ] Analyze search queries for content gaps

### Quarterly Tasks
- [ ] Full documentation audit
- [ ] Update version compatibility matrix
- [ ] Refresh getting started guide
- [ ] User survey

### Yearly Tasks
- [ ] Major documentation restructure (if needed)
- [ ] Design refresh
- [ ] Framework upgrade (Docusaurus, etc.)
- [ ] Comprehensive accessibility audit

---

## 📚 Resources & References

### Docusaurus Resources
- [Official Documentation](https://docusaurus.io/)
- [Showcase](https://docusaurus.io/showcase) - Examples of great sites
- [Plugin Marketplace](https://docusaurus.io/community/resources#plugins)

### Design Inspiration
- [Stripe Docs](https://stripe.com/docs) - API reference excellence
- [React Docs](https://react.dev/) - Interactive examples
- [Next.js Docs](https://nextjs.org/docs) - Clean design
- [Supabase Docs](https://supabase.com/docs) - Great search UX

### Tools & Plugins
- [Swagger UI](https://swagger.io/tools/swagger-ui/) - REST API docs
- [GraphQL Playground](https://github.com/graphql/graphql-playground) - GraphQL explorer
- [Mermaid.js](https://mermaid.js.org/) - Diagrams
- [Prism](https://prismjs.com/) - Syntax highlighting

---

## ✅ Implementation Checklist

### Initial Setup
- [ ] Create `docs-portal/` directory
- [ ] Initialize Docusaurus project
- [ ] Configure custom domain
- [ ] Set up GitHub Actions for deployment
- [ ] Apply for Algolia DocSearch

### Content Migration
- [ ] Migrate README.md content
- [ ] Migrate PROJECT-SCOPE.md
- [ ] Migrate planning docs to user-facing format
- [ ] Create getting started guide
- [ ] Add code examples

### API Documentation
- [ ] Set up OpenAPI spec generation
- [ ] Integrate Swagger UI
- [ ] Configure GraphQL Playground
- [ ] Create Avro schema viewer
- [ ] Link to godoc

### Styling & UX
- [ ] Implement custom theme
- [ ] Configure dark mode
- [ ] Add hero section
- [ ] Create custom components
- [ ] Add feedback widget

### Search & Discovery
- [ ] Configure Algolia
- [ ] Add local search fallback
- [ ] Implement search analytics
- [ ] Test search accuracy

### Quality Assurance
- [ ] Set up link checker
- [ ] Configure Lighthouse CI
- [ ] Add spell checker
- [ ] Test mobile responsive
- [ ] Verify accessibility

### Launch
- [ ] Final review
- [ ] Announce on GitHub
- [ ] Share with community
- [ ] Collect initial feedback

---

[← Previous: Contribution Guidelines](20-CONTRIBUTION-GUIDELINES.md)

---

**Document Metadata**
- Version: 1.0.0
- Last Updated: 2025-11-08
- Author: DictaMesh Documentation Team
- Status: Planning → Implementation
