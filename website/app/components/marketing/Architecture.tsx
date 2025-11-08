// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

export function Architecture() {
  return (
    <div className="bg-dictamesh-blue-50 section-padding" id="architecture">
      <div className="container-custom">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-base font-semibold leading-7 text-dictamesh-blue-600">
            Proven Architecture
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-4xl">
            Built on Validated Enterprise Patterns
          </p>
          <p className="mt-6 text-lg leading-8 text-gray-600">
            DictaMesh synthesizes patterns validated at Fortune 500 scale:
            Data Mesh (ThoughtWorks), CQRS/Event Sourcing (Microsoft, AWS),
            Apollo Federation (Netflix, PayPal), and proven microservices architecture.
          </p>
        </div>

        <div className="mt-16 space-y-16">
          {/* Layered Architecture Diagram */}
          <div className="card bg-white p-8">
            <h3 className="text-xl font-semibold text-gray-900 mb-6">
              Layered Architecture
            </h3>
            <div className="space-y-4">
              {/* Services Layer */}
              <div className="border-2 border-dictamesh-purple-300 rounded-lg p-4 bg-dictamesh-purple-50">
                <div className="font-semibold text-dictamesh-purple-900 mb-2">
                  Services Layer
                </div>
                <div className="text-sm text-gray-600">
                  Your business applications: APIs, pipelines, workflows, AI/ML services
                </div>
              </div>

              {/* Arrow */}
              <div className="flex justify-center">
                <svg className="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </div>

              {/* Core Framework */}
              <div className="border-2 border-dictamesh-blue-500 rounded-lg p-4 bg-dictamesh-blue-50">
                <div className="font-semibold text-dictamesh-blue-900 mb-2">
                  Core Framework (DictaMesh Provides)
                </div>
                <div className="grid grid-cols-3 gap-2 mt-3 text-xs text-gray-700">
                  <div>GraphQL Gateway</div>
                  <div>Event Bus</div>
                  <div>Metadata Catalog</div>
                  <div>Observability</div>
                  <div>Governance</div>
                  <div>Resilience</div>
                </div>
              </div>

              {/* Arrow */}
              <div className="flex justify-center">
                <svg className="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </div>

              {/* Adapters Layer */}
              <div className="border-2 border-dictamesh-teal-400 rounded-lg p-4 bg-dictamesh-teal-50">
                <div className="font-semibold text-dictamesh-teal-900 mb-2">
                  Adapters Layer (You Build)
                </div>
                <div className="text-sm text-gray-600">
                  Domain-specific implementations: CMS, ERP, APIs, databases
                </div>
              </div>

              {/* Arrow */}
              <div className="flex justify-center">
                <svg className="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </div>

              {/* Connectors Layer */}
              <div className="border-2 border-gray-300 rounded-lg p-4 bg-gray-50">
                <div className="font-semibold text-gray-900 mb-2">
                  Connectors Layer
                </div>
                <div className="text-sm text-gray-600">
                  Protocol drivers: REST, GraphQL, gRPC, SOAP, PostgreSQL, MongoDB, etc.
                </div>
              </div>

              {/* Arrow */}
              <div className="flex justify-center">
                <svg className="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </div>

              {/* Data Sources */}
              <div className="border-2 border-gray-300 rounded-lg p-4 bg-white">
                <div className="font-semibold text-gray-900 mb-2">
                  Data Sources
                </div>
                <div className="text-sm text-gray-600">
                  Your existing systems: APIs, databases, files, legacy systems
                </div>
              </div>
            </div>
          </div>

          {/* Key Principles */}
          <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
            <div className="card">
              <h4 className="font-semibold text-lg text-gray-900 mb-3">
                Domain-Oriented
              </h4>
              <p className="text-sm text-gray-600">
                Each adapter owns its domain's data product. Teams maintain autonomy
                while framework provides consistency.
              </p>
            </div>
            <div className="card">
              <h4 className="font-semibold text-lg text-gray-900 mb-3">
                Event-Driven
              </h4>
              <p className="text-sm text-gray-600">
                Immutable event log enables real-time sync, audit trails, and
                time-travel queries. Source of truth always clear.
              </p>
            </div>
            <div className="card">
              <h4 className="font-semibold text-lg text-gray-900 mb-3">
                Self-Service
              </h4>
              <p className="text-sm text-gray-600">
                Platform team maintains framework. Domain teams build adapters
                independently. No bottlenecks.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
