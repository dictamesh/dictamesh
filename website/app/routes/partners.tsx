// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { Header } from "~/components/layout/Header";
import { Footer } from "~/components/layout/Footer";

export const meta: MetaFunction = () => {
  return [
    {
      title: "Partner Programs | DictaMesh - Grow Together",
    },
    {
      name: "description",
      content:
        "Join the DictaMesh partner ecosystem. Affiliate, reseller, integration partner, and white-label programs available.",
    },
  ];
};

const programs = [
  {
    title: "Affiliate Program",
    description:
      "Earn recurring commissions by referring customers to DictaMesh managed hosting.",
    icon: "💰",
    benefits: [
      "20% recurring commission on all referrals",
      "90-day cookie duration",
      "Professional marketing materials",
      "Dedicated affiliate dashboard",
      "Monthly payouts via PayPal/Stripe",
      "No approval required to join",
    ],
    ideal: "Content creators, bloggers, tech influencers",
    cta: "Join Affiliate Program",
    link: "/partners/affiliates/apply",
  },
  {
    title: "Reseller Program",
    description:
      "Sell DictaMesh managed hosting under your brand with white-label options.",
    icon: "🤝",
    benefits: [
      "Up to 30% partner discount",
      "White-label hosting options",
      "Co-branded marketing materials",
      "Dedicated partner portal",
      "Pre-sales technical support",
      "Joint go-to-market strategies",
      "Quarterly business reviews",
    ],
    ideal: "System integrators, consulting firms, agencies",
    cta: "Become a Reseller",
    link: "/partners/resellers/apply",
  },
  {
    title: "Integration Partners",
    description:
      "Build and certify connectors/adapters for DictaMesh framework.",
    icon: "🔌",
    benefits: [
      "Listed in official marketplace",
      "Technical enablement & training",
      "Co-marketing opportunities",
      "Early access to new features",
      "Partner certification program",
      "Revenue share on paid integrations",
      "Joint customer success stories",
    ],
    ideal: "Software vendors, platform providers, SaaS companies",
    cta: "Build Integration",
    link: "/partners/integrations/apply",
  },
  {
    title: "White Label",
    description:
      "Rebrand and resell DictaMesh as your own data mesh solution.",
    icon: "🏷️",
    benefits: [
      "Full source code access (AGPL-3.0)",
      "Remove all DictaMesh branding",
      "Custom domain & SSL",
      "Your logo and color scheme",
      "Dedicated infrastructure",
      "Priority support & updates",
      "Custom feature development available",
    ],
    ideal: "Enterprises, cloud providers, consulting firms",
    cta: "Discuss White Label",
    link: "/partners/white-label/apply",
  },
];

const partnerBenefits = [
  {
    title: "Marketing Support",
    description: "Access to co-branded materials, case studies, and joint campaigns.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" d="M10.34 15.84c-.688-.06-1.386-.09-2.09-.09H7.5a4.5 4.5 0 110-9h.75c.704 0 1.402-.03 2.09-.09m0 9.18c.253.962.584 1.892.985 2.783.247.55.06 1.21-.463 1.511l-.657.38c-.551.318-1.26.117-1.527-.461a20.845 20.845 0 01-1.44-4.282m3.102.069a18.03 18.03 0 01-.59-4.59c0-1.586.205-3.124.59-4.59m0 9.18a23.848 23.848 0 018.835 2.535M10.34 6.66a23.847 23.847 0 008.835-2.535m0 0A23.74 23.74 0 0018.795 3m.38 1.125a23.91 23.91 0 011.014 5.395m-1.014 8.855c-.118.38-.245.754-.38 1.125m.38-1.125a23.91 23.91 0 001.014-5.395m0-3.46c.495.413.811 1.035.811 1.73 0 .695-.316 1.317-.811 1.73m0-3.46a24.347 24.347 0 010 3.46" />
      </svg>
    ),
  },
  {
    title: "Technical Enablement",
    description: "Training, documentation, and dedicated technical support for partners.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" d="M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.627 48.627 0 0112 20.904a48.627 48.627 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.57 50.57 0 00-2.658-.813A59.905 59.905 0 0112 3.493a59.902 59.902 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.697 50.697 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443m-7.007 11.55A5.981 5.981 0 006.75 15.75v-1.5" />
      </svg>
    ),
  },
  {
    title: "Revenue Opportunities",
    description: "Attractive margins, recurring revenue, and growth incentives.",
    icon: (
      <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
      </svg>
    ),
  },
];

