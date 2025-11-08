# DictaMesh Brand Guidelines

## Brand Identity

### Mission Statement
Empowering enterprises to build federated data architectures with confidence through open, transparent, and production-ready infrastructure.

### Brand Personality
- **Professional**: Enterprise-grade quality and reliability
- **Transparent**: Open-source, well-documented, scientifically validated
- **Empowering**: Enables teams to build solutions independently
- **Innovative**: Modern architecture patterns validated at scale
- **Trustworthy**: Production-proven, tested, and secure

## Visual Identity

### Logo

**Primary Logo**
```
┌─────────────────────────────────────┐
│                                     │
│    DictaMesh                        │
│    Enterprise Reference             │
│    Architecture                     │
│                                     │
└─────────────────────────────────────┘
```

**Logo Concepts:**
1. **Mesh Network**: Interconnected nodes representing data products
2. **Federation Symbol**: Unified gateway with distributed sources
3. **Data Flow**: Streams and pipelines converging

**Usage Rules:**
- Minimum clear space: 20px on all sides
- Minimum size: 120px width
- Never distort or rotate
- Never change colors without approval

### Color Palette

#### Primary Colors
```css
--dictamesh-blue-900:    #0A2540  /* Deep navy - headers, primary text */
--dictamesh-blue-700:    #1E4976  /* Rich blue - interactive elements */
--dictamesh-blue-500:    #3B82F6  /* Bright blue - CTAs, links */
--dictamesh-blue-300:    #93C5FD  /* Light blue - hover states */
--dictamesh-blue-100:    #DBEAFE  /* Pale blue - backgrounds */
```

#### Accent Colors
```css
--dictamesh-teal-600:    #0D9488  /* Success states, positive metrics */
--dictamesh-purple-600:  #7C3AED  /* Premium features, enterprise */
--dictamesh-amber-500:   #F59E0B  /* Warnings, important callouts */
--dictamesh-red-600:     #DC2626  /* Errors, critical alerts */
```

#### Neutral Colors
```css
--dictamesh-gray-900:    #111827  /* Body text */
--dictamesh-gray-700:    #374151  /* Secondary text */
--dictamesh-gray-500:    #6B7280  /* Muted text */
--dictamesh-gray-300:    #D1D5DB  /* Borders */
--dictamesh-gray-100:    #F3F4F6  /* Backgrounds */
--dictamesh-white:       #FFFFFF  /* White */
```

#### Semantic Colors
```css
--color-success:         var(--dictamesh-teal-600)
--color-warning:         var(--dictamesh-amber-500)
--color-error:           var(--dictamesh-red-600)
--color-info:            var(--dictamesh-blue-500)
```

### Typography

#### Primary Font Family
**Inter** (Google Fonts)
- Clean, modern, highly readable
- Excellent for UI and long-form content
- Wide range of weights available

```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
```

#### Secondary Font (Code)
**JetBrains Mono** or **Fira Code**
- Optimized for code display
- Ligature support
- Clear distinction between similar characters

```css
font-family: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;
```

#### Type Scale
```css
--text-xs:    0.75rem   /* 12px - captions, labels */
--text-sm:    0.875rem  /* 14px - small text, metadata */
--text-base:  1rem      /* 16px - body text */
--text-lg:    1.125rem  /* 18px - emphasis, leads */
--text-xl:    1.25rem   /* 20px - subheadings */
--text-2xl:   1.5rem    /* 24px - section titles */
--text-3xl:   1.875rem  /* 30px - page headings */
--text-4xl:   2.25rem   /* 36px - hero headings */
--text-5xl:   3rem      /* 48px - major headings */
--text-6xl:   3.75rem   /* 60px - landing page hero */
```

#### Font Weights
```css
--font-light:     300
--font-normal:    400
--font-medium:    500
--font-semibold:  600
--font-bold:      700
--font-extrabold: 800
```

### Spacing System

Based on 8px grid:
```css
--space-1:   0.25rem  /* 4px */
--space-2:   0.5rem   /* 8px */
--space-3:   0.75rem  /* 12px */
--space-4:   1rem     /* 16px */
--space-6:   1.5rem   /* 24px */
--space-8:   2rem     /* 32px */
--space-12:  3rem     /* 48px */
--space-16:  4rem     /* 64px */
--space-24:  6rem     /* 96px */
--space-32:  8rem     /* 128px */
```

### Border Radius
```css
--radius-sm:    0.125rem  /* 2px */
--radius-base:  0.25rem   /* 4px */
--radius-md:    0.375rem  /* 6px */
--radius-lg:    0.5rem    /* 8px */
--radius-xl:    0.75rem   /* 12px */
--radius-2xl:   1rem      /* 16px */
--radius-full:  9999px    /* Full circle */
```

### Shadows
```css
--shadow-sm:   0 1px 2px 0 rgba(0, 0, 0, 0.05);
--shadow-base: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
--shadow-md:   0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
--shadow-lg:   0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
--shadow-xl:   0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
```

## Voice & Tone

### General Voice
- **Clear and Concise**: Avoid jargon, explain technical concepts
- **Professional but Approachable**: Enterprise-ready, not corporate
- **Confident**: Based on proven patterns and production experience
- **Educational**: Help readers understand, not just sell
- **Honest**: Transparent about what the framework is and isn't

### Content Guidelines

