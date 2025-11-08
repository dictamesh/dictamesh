// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  icon: string;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Framework, Not Implementation',
    icon: '🏗️',
    description: (
      <>
        DictaMesh is a comprehensive framework that provides the infrastructure
        for building data mesh adapters. You build adapters for your specific
        data sources - we handle the rest.
      </>
    ),
  },
  {
    title: 'Event-Driven Architecture',
    icon: '⚡',
    description: (
      <>
        Built on Kafka with structured event schemas, automatic publishing,
        and consumer patterns. Real-time data synchronization across your
        entire data mesh.
      </>
    ),
  },
  {
    title: 'Federated GraphQL Gateway',
    icon: '🔗',
    description: (
      <>
        Automatic API composition from your adapters using Apollo Federation.
        Query data across multiple sources with a single GraphQL query.
      </>
    ),
  },
  {
    title: 'Metadata Catalog',
    icon: '📊',
    description: (
      <>
        Complete entity registry, relationships, and lineage tracking.
        Know where your data lives and how it's connected across systems.
      </>
    ),
  },
  {
    title: 'Built-in Observability',
    icon: '👁️',
    description: (
      <>
        Distributed tracing with OpenTelemetry, Prometheus metrics,
        and structured logging. Full visibility into your data mesh.
      </>
    ),
  },
  {
    title: 'Production-Ready Patterns',
    icon: '🛡️',
    description: (
      <>
        Circuit breakers, retry logic, rate limiting, caching, and
        graceful degradation - all built in and battle-tested.
      </>
    ),
  },
];

function Feature({title, icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <div className={styles.featureIcon}>{icon}</div>
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