export default function Partners() {
  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main>
        {/* Hero */}
        <div className="bg-gradient-to-b from-dictamesh-blue-50 to-white section-padding">
          <div className="container-custom">
            <div className="mx-auto max-w-3xl text-center">
              <h1 className="text-4xl font-bold tracking-tight text-dictamesh-blue-900 sm:text-5xl">
                Grow Your Business with DictaMesh
              </h1>
              <p className="mt-6 text-lg leading-8 text-gray-600">
                Join our partner ecosystem and help organizations build better data
                architectures. Choose the program that fits your business model.
              </p>
            </div>
          </div>
        </div>

        {/* Partner Programs */}
        <div className="section-padding">
          <div className="container-custom">
            <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
              {programs.map((program) => (
                <div key={program.title} className="card hover:shadow-xl transition-all">
                  <div className="text-5xl mb-4">{program.icon}</div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-3">
                    {program.title}
                  </h3>
                  <p className="text-gray-600 mb-6">{program.description}</p>

                  <div className="mb-6">
                    <h4 className="font-semibold text-sm text-gray-900 mb-3">
                      Program Benefits:
                    </h4>
                    <ul className="space-y-2">
                      {program.benefits.map((benefit) => (
                        <li key={benefit} className="flex items-start">
                          <svg
                            className="h-5 w-5 flex-shrink-0 text-dictamesh-teal-600 mt-0.5"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={2}
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                          </svg>
                          <span className="ml-3 text-sm text-gray-700">{benefit}</span>
                        </li>
                      ))}
                    </ul>
                  </div>

                  <div className="border-t border-gray-200 pt-4 mb-6">
                    <div className="text-sm text-gray-600">
                      <span className="font-medium">Ideal for: </span>
                      {program.ideal}
                    </div>
                  </div>

                  <Link to={program.link} className="btn-primary w-full text-center">
                    {program.cta}
                    <svg
                      className="ml-2 h-5 w-5 inline"
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
              ))}
            </div>
          </div>
        </div>

        {/* Partner Benefits */}
        <div className="bg-dictamesh-blue-50 section-padding">
          <div className="container-custom">
            <div className="mx-auto max-w-2xl text-center mb-16">
              <h2 className="text-3xl font-bold tracking-tight text-dictamesh-blue-900">
                Why Partner with DictaMesh?
              </h2>
              <p className="mt-4 text-lg text-gray-600">
                We're committed to your success with comprehensive support and resources.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
              {partnerBenefits.map((benefit) => (
                <div key={benefit.title} className="card bg-white">
                  <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-dictamesh-blue-600 text-white mb-4">
                    {benefit.icon}
                  </div>
                  <h4 className="text-lg font-semibold text-gray-900 mb-2">
                    {benefit.title}
                  </h4>
                  <p className="text-gray-600">{benefit.description}</p>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Success Stories */}
        <div className="section-padding">
          <div className="container-custom">
            <div className="mx-auto max-w-2xl text-center mb-16">
              <h2 className="text-3xl font-bold tracking-tight text-gray-900">
                Partner Success Stories
              </h2>
            </div>

            <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
              <div className="card bg-gray-50">
                <div className="text-4xl mb-4">📈</div>
                <blockquote className="text-gray-700 mb-4">
                  "Joining as a reseller increased our revenue by 40% in the first year.
                  DictaMesh's technology and support made the difference."
                </blockquote>
                <div className="border-t border-gray-200 pt-4">
                  <div className="font-semibold text-gray-900">Integration Solutions Inc.</div>
                  <div className="text-sm text-gray-600">Reseller Partner</div>
                </div>
              </div>

              <div className="card bg-gray-50">
                <div className="text-4xl mb-4">🎯</div>
                <blockquote className="text-gray-700 mb-4">
                  "The affiliate program is straightforward and pays well. Great passive
                  income from content I was already creating."
                </blockquote>
                <div className="border-t border-gray-200 pt-4">
                  <div className="font-semibold text-gray-900">Tech Blog Pro</div>
                  <div className="text-sm text-gray-600">Affiliate Partner</div>
                </div>
              </div>

              <div className="card bg-gray-50">
                <div className="text-4xl mb-4">🚀</div>
                <blockquote className="text-gray-700 mb-4">
                  "Building our connector opened doors to enterprise customers. The
                  marketplace visibility is invaluable."
                </blockquote>
                <div className="border-t border-gray-200 pt-4">
                  <div className="font-semibold text-gray-900">DataConnect Platform</div>
                  <div className="text-sm text-gray-600">Integration Partner</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* CTA */}
        <div className="bg-dictamesh-blue-900 section-padding">
          <div className="container-custom">
            <div className="mx-auto max-w-3xl text-center">
              <h2 className="text-3xl font-bold text-white">
                Ready to Partner with Us?
              </h2>
              <p className="mt-6 text-lg text-gray-300">
                Choose the program that fits your business, or contact us to discuss
                a custom partnership.
              </p>
              <div className="mt-10 flex flex-col sm:flex-row gap-4 justify-center">
                <Link
                  to="/partners/apply"
                  className="rounded-lg bg-white px-6 py-3 text-base font-semibold text-dictamesh-blue-900 hover:bg-gray-100 transition-colors"
                >
                  Apply Now
                </Link>
                <Link
                  to="/contact"
                  className="rounded-lg border-2 border-white px-6 py-3 text-base font-semibold text-white hover:bg-white/10 transition-colors"
                >
                  Contact Partnership Team
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
