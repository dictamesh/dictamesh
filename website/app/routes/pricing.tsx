// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { Header } from "~/components/layout/Header";
import { Footer } from "~/components/layout/Footer";

export const meta: MetaFunction = () => {
  return [
    {
      title: "Pricing | DictaMesh - Managed Hosting & Enterprise Support",
    },
    {
      name: "description",
      content:
        "Choose the right plan for your data mesh. Self-hosted open source, managed cloud hosting, or enterprise support with SLAs and dedicated assistance.",
    },
  ];
};

const tiers = [
  {
    name: "Open Source",
    id: "tier-opensource",
    price: "$0",
    description: "Perfect for individuals and small teams getting started.",
    features: [
      "Full framework source code",
      "Community support via GitHub",
      "Complete documentation",
      "Reference implementations",
      "Docker Compose templates",
      "MIT license adapters library",
      "Unlimited adapters",
      "Self-hosted deployment",
    ],
    cta: "Get Started",
    ctaLink: "/get-started",
    highlighted: false,
  },
  {
    name: "Professional",
    id: "tier-professional",
    price: "$499",
    priceDetail: "/month",
    description: "Managed cloud hosting with professional support.",
    features: [
      "Everything in Open Source",
      "Managed Kafka cluster",
      "Managed PostgreSQL database",
      "Managed Redis cache",
      "99.9% uptime SLA",
      "Up to 1M events/month",
      "Up to 10M API calls/month",
      "Email support (24h response)",
      "Monthly security updates",
      "Automated backups",
      "Monitoring dashboards",
      "SSL certificates included",
    ],
    cta: "Start Trial",
    ctaLink: "/trial",
    highlighted: true,
  },
  {
    name: "Enterprise",
    id: "tier-enterprise",
    price: "Custom",
    description: "For organizations requiring dedicated support and SLAs.",
    features: [
      "Everything in Professional",
      "99.99% uptime SLA",
      "Unlimited events & API calls",
      "Dedicated support engineer",
      "24/7 phone & Slack support",
      "1-hour critical response",
      "Custom adapter development",
      "Architecture consulting",
      "On-site training available",
      "Private cloud deployment",
      "Multi-region setup",
      "Compliance certifications",
      "Custom SLAs available",
    ],
    cta: "Contact Sales",
    ctaLink: "/contact",
    highlighted: false,
  },
];

const addons = [
  {
    name: "Additional Events",
    price: "$50",
    unit: "per 1M events/month",
  },
  {
    name: "Additional API Calls",
    price: "$25",
    unit: "per 10M calls/month",
  },
  {
    name: "Adapter Development",
    price: "$5,000",
    unit: "per custom adapter",
  },
  {
    name: "Professional Services",
    price: "$200",
    unit: "per hour",
  },
];

