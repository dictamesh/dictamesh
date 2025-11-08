// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import { Link } from "@remix-run/react";

export function Hero() {
  return (
    <div className="relative isolate overflow-hidden bg-gradient-to-b from-dictamesh-blue-50 to-white">
      <div className="container-custom section-padding">
        <div className="mx-auto max-w-4xl text-center">
          {/* Badge */}
          <div className="mb-8 flex justify-center">
            <div className="relative rounded-full px-4 py-1.5 text-sm leading-6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20 transition-all">
              <span className="font-semibold text-dictamesh-blue-600">
                Open Source
              </span>{" "}
              <span className="inline-block h-1 w-1 rounded-full bg-gray-400 mx-2"></span>
              <span>Production-ready framework</span>
              <span className="inline-block h-1 w-1 rounded-full bg-gray-400 mx-2"></span>
              <span className="font-semibold">AGPL-3.0</span>
            </div>
          </div>

          {/* Headline */}
          <h1 className="text-5xl font-extrabold tracking-tight text-dictamesh-blue-900 sm:text-6xl lg:text-7xl mb-6">
            Build Your Data Mesh with{" "}
            <span className="gradient-text">Confidence</span>
          </h1>

          {/* Subheadline */}
          <p className="mx-auto mt-6 max-w-2xl text-lg leading-8 text-gray-600 sm:text-xl">
            Enterprise-grade framework for building federated data integrations.
            Integrate any data source with event-driven architecture, GraphQL
            federation, and built-in governance—all production-ready from day one.
          </p>

          {/* CTA Buttons */}
          <div className="mt-10 flex items-center justify-center gap-x-6">
            <Link to="/get-started" className="btn-primary text-lg">
              Get Started
              <svg
                className="ml-2 h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={2}
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3"
                />
              </svg>
            </Link>
            <Link to="/docs" className="btn-secondary text-lg">
              Read Documentation
            </Link>
          </div>

          {/* Trust Indicators */}
          <div className="mt-16 flex flex-wrap items-center justify-center gap-8 text-sm text-gray-600">
            <div className="flex items-center space-x-2">
              <svg
                className="h-5 w-5 text-dictamesh-teal-600"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                  clipRule="evenodd"
                />
              </svg>
              <span className="font-medium">Production-Proven</span>
            </div>
            <div className="flex items-center space-x-2">
              <svg
                className="h-5 w-5 text-dictamesh-teal-600"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                  clipRule="evenodd"
                />
              </svg>
              <span className="font-medium">Enterprise-Ready</span>
            </div>
            <div className="flex items-center space-x-2">
              <svg
                className="h-5 w-5 text-dictamesh-teal-600"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                  clipRule="evenodd"
                />
              </svg>
              <span className="font-medium">Fully Open Source</span>
            </div>
            <div className="flex items-center space-x-2">
              <svg
                className="h-5 w-5 text-dictamesh-teal-600"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                  clipRule="evenodd"
                />
              </svg>
              <span className="font-medium">Active Support</span>
            </div>
          </div>

          {/* Code Example */}
          <div className="mt-16 rounded-2xl bg-dictamesh-blue-900 p-8 shadow-2xl ring-1 ring-white/10">
            <div className="flex items-center justify-between mb-4">
              <div className="flex space-x-2">
                <div className="h-3 w-3 rounded-full bg-red-500"></div>
                <div className="h-3 w-3 rounded-full bg-yellow-500"></div>
                <div className="h-3 w-3 rounded-full bg-green-500"></div>
              </div>
              <span className="text-xs text-gray-400 font-mono">main.go</span>
            </div>
            <pre className="text-left text-sm text-gray-300 overflow-x-auto">
              <code className="font-mono">
{`// Build your adapter using DictaMesh
type CustomerAdapter struct {
    connector *rest.Connector
}

func (a *CustomerAdapter) GetEntity(
    ctx context.Context,
    id string,
) (*Entity, error) {
    // Framework provides:
    // ✓ Circuit breakers
    // ✓ Distributed tracing
    // ✓ Automatic caching
    // ✓ Event publishing
    // ✓ Governance policies

    return a.connector.Get(ctx, "/customers/" + id)
}

// Register with framework
app.RegisterAdapter("customers", customerAdapter)

// Framework automatically provides:
// → GraphQL API
// → Event streaming
// → Observability
// → Resilience patterns`}
              </code>
            </pre>
          </div>
        </div>
      </div>

      {/* Background decoration */}
      <div className="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80" aria-hidden="true">
        <div className="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-dictamesh-blue-500 to-dictamesh-purple-600 opacity-20 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"></div>
      </div>
    </div>
  );
}
