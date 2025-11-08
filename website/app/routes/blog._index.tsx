// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { Header } from "~/components/layout/Header";
import { Footer } from "~/components/layout/Footer";

export const meta: MetaFunction = () => {
  return [
    {
      title: "Blog | DictaMesh - Data Mesh Insights & Best Practices",
    },
    {
      name: "description",
      content:
        "Learn about data mesh architecture, event-driven integration, and production patterns from the DictaMesh team.",
    },
  ];
};

// Example blog posts - in production, these would come from a CMS or markdown files
const featuredPost = {
  title: "Building Production-Ready Data Mesh Adapters",
  slug: "building-production-ready-adapters",
  excerpt:
    "Learn the essential patterns and best practices for building resilient, observable, and maintainable data mesh adapters using the DictaMesh framework.",
  author: "Engineering Team",
  date: "2025-01-08",
  readTime: "12 min read",
  category: "Tutorial",
  image: "/og/blog-featured.png",
};

const posts = [
  {
    title: "Event-Driven Architecture with Kafka and Avro",
    slug: "event-driven-architecture-kafka-avro",
    excerpt:
      "Deep dive into designing event schemas, topic taxonomy, and handling schema evolution in a data mesh architecture.",
    author: "Platform Team",
    date: "2025-01-06",
    readTime: "10 min read",
    category: "Architecture",
  },
  {
    title: "Implementing Circuit Breakers for Resilient Adapters",
    slug: "circuit-breakers-resilient-adapters",
    excerpt:
      "Protect your data mesh from cascading failures with adaptive circuit breakers and exponential backoff strategies.",
    author: "DevOps Team",
    date: "2025-01-04",
    readTime: "8 min read",
    category: "Best Practices",
  },
  {
    title: "GraphQL Federation: Unified API Over Distributed Data",
    slug: "graphql-federation-unified-api",
    excerpt:
      "How DictaMesh uses Apollo Federation to compose a single GraphQL API from multiple domain adapters without N+1 queries.",
    author: "API Team",
    date: "2025-01-02",
    readTime: "15 min read",
    category: "Tutorial",
  },
  {
    title: "Data Governance and PII Tracking in the Mesh",
    slug: "data-governance-pii-tracking",
    excerpt:
      "Implement comprehensive data governance with automatic PII detection, access control, and audit logging.",
    author: "Compliance Team",
    date: "2024-12-30",
    readTime: "10 min read",
    category: "Governance",
  },
  {
    title: "Observability: Tracing Requests Across the Mesh",
    slug: "observability-distributed-tracing",
    excerpt:
      "Use OpenTelemetry to gain complete visibility into request flows from GraphQL gateway through adapters to source systems.",
    author: "SRE Team",
    date: "2024-12-28",
    readTime: "12 min read",
    category: "Observability",
  },
  {
    title: "Optimizing Cache Strategies for Sub-10ms Latency",
    slug: "cache-strategies-low-latency",
    excerpt:
      "Multi-layer caching patterns with L1 memory, L2 Redis, and intelligent invalidation to achieve millisecond response times.",
    author: "Performance Team",
    date: "2024-12-26",
    readTime: "9 min read",
    category: "Performance",
  },
];

const categories = ["All", "Tutorial", "Architecture", "Best Practices", "Governance", "Observability", "Performance"];

