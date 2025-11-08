// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

const testimonials = [
  {
    quote:
      "DictaMesh gave us the foundation to build our data mesh in weeks, not months. The built-in observability and governance saved us from reinventing the wheel.",
    author: "Engineering Leader",
    role: "VP of Engineering",
    company: "Fortune 500 Retail Company",
  },
  {
    quote:
      "Finally, an open-source framework that doesn't compromise on enterprise features. The event-driven architecture and GraphQL federation work flawlessly together.",
    author: "Platform Architect",
    role: "Principal Architect",
    company: "Global SaaS Provider",
  },
  {
    quote:
      "We integrated 12 legacy systems in 3 months using DictaMesh. The circuit breakers and retry patterns prevented countless production issues.",
    author: "Integration Team Lead",
    role: "Senior Engineering Manager",
    company: "Financial Services Firm",
  },
];

export function Testimonials() {
  return (
    <div className="bg-dictamesh-blue-50 section-padding">
      <div className="container-custom">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-base font-semibold leading-7 text-dictamesh-blue-600">
            Trusted by Teams
          </h2>
          <p className="mt-2 text-3xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-4xl">
            Built for Production
          </p>
        </div>

        <div className="mx-auto grid max-w-2xl grid-cols-1 gap-8 lg:max-w-none lg:grid-cols-3">
          {testimonials.map((testimonial, idx) => (
            <div key={idx} className="card bg-white">
              <div className="flex items-start space-x-1 mb-4">
                {[...Array(5)].map((_, i) => (
                  <svg
                    key={i}
                    className="h-5 w-5 text-dictamesh-amber-500"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                ))}
              </div>
              <blockquote className="text-gray-700 mb-6">
                "{testimonial.quote}"
              </blockquote>
              <div className="border-t border-gray-200 pt-4">
                <div className="font-semibold text-gray-900">
                  {testimonial.author}
                </div>
                <div className="text-sm text-gray-600">{testimonial.role}</div>
                <div className="text-sm text-gray-500">{testimonial.company}</div>
              </div>
            </div>
          ))}
        </div>

        {/* Trust indicators */}
        <div className="mt-16 border-t border-gray-200 pt-16">
          <div className="text-center mb-8">
            <p className="text-sm font-semibold text-gray-600 uppercase tracking-wide">
              Built on Validated Patterns From
            </p>
          </div>
          <div className="grid grid-cols-2 gap-8 md:grid-cols-4 lg:grid-cols-6 opacity-60 items-center justify-items-center">
            <div className="text-center text-sm font-medium text-gray-600">Netflix</div>
            <div className="text-center text-sm font-medium text-gray-600">Uber</div>
            <div className="text-center text-sm font-medium text-gray-600">LinkedIn</div>
            <div className="text-center text-sm font-medium text-gray-600">Airbnb</div>
            <div className="text-center text-sm font-medium text-gray-600">PayPal</div>
            <div className="text-center text-sm font-medium text-gray-600">Microsoft</div>
          </div>
          <p className="text-center text-xs text-gray-500 mt-4">
            Architecture patterns validated at Fortune 500 scale
          </p>
        </div>
      </div>
    </div>
  );
}
