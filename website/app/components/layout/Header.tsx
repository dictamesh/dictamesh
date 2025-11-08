// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import { Link } from "@remix-run/react";
import { useState } from "react";

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <nav className="container-custom" aria-label="Global">
        <div className="flex items-center justify-between py-4">
          {/* Logo */}
          <div className="flex lg:flex-1">
            <Link to="/" className="-m-1.5 p-1.5">
              <span className="sr-only">DictaMesh</span>
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-gradient-to-br from-dictamesh-blue-500 to-dictamesh-purple-600 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-lg">DM</span>
                </div>
                <span className="text-xl font-bold text-dictamesh-blue-900">
                  DictaMesh
                </span>
              </div>
            </Link>
          </div>

          {/* Mobile menu button */}
          <div className="flex lg:hidden">
            <button
              type="button"
              className="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              <span className="sr-only">Open main menu</span>
              <svg
                className="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                />
              </svg>
            </button>
          </div>

          {/* Desktop navigation */}
          <div className="hidden lg:flex lg:gap-x-8">
            <Link
              to="/features"
              className="text-sm font-semibold leading-6 text-gray-900 hover:text-dictamesh-blue-600 transition-colors"
            >
              Features
            </Link>
            <Link
              to="/docs"
              className="text-sm font-semibold leading-6 text-gray-900 hover:text-dictamesh-blue-600 transition-colors"
            >
              Documentation
            </Link>
            <Link
              to="/pricing"
              className="text-sm font-semibold leading-6 text-gray-900 hover:text-dictamesh-blue-600 transition-colors"
            >
              Pricing
            </Link>
            <Link
              to="/partners"
              className="text-sm font-semibold leading-6 text-gray-900 hover:text-dictamesh-blue-600 transition-colors"
            >
              Partners
            </Link>
            <Link
              to="/blog"
              className="text-sm font-semibold leading-6 text-gray-900 hover:text-dictamesh-blue-600 transition-colors"
            >
              Blog
            </Link>
          </div>

          {/* CTA buttons */}
          <div className="hidden lg:flex lg:flex-1 lg:justify-end lg:gap-x-4">
            <a
              href="https://github.com/click2-run/dictamesh"
              target="_blank"
              rel="noopener noreferrer"
              className="btn-ghost text-sm"
            >
              GitHub
            </a>
            <Link to="/get-started" className="btn-primary text-sm">
              Get Started
            </Link>
          </div>
        </div>

        {/* Mobile menu */}
        {mobileMenuOpen && (
          <div className="lg:hidden border-t border-gray-200 py-4">
            <div className="space-y-2 py-6">
              <Link
                to="/features"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                Features
              </Link>
              <Link
                to="/docs"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                Documentation
              </Link>
              <Link
                to="/pricing"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                Pricing
              </Link>
              <Link
                to="/partners"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                Partners
              </Link>
              <Link
                to="/blog"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                Blog
              </Link>
              <div className="border-t border-gray-200 my-4"></div>
              <a
                href="https://github.com/click2-run/dictamesh"
                target="_blank"
                rel="noopener noreferrer"
                className="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
              >
                GitHub
              </a>
              <Link
                to="/get-started"
                className="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-white bg-dictamesh-blue-500 hover:bg-dictamesh-blue-600 transition-colors"
              >
                Get Started
              </Link>
            </div>
          </div>
        )}
      </nav>
    </header>
  );
}