export default function BlogIndex() {
  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main>
        {/* Header */}
        <div className="bg-gradient-to-b from-dictamesh-blue-50 to-white py-16">
          <div className="container-custom">
            <div className="mx-auto max-w-2xl text-center">
              <h1 className="text-4xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-5xl">
                DictaMesh Blog
              </h1>
              <p className="mt-6 text-lg leading-8 text-gray-600">
                Insights, tutorials, and best practices for building data mesh architectures.
              </p>
            </div>
          </div>
        </div>

        {/* Featured Post */}
        <div className="section-padding bg-white">
          <div className="container-custom">
            <Link to={`/blog/${featuredPost.slug}`} className="group">
              <div className="card hover:shadow-2xl transition-all p-0 overflow-hidden">
                <div className="grid grid-cols-1 lg:grid-cols-2">
                  <div className="bg-gradient-to-br from-dictamesh-blue-500 to-dictamesh-purple-600 aspect-video lg:aspect-auto"></div>
                  <div className="p-8 lg:p-12">
                    <div className="flex items-center space-x-2 mb-4">
                      <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-dictamesh-blue-100 text-dictamesh-blue-800">
                        Featured
                      </span>
                      <span className="text-sm text-gray-600">{featuredPost.category}</span>
                    </div>
                    <h2 className="text-3xl font-bold text-gray-900 mb-4 group-hover:text-dictamesh-blue-600 transition-colors">
                      {featuredPost.title}
                    </h2>
                    <p className="text-gray-600 mb-6">{featuredPost.excerpt}</p>
                    <div className="flex items-center text-sm text-gray-500">
                      <span>{featuredPost.author}</span>
                      <span className="mx-2">·</span>
                      <time>{new Date(featuredPost.date).toLocaleDateString()}</time>
                      <span className="mx-2">·</span>
                      <span>{featuredPost.readTime}</span>
                    </div>
                  </div>
                </div>
              </div>
            </Link>
          </div>
        </div>

        {/* Category Filter */}
        <div className="border-b border-gray-200 bg-white sticky top-16 z-40">
          <div className="container-custom py-4">
            <div className="flex flex-wrap gap-2">
              {categories.map((category) => (
                <button
                  key={category}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    category === "All"
                      ? "bg-dictamesh-blue-600 text-white"
                      : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
                >
                  {category}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Blog Posts Grid */}
        <div className="section-padding bg-white">
          <div className="container-custom">
            <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
              {posts.map((post) => (
                <Link
                  key={post.slug}
                  to={`/blog/${post.slug}`}
                  className="card hover:shadow-lg transition-all p-0 overflow-hidden group"
                >
                  <div className="bg-gradient-to-br from-dictamesh-blue-400 to-dictamesh-teal-500 aspect-video"></div>
                  <div className="p-6">
                    <div className="flex items-center space-x-2 mb-3">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-dictamesh-blue-100 text-dictamesh-blue-800">
                        {post.category}
                      </span>
                    </div>
                    <h3 className="text-xl font-semibold text-gray-900 mb-3 group-hover:text-dictamesh-blue-600 transition-colors">
                      {post.title}
                    </h3>
                    <p className="text-gray-600 text-sm mb-4 line-clamp-3">{post.excerpt}</p>
                    <div className="flex items-center text-xs text-gray-500">
                      <span>{post.author}</span>
                      <span className="mx-2">·</span>
                      <time>{new Date(post.date).toLocaleDateString()}</time>
                      <span className="mx-2">·</span>
                      <span>{post.readTime}</span>
                    </div>
                  </div>
                </Link>
              ))}
            </div>

            {/* Load More */}
            <div className="mt-12 text-center">
              <button className="btn-secondary">
                Load More Posts
              </button>
            </div>
          </div>
        </div>

        {/* Newsletter Subscription */}
        <div className="bg-dictamesh-blue-50 section-padding">
          <div className="container-custom">
            <div className="mx-auto max-w-2xl text-center">
              <h3 className="text-2xl font-bold text-gray-900 mb-4">
                Subscribe to Our Newsletter
              </h3>
              <p className="text-gray-600 mb-8">
                Get the latest posts, tutorials, and data mesh insights delivered to your inbox.
              </p>
              <form className="flex flex-col sm:flex-row gap-4 max-w-md mx-auto">
                <input
                  type="email"
                  placeholder="Enter your email"
                  className="input flex-1"
                  required
                />
                <button type="submit" className="btn-primary whitespace-nowrap">
                  Subscribe
                </button>
              </form>
              <p className="text-xs text-gray-500 mt-4">
                No spam. Unsubscribe anytime. We respect your privacy.
              </p>
            </div>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
