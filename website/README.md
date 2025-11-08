# DictaMesh Website

Official website for DictaMesh - Enterprise-Grade Data Mesh Adapter Framework

## Technology Stack

- **Framework**: Remix (React-based SSR framework)
- **Admin/Dashboard**: Refine.dev integration
- **Styling**: Tailwind CSS
- **Animations**: Framer Motion
- **SEO**: Built-in Remix meta exports with Open Graph support
- **Blog**: Markdown-based with gray-matter
- **Analytics**: Ready for Google Analytics, Plausible, or self-hosted

## Project Structure

```
website/
├── app/                          # Remix application
│   ├── routes/                   # Page routes (file-based routing)
│   │   ├── _index.tsx           # Homepage
│   │   ├── pricing.tsx          # Pricing page
│   │   ├── partners.tsx         # Partner programs
│   │   ├── blog._index.tsx      # Blog index
│   │   ├── blog.$slug.tsx       # Blog posts
│   │   └── docs.$slug.tsx       # Documentation
│   ├── components/              # React components
│   │   ├── layout/             # Layout components
│   │   ├── marketing/          # Marketing sections
│   │   ├── ui/                 # UI primitives
│   │   └── seo/                # SEO components
│   ├── styles/                 # Global styles
│   ├── lib/                    # Utilities and libraries
│   └── utils/                  # Helper functions
├── public/                     # Static assets
│   ├── og/                    # Open Graph images
│   ├── images/                # Images and graphics
│   └── fonts/                 # Custom fonts
├── content/                   # Content (Markdown)
│   ├── blog/                  # Blog posts
│   └── docs/                  # Documentation
└── branding/                  # Brand assets and guidelines
```

## Features

### Enterprise Marketing
- Professional homepage with clear value proposition
- Detailed feature pages
- Use case demonstrations
- Customer testimonials (ready for integration)
- Trust indicators and certifications

### Hosted Cloud Services
- Pricing tiers (Free, Professional, Enterprise)
- Managed hosting options
- Support packages
- SLA guarantees

### Partner Programs
- Affiliate program
- Reseller partnerships
- Integration partners
- White-label offerings

### SEO & Performance
- Server-side rendering (SSR)
- Optimized meta tags
- Open Graph images
- Structured data (JSON-LD)
- Sitemap generation
- Fast page loads (<1s FCP)

### Blog Platform
- Markdown-based blog posts
- Syntax highlighting for code
- Author profiles
- Categories and tags
- RSS feed

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

## Deployment

The website is designed to deploy to:
- **Vercel** (recommended for Remix)
- **Netlify**
- **Cloudflare Pages**
- **Self-hosted** (Docker container included)

## Environment Variables

```bash
# Analytics
ANALYTICS_ID=

# Sentry (error tracking)
SENTRY_DSN=

# API endpoints (if using Refine admin)
API_URL=

# Newsletter/Email
NEWSLETTER_API_KEY=
```

## Brand Guidelines

See `branding/BRAND-GUIDELINES.md` for:
- Logo usage
- Color palette
- Typography
- Voice and tone
- Visual style

## License

SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