export default function Pricing() {
  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main className="section-padding">
        <div className="container-custom">
          {/* Header */}
          <div className="mx-auto max-w-2xl text-center">
            <h1 className="text-4xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-5xl">
              Choose Your Plan
            </h1>
            <p className="mt-6 text-lg leading-8 text-gray-600">
              Start with open source and scale to managed hosting when you're ready.
              All plans include the full framework with production-ready features.
            </p>
          </div>

          {/* Pricing Cards */}
          <div className="mx-auto mt-16 grid max-w-lg grid-cols-1 gap-8 lg:max-w-none lg:grid-cols-3">
            {tiers.map((tier) => (
              <div
                key={tier.id}
                className={`card ${
                  tier.highlighted
                    ? "ring-2 ring-dictamesh-blue-500 shadow-xl scale-105"
                    : "shadow-md"
                }`}
              >
                {tier.highlighted && (
                  <div className="absolute -top-5 left-0 right-0 mx-auto w-fit">
                    <span className="inline-flex rounded-full bg-dictamesh-blue-500 px-4 py-1 text-sm font-semibold text-white">
                      Most Popular
                    </span>
                  </div>
                )}

                <div className="p-8">
                  <h3 className="text-2xl font-semibold text-gray-900">
                    {tier.name}
                  </h3>
                  <p className="mt-4 text-sm text-gray-600">{tier.description}</p>
                  <p className="mt-8">
                    <span className="text-4xl font-bold text-dictamesh-blue-900">
                      {tier.price}
                    </span>
                    {tier.priceDetail && (
                      <span className="text-sm text-gray-600">{tier.priceDetail}</span>
                    )}
                  </p>

                  <Link
                    to={tier.ctaLink}
                    className={`mt-8 block w-full text-center ${
                      tier.highlighted ? "btn-primary" : "btn-secondary"
                    }`}
                  >
                    {tier.cta}
                  </Link>

                  <ul className="mt-8 space-y-3">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex items-start">
                        <svg
                          className="h-6 w-6 flex-shrink-0 text-dictamesh-teal-600"
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
                        <span className="ml-3 text-sm text-gray-700">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            ))}
          </div>

          {/* Add-ons Section */}
          <div className="mt-24 border-t border-gray-200 pt-16">
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="text-3xl font-bold tracking-tight text-gray-900">
                Add-ons & Services
              </h2>
              <p className="mt-4 text-lg text-gray-600">
                Extend your plan with additional capacity and professional services.
              </p>
            </div>

            <div className="mx-auto mt-12 grid max-w-4xl grid-cols-1 gap-6 sm:grid-cols-2">
              {addons.map((addon) => (
                <div key={addon.name} className="card bg-gray-50">
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-semibold text-gray-900">{addon.name}</h4>
                      <p className="text-sm text-gray-600 mt-1">{addon.unit}</p>
                    </div>
                    <div className="text-right">
                      <div className="font-bold text-dictamesh-blue-900">{addon.price}</div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* FAQ Section */}
          <div className="mt-24 border-t border-gray-200 pt-16">
            <div className="mx-auto max-w-2xl text-center mb-12">
              <h2 className="text-3xl font-bold tracking-tight text-gray-900">
                Frequently Asked Questions
              </h2>
            </div>

            <div className="mx-auto max-w-3xl space-y-8">
              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  Can I start with open source and upgrade later?
                </h4>
                <p className="text-gray-600">
                  Absolutely! Start self-hosting with the open-source framework. When you're
                  ready for managed hosting, we'll migrate your data with zero downtime.
                </p>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  What's included in the SLA?
                </h4>
                <p className="text-gray-600">
                  Professional tier includes 99.9% uptime (43 minutes downtime/month max).
                  Enterprise tier offers 99.99% (4 minutes/month max) with financial credits
                  for any breaches.
                </p>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  Can I deploy on my own infrastructure?
                </h4>
                <p className="text-gray-600">
                  Yes! The open-source framework can be deployed anywhere: your data center,
                  AWS, Azure, GCP, or any Kubernetes cluster. Enterprise plans include
                  deployment assistance.
                </p>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  What if I exceed my event/API limits?
                </h4>
                <p className="text-gray-600">
                  We'll notify you at 80% usage. You can add capacity via add-ons or upgrade
                  to Enterprise for unlimited usage. No surprise charges or service interruptions.
                </p>
              </div>

              <div>
                <h4 className="font-semibold text-gray-900 mb-2">
                  Do you offer custom development?
                </h4>
                <p className="text-gray-600">
                  Yes! Enterprise customers can request custom adapter development, specialized
                  connectors, or framework extensions. Contact our sales team for pricing.
                </p>
              </div>
            </div>
          </div>

          {/* CTA */}
          <div className="mt-24 rounded-2xl bg-dictamesh-blue-900 p-12 text-center">
            <h3 className="text-3xl font-bold text-white">
              Not sure which plan is right?
            </h3>
            <p className="mt-4 text-lg text-gray-300">
              Schedule a call with our team to discuss your requirements.
            </p>
            <div className="mt-8">
              <Link
                to="/contact"
                className="inline-flex items-center rounded-lg bg-white px-6 py-3 text-base font-semibold text-dictamesh-blue-900 hover:bg-gray-100 transition-colors"
              >
                Schedule a Demo
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
            </div>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