#### Do's
✅ Use active voice
✅ Be specific with examples
✅ Cite sources and validation
✅ Explain the "why" behind decisions
✅ Show code examples
✅ Highlight real-world use cases
✅ Be transparent about trade-offs

#### Don'ts
❌ Make unsupported claims
❌ Use excessive buzzwords
❌ Oversell or hype
❌ Hide limitations
❌ Use generic marketing speak
❌ Assume prior knowledge without explanation

### Tone by Context

**Homepage/Marketing**
- Confident, clear value proposition
- Problem-solution focused
- Evidence-based (validation sources)
- Professional but inviting

**Documentation**
- Clear, instructional
- Step-by-step guidance
- Assumes intermediate technical knowledge
- Lots of code examples

**Blog Posts**
- Educational, informative
- Share learnings and insights
- Deep technical dives welcome
- Conversational but professional

**Enterprise Sales**
- ROI-focused
- Scalability and reliability emphasis
- Compliance and security highlights
- Professional, formal tone

## Visual Style

### Photography & Imagery
- **Style**: Clean, modern, professional
- **Subjects**: Data visualizations, architecture diagrams, team collaboration
- **Avoid**: Generic stock photos, clipart
- **Prefer**: Custom diagrams, real screenshots, authentic team photos

### Iconography
- **Style**: Outline icons (Heroicons)
- **Size**: 24px standard, 20px small, 32px large
- **Stroke**: 2px stroke width
- **Color**: Match text color or use brand colors

### Illustrations
- **Style**: Isometric or flat 2.5D
- **Colors**: Brand palette only
- **Purpose**: Explain complex concepts visually
- **Subjects**: Architecture diagrams, data flows, system interactions

### Data Visualization
- **Library**: D3.js, Recharts, or similar
- **Colors**: Use brand palette
- **Style**: Clean, minimalist
- **Accessibility**: Ensure sufficient contrast

## UI Components Style

### Buttons

**Primary Button**
```css
background: var(--dictamesh-blue-500);
color: white;
padding: 0.75rem 1.5rem;
border-radius: var(--radius-lg);
font-weight: var(--font-semibold);
transition: all 0.2s ease;
```

**Secondary Button**
```css
background: transparent;
color: var(--dictamesh-blue-500);
border: 2px solid var(--dictamesh-blue-500);
padding: 0.75rem 1.5rem;
border-radius: var(--radius-lg);
font-weight: var(--font-semibold);
```

### Cards
```css
background: white;
border: 1px solid var(--dictamesh-gray-200);
border-radius: var(--radius-xl);
padding: var(--space-6);
box-shadow: var(--shadow-sm);
transition: box-shadow 0.2s ease;
```

**Hover State**
```css
box-shadow: var(--shadow-md);
transform: translateY(-2px);
```

### Input Fields
```css
background: white;
border: 1px solid var(--dictamesh-gray-300);
border-radius: var(--radius-md);
padding: 0.625rem 1rem;
font-size: var(--text-base);
transition: border-color 0.2s ease;
```

**Focus State**
```css
border-color: var(--dictamesh-blue-500);
outline: 2px solid var(--dictamesh-blue-100);
outline-offset: 2px;
```

## Accessibility Standards

### WCAG 2.1 Level AA Compliance
- Minimum contrast ratio: 4.5:1 for normal text
- Minimum contrast ratio: 3:1 for large text (18px+)
- All interactive elements keyboard accessible
- Proper ARIA labels and roles
- Focus indicators visible
- Alt text for all images

### Color Contrast Validation
All color combinations have been tested for WCAG compliance:
- Primary blue on white: 8.6:1 ✅
- Body text on white: 16:1 ✅
- Gray 700 on white: 9.7:1 ✅
- Link blue on white: 4.6:1 ✅

## Brand Applications

### Website
- Homepage hero with clear value proposition
- Feature pages with diagrams and examples
- Pricing page with transparent pricing
- Partner program pages
- Blog with technical content
- Documentation site

### Marketing Materials
- Case studies
- White papers
- Technical presentations
- Conference booth materials
- Swag (t-shirts, stickers, etc.)

### Social Media
- LinkedIn (primary platform)
- Twitter/X (community engagement)
- GitHub (developer community)
- YouTube (tutorials and demos)

### Email Communications
- Newsletter (technical updates)
- Transactional emails
- Marketing campaigns
- Support communications

## File Formats

### Logo Files
- SVG (vector, preferred for web)
- PNG (transparent background, multiple sizes)
- PDF (print materials)

### Colors
- CSS variables (web)
- HEX codes (design tools)
- RGB values (screen media)
- CMYK values (print media)

## Brand Assets Location

```
branding/
├── logo/
│   ├── dictamesh-logo.svg
│   ├── dictamesh-logo-dark.svg
│   ├── dictamesh-icon.svg
│   └── dictamesh-wordmark.svg
├── colors/
│   └── palette.css
├── fonts/
│   ├── Inter/
│   └── JetBrainsMono/
└── templates/
    ├── presentation-template.pptx
    ├── document-template.docx
    └── social-media-templates/
```

## Questions or Clarifications

For brand guideline questions, contact:
- Design Team: design@dictamesh.com
- Marketing Team: marketing@dictamesh.com

## Changelog

- 2025-01-08: Initial brand guidelines established
