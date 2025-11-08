# Open Graph Images

This directory contains Open Graph (OG) images for social media sharing.

## Image Specifications

- **Size**: 1200x630 pixels
- **Format**: PNG or JPG
- **File size**: Under 1MB for optimal loading
- **Aspect ratio**: 1.91:1

## Required Images

1. `home.png` - Homepage OG image
2. `blog-featured.png` - Featured blog post image
3. `pricing.png` - Pricing page OG image
4. `partners.png` - Partners page OG image
5. `docs.png` - Documentation OG image

## Design Guidelines

### Branding
- Use DictaMesh color palette (blue-900, blue-500, teal-600)
- Include logo in bottom left corner
- Maintain 80px padding on all sides

### Typography
- Headline: Inter Bold, 64px
- Subheadline: Inter Medium, 32px
- Use white text on dark backgrounds
- Ensure WCAG AA contrast compliance

### Visual Elements
- Use mesh network patterns or data flow diagrams
- Include code snippets for technical content
- Add relevant icons (Heroicons style)

## Tools for Creation

- **Figma**: Use provided templates in `/design/og-templates.fig`
- **Canva**: Import brand assets and use custom dimensions
- **Code-based**: Generate programmatically with Puppeteer/Playwright

## Dynamic OG Images (Future)

Consider implementing dynamic OG image generation for:
- Blog posts (title, author, date)
- Documentation pages (section title)
- Custom landing pages

Example services:
- Vercel OG Image Generation
- Cloudinary transformations
- Self-hosted with Puppeteer

## SEO Best Practices

- Include relevant keywords in image filename
- Add descriptive alt text in meta tags
- Ensure images load quickly (<1s)
- Test with [Open Graph Debugger](https://www.opengraph.xyz/)

## Testing

Test OG images on:
- [Facebook Debugger](https://developers.facebook.com/tools/debug/)
- [Twitter Card Validator](https://cards-dev.twitter.com/validator)
- [LinkedIn Post Inspector](https://www.linkedin.com/post-inspector/)

## Template Structure

```
┌────────────────────────────────────────────────┐
│  [80px padding]                                │
│                                                │
│     [Headline Text - Bold, Large]              │
│     [Subheadline - Medium weight]              │
│                                                │
│     [Visual Element or Code Snippet]           │
│                                                │
│                                                │
│  [Logo]                            [80px]      │
└────────────────────────────────────────────────┘
```

## Placeholder Images

Until custom images are created, use:
1. Gradient backgrounds (blue-500 to purple-600)
2. DictaMesh logo centered
3. Page title in white text

Generate placeholders with:
```bash
# Using ImageMagick
convert -size 1200x630 gradient:blue-purple \
  -gravity center -pointsize 72 -fill white \
  -annotate +0+0 "DictaMesh" \
  home.png
```
