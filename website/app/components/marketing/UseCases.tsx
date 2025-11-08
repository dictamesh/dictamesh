// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

const useCases = [
  {
    title: "E-commerce Integration",
    description:
      "Unify product catalogs, orders, customers, and inventory from CMS, payment gateways, and ERPs into a single data mesh.",
    tech: ["REST APIs", "PostgreSQL", "Event Streaming"],
    icon: "🛒",
  },
  {
    title: "Healthcare Data Exchange",
    description:
      "HIPAA-compliant integration of EHR systems, lab results, pharmacy data, and insurance claims with built-in governance.",
    tech: ["HL7/FHIR", "SOAP", "PII Tracking"],
    icon: "🏥",
  },
  {
    title: "Financial Services",
    description:
      "Integrate trading systems, risk engines, customer data, and regulatory reporting with SOC 2 compliance.",
    tech: ["Oracle", "Mainframes", "Audit Logs"],
    icon: "🏦",
  },
  {
    title: "SaaS Platform",
    description:
      "Connect CRM, analytics, customer databases, and third-party APIs into unified customer 360° view.",
    tech: ["GraphQL", "Multi-tenant DB", "OpenAPI"],
    icon: "☁️",
  },
  {
    title: "IoT & Telemetry",
    description:
      "Real-time ingestion of sensor data, device telemetry, and environmental monitoring with time-series optimization.",
    tech: ["MQTT", "Time-series DB", "Stream Processing"],
    icon: "📡",
  },
  {
    title: "Legacy Modernization",
    description:
      "Gradually extract data from COBOL mainframes, AS/400, and legacy databases without big-bang migration.",
    tech: ["ODBC", "File Transfer", "Change Data Capture"],
    icon: "🏛️",
  },
];

export function UseCases() {
  return (
    <div className="bg-white section-padding" id="use-cases">
      <div className="container-custom">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-base font-semibold leading-7 text-dictamesh-blue-600">
            Use Cases
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-4xl">
            Proven Across Industries
          </p>
          <p className="mt-6 text-lg leading-8 text-gray-600">
            From startups to Fortune 500 enterprises, DictaMesh adapts to your integration needs.
          </p>
        </div>

        <div className="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-8 lg:max-w-none lg:grid-cols-3">
          {useCases.map((useCase) => (
            <div key={useCase.title} className="card hover:shadow-lg transition-all">
              <div className="text-4xl mb-4">{useCase.icon}</div>
              <h3 className="text-lg font-semibold text-gray-900 mb-3">
                {useCase.title}
              </h3>
              <p className="text-gray-600 mb-4">{useCase.description}</p>
              <div className="flex flex-wrap gap-2">
                {useCase.tech.map((tech) => (
                  <span
                    key={tech}
                    className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-dictamesh-blue-100 text-dictamesh-blue-800"
                  >
                    {tech}
                  </span>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* CTA */}
        <div className="mt-16 text-center">
          <a href="/solutions" className="btn-secondary">
            Explore All Solutions
            <svg className="ml-2 h-5 w-5 inline" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
            </svg>
          </a>
        </div>
      </div>
    </div>
  );
}
