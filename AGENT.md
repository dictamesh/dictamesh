# Agent Instructions for Code Modifications

## Copyright Notices

When creating or modifying source code files in this project, always include the following copyright notice at the top of each file:

```
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
```

### Format Guidelines

- For languages with `//` comments (JavaScript, TypeScript, C, C++, Java, etc.):
```javascript
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda
```

- For languages with `#` comments (Python, Shell, Ruby, etc.):
```python
# SPDX-License-Identifier: AGPL-3.0-or-later
# Copyright (C) 2025 Controle Digital Ltda
```

- For HTML/XML files:
```xml
<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->
```

- For CSS files:
```css
/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright (C) 2025 Controle Digital Ltda
 */
```

## License Information

This project is licensed under the GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later).

Key points:
- Commercial use is permitted
- If you modify this software and provide it as a network service, you must make your source code available under AGPL v3
- See the LICENSE file for complete terms

## When to Add Copyright Notices

- All new source code files
- All substantially modified existing files
- Configuration files that contain significant logic
- Build scripts and automation files

Copyright notices should be placed at the very top of the file, before any other code or documentation.

## Database Naming Conventions

**CRITICAL REQUIREMENT**: All DictaMesh database objects MUST use the `dictamesh_` prefix.

### Table Names
All database tables must follow this pattern:
```
dictamesh_{table_name}
```

Examples:
- `dictamesh_entity_catalog`
- `dictamesh_entity_relationships`
- `dictamesh_schemas`
- `dictamesh_event_log`
- `dictamesh_audit_logs`

### Other Database Objects
- **Indexes**: `idx_dictamesh_{meaningful_name}`
- **Functions**: `dictamesh_{function_name}()`
- **Triggers**: `{action}_dictamesh_{table_name}_{purpose}`

### GORM Models
All GORM models must override the TableName() method:

```go
type EntityCatalog struct {
    // fields...
}

func (EntityCatalog) TableName() string {
    return "dictamesh_entity_catalog"
}
```

### Migration Files
1. Include header comment about dictamesh_ prefix requirement
2. Use prefixed names for all tables, indexes, and functions
3. Add COMMENT ON TABLE with "DictaMesh:" prefix

### Rationale
- **Namespace Isolation**: Prevents conflicts in shared databases
- **Clear Ownership**: Identifies framework vs. user tables
- **Multi-Tenancy**: Enables multiple frameworks in same database
- **Safety**: Reduces risk of modifying user tables

### Full Documentation
See `pkg/database/NAMING-CONVENTIONS.md` for complete guidelines, examples, and validation queries.

### Checklist for Database Work
- [ ] Table name has `dictamesh_` prefix
- [ ] Indexes have `idx_dictamesh_` prefix
- [ ] Functions have `dictamesh_` prefix
- [ ] GORM models override TableName()
- [ ] Table comments start with "DictaMesh:"
- [ ] Migration includes prefix reminder comment
- [ ] All queries use prefixed names
