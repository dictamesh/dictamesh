// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import { Link } from "@remix-run/react";

export function CTASection() {
  return (
    <div className="bg-dictamesh-blue-900 section-padding">
      <div className="container-custom">
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Ready to Build Your Data Mesh?
          </h2>
          <p className="mx-auto mt-6 max-w-xl text-lg leading-8 text-gray-300">
            Start with our open-source framework or get enterprise support with
            managed hosting, SLAs, and dedicated assistance.
          </p>
          <div className="mt-10 flex items-center justify-center gap-x-6">
            <Link
              to="/get-started"
              className="rounded-lg bg-white px-6 py-3 text-base font-semibold text-dictamesh-blue-900 shadow-sm hover:bg-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white transition-colors"
            >
              Get Started Free
            </Link>
            <Link
              to="/contact"
              className="rounded-lg border-2 border-white px-6 py-3 text-base font-semibold text-white hover:bg-white/10 transition-colors"
            >
              Talk to Sales
            </Link>
          </div>

          {/* Features grid */}
          <div className="mt-16 grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div className="flex flex-col items-center">
              <svg
                className="h-8 w-8 text-dictamesh-teal-400 mb-2"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div className="text-white font-semibold">Open Source</div>
              <div className="text-sm text-gray-400 mt-1">AGPL-3.0 License</div>
            </div>
            <div className="flex flex-col items-center">
              <svg
                className="h-8 w-8 text-dictamesh-teal-400 mb-2"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div className="text-white font-semibold">Production Ready</div>
              <div className="text-sm text-gray-400 mt-1">Enterprise Tested</div>
            </div>
            <div className="flex flex-col items-center">
              <svg
                className="h-8 w-8 text-dictamesh-teal-400 mb-2"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div className="text-white font-semibold">Full Support</div>
              <div className="text-sm text-gray-400 mt-1">Available 24/7</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
