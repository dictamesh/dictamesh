// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import { Hero } from "~/components/marketing/Hero";
import { Features } from "~/components/marketing/Features";
import { Architecture } from "~/components/marketing/Architecture";
import { UseCases } from "~/components/marketing/UseCases";
import { Testimonials } from "~/components/marketing/Testimonials";
import { CTASection } from "~/components/marketing/CTASection";
import { Header } from "~/components/layout/Header";
import { Footer } from "~/components/layout/Footer";

export const meta: MetaFunction = () => {
  return [
    {
      title: "DictaMesh | Enterprise-Grade Reference Architecture: Integration of Federated Authority Sources with Event-Driven Coordination",
    },
    {
      name: "description",
      content:
        "Build federated data integrations with confidence. Open-source framework providing event-driven architecture, GraphQL federation, and enterprise governance for your data mesh.",
    },
    {
      property: "og:title",
      content: "DictaMesh | Enterprise-Grade Reference Architecture: Integration of Federated Authority Sources with Event-Driven Coordination",
    },
    {
      property: "og:description",
      content:
        "Production-ready framework for building data mesh adapters. Integrate any data source with built-in observability, governance, and resilience patterns.",
    },
    {
      property: "og:image",
      content: "https://dictamesh.com/og/home.png",
    },
    {
      property: "og:type",
      content: "website",
    },
    {
      name: "twitter:card",
      content: "summary_large_image",
    },
    {
      name: "twitter:title",
      content: "DictaMesh | Enterprise Data Mesh Framework",
    },
    {
      name: "twitter:description",
      content:
        "Open-source framework for building production-ready data mesh adapters with event-driven architecture and federated GraphQL.",
    },
    {
      name: "twitter:image",
      content: "https://dictamesh.com/og/home.png",
    },
  ];
};

export default function Index() {
  return (
    <div className="min-h-screen bg-white">
      <Header />
      <main>
        <Hero />
        <Features />
        <Architecture />
        <UseCases />
        <Testimonials />
        <CTASection />
      </main>
      <Footer />
    </div>
  );
}
