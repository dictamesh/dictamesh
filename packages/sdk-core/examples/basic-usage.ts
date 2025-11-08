// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

/**
 * Basic usage example for DictaMesh SDK
 */

import {
  DictaMeshClient,
  DictaMeshAdapter,
  HTTPConnector,
  MemoryCache,
} from '../src';

async function main() {
  console.log('=== DictaMesh SDK Basic Usage Example ===\n');

  // 1. Create HTTP connector
  console.log('1. Creating HTTP connector...');
  const connector = new HTTPConnector({
    endpoint: 'https://api.example.com/dictamesh',
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // 2. Create and initialize adapter
  console.log('2. Creating DictaMesh adapter...');
  const adapter = new DictaMeshAdapter({
    name: 'main',
    version: '1.0.0',
    endpoint: 'https://api.example.com',
    connector,
  });

  await adapter.initialize({
    name: 'main',
    version: '1.0.0',
    connector,
  });

  // 3. Create client
  console.log('3. Creating client...');
  const client = new DictaMeshClient({
    endpoint: 'https://api.example.com',
    defaultAdapter: 'main',
    cache: {
      enabled: true,
      type: 'memory',
      ttl: 300000, // 5 minutes
    },
    auth: {
      type: 'bearer',
      token: 'your-api-token-here',
    },
    timeout: 30000,
    tracing: true,
  });

  // 4. Register adapter
  console.log('4. Registering adapter...');
  client.registerAdapter('main', adapter);
  client.setDefaultAdapter('main');

  // 5. Set cache
  console.log('5. Setting up cache...');
  const cache = new MemoryCache(300000); // 5 minutes TTL
  client.setCache(cache);

  // 6. Connect
  console.log('6. Connecting...\n');
  await client.connect();

  console.log('=== Executing Operations ===\n');

  // Example 1: Get single entity
  console.log('Example 1: Get single customer');
  try {
    const customer = await client.query({
      type: 'get',
      entity: 'customer',
      params: { id: '123' },
      options: {
        select: ['id', 'name', 'email', 'status'],
        cache: {
          enabled: true,
          strategy: 'cache-first',
        },
      },
    });
    console.log('Result:', customer.data);
    console.log('Meta:', customer.meta);
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 2: List entities with filtering
  console.log('Example 2: List active customers');
  try {
    const customers = await client.query({
      type: 'list',
      entity: 'customer',
      options: {
        where: {
          status: 'active',
          createdAt: { $gte: '2025-01-01' },
        },
        orderBy: [
          { field: 'createdAt', direction: 'desc' },
          { field: 'name', direction: 'asc' },
        ],
        limit: 10,
        offset: 0,
        select: ['id', 'name', 'email', 'createdAt'],
      },
    });
    console.log('Count:', customers.meta?.count);
    console.log('Total:', customers.meta?.total);
    console.log('Has more:', customers.meta?.hasMore);
    console.log('First customer:', customers.data?.[0]);
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 3: Search
  console.log('Example 3: Search products');
  try {
    const results = await client.query({
      type: 'search',
      entity: 'product',
      params: {
        query: 'laptop',
        fields: ['name', 'description', 'tags'],
      },
      options: {
        limit: 5,
        where: {
          status: 'available',
          price: { $gte: 500, $lte: 2000 },
        },
      },
    });
    console.log('Results:', results.data);
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 4: Create entity
  console.log('Example 4: Create new customer');
  try {
    const newCustomer = await client.mutate({
      type: 'create',
      entity: 'customer',
      data: {
        name: 'Jane Doe',
        email: 'jane@example.com',
        status: 'active',
        metadata: {
          source: 'api',
          campaign: 'spring-2025',
        },
      },
    });
    console.log('Created:', newCustomer.data);
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 5: Update entity
  console.log('Example 5: Update customer');
  try {
    const updated = await client.mutate({
      type: 'update',
      entity: 'customer',
      params: { id: '123' },
      data: {
        status: 'premium',
        lastUpdated: new Date().toISOString(),
      },
    });
    console.log('Updated:', updated.data);
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 6: Batch operations
  console.log('Example 6: Batch operations');
  try {
    const batchResults = await client.batch([
      {
        type: 'get',
        entity: 'customer',
        params: { id: '123' },
      },
      {
        type: 'get',
        entity: 'customer',
        params: { id: '456' },
      },
      {
        type: 'list',
        entity: 'product',
        options: { limit: 5 },
      },
    ]);
    console.log('Total operations:', batchResults.meta?.total);
    console.log('Successful:', batchResults.meta?.successful);
    console.log('Failed:', batchResults.meta?.failed);
    console.log('Duration:', batchResults.meta?.took, 'ms');
  } catch (error) {
    console.error('Error:', error);
  }
  console.log('');

  // Example 7: Cache statistics
  console.log('Example 7: Cache statistics');
  const cacheInstance = client.getCache();
  if (cacheInstance) {
    const stats = cacheInstance.getStats();
    console.log('Cache hits:', stats.hits);
    console.log('Cache misses:', stats.misses);
    console.log('Cache hit rate:', (stats.hitRate * 100).toFixed(2) + '%');
    console.log('Cache size:', stats.size);
  }
  console.log('');

  // Example 8: Clear cache
  console.log('Example 8: Clear cache');
  await client.clearCache({ entity: 'customer' });
  console.log('Cleared customer cache');
  console.log('');

  // Example 9: Client status
  console.log('Example 9: Client status');
  const status = client.getStatus();
  console.log('Connected:', status.connected);
  console.log('Default adapter:', status.adapter);
  console.log('Connection state:', status.connectionState);
  console.log('Request count:', status.metrics?.requestCount);
  console.log('Error count:', status.metrics?.errorCount);
  console.log('');

  // Cleanup
  console.log('=== Cleanup ===');
  await client.disconnect();
  await adapter.dispose();
  console.log('Disconnected and cleaned up');
}

// Run example
if (require.main === module) {
  main()
    .then(() => {
      console.log('\nExample completed successfully');
      process.exit(0);
    })
    .catch(error => {
      console.error('\nExample failed:', error);
      process.exit(1);
    });
}

export default main;
