<!--
SPDX-License-Identifier: AGPL-3.0-or-later
Copyright (C) 2025 Controle Digital Ltda
-->

---
sidebar_position: 1
---

# Contributing to DictaMesh

Thank you for your interest in contributing to DictaMesh! We welcome contributions from the community and are grateful for your support in making this project better.

This guide will help you understand how to contribute effectively to the project, whether you're fixing bugs, adding features, improving documentation, or helping in other ways.

## Ways to Contribute

There are many ways you can contribute to DictaMesh:

### Code Contributions
- Fix bugs and issues
- Implement new features
- Improve performance
- Add new connectors or adapters
- Enhance test coverage

### Documentation
- Fix typos and improve clarity
- Add examples and tutorials
- Translate documentation
- Create video tutorials or blog posts

### Community Support
- Answer questions in GitHub Discussions
- Help review pull requests
- Report bugs and suggest features
- Share your DictaMesh use cases

### Testing
- Test new releases
- Report bugs with detailed reproduction steps
- Improve test coverage
- Add integration tests

## Getting Started

### 1. Fork and Clone

Start by forking the repository and cloning it locally:

```bash
# Fork the repository on GitHub first, then:
git clone https://github.com/YOUR-USERNAME/dictamesh.git
cd dictamesh

# Add upstream remote
git remote add upstream https://github.com/Click2-Run/dictamesh.git
```

### 2. Set Up Development Environment

Follow our [Development Setup Guide](./development-setup.md) to configure your local environment.

### 3. Create a Branch

Create a new branch for your changes:

```bash
# Update your local main branch
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/my-new-feature

# Or for bug fixes
git checkout -b fix/issue-123
```

### Branch Naming Conventions

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring
- `test/description` - Test improvements
- `chore/description` - Maintenance tasks

## Development Workflow

### 1. Make Your Changes

- Write clean, maintainable code
- Follow the project's code style (see below)
- Add tests for new functionality
- Update documentation as needed
- Follow the copyright header requirements (see AGENT.md)

### 2. Test Your Changes

Before submitting, ensure all tests pass:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linters
make lint

# Run integration tests
make test-integration
```

### 3. Commit Your Changes

Write clear, descriptive commit messages:

```bash
# Stage your changes
git add .

# Commit with a descriptive message
git commit -m "feat: add support for MongoDB connector

- Implement MongoDB connection interface
- Add authentication handling
- Include comprehensive tests
- Update documentation

Closes #123"
```

### Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

**Format:**
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

**Examples:**
```
feat(connectors): add MongoDB support

Implement MongoDB connector with connection pooling,
authentication, and error handling.

Closes #123
```

```
fix(adapters): handle null values in GraphQL response

Previously, null values in API responses caused panics.
Now properly handles null values with default empty values.

Fixes #456
```

### 4. Push and Create Pull Request

```bash
# Push your branch
git push origin feature/my-new-feature

# Create a pull request on GitHub
```

## Pull Request Process

### Before Submitting

- [ ] All tests pass locally
- [ ] Code follows project style guidelines
- [ ] Documentation is updated
- [ ] Copyright headers are added to new files
- [ ] Commit messages follow conventions
- [ ] PR description clearly explains the changes

### PR Title Format

Use the same format as commit messages:

```
feat: add MongoDB connector support
fix: resolve GraphQL null pointer exception
docs: improve adapter building guide
```

### PR Description Template

When creating a PR, include:

```markdown
## Description
Brief description of what this PR does.

## Motivation and Context
Why is this change needed? What problem does it solve?

## Type of Change
- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to change)
- [ ] Documentation update

## How Has This Been Tested?
Describe the tests you ran and how to reproduce them.

## Checklist
- [ ] My code follows the project's code style
- [ ] I have added tests covering my changes
- [ ] All new and existing tests pass
- [ ] I have updated the documentation
- [ ] I have added copyright headers to new files
- [ ] My changes generate no new warnings
```

### Review Process

1. **Automated Checks**: CI/CD runs tests and linters automatically
2. **Code Review**: Maintainers review your code for quality and correctness
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, a maintainer will merge your PR

**Review Timeline:**
- Initial review: Within 3-5 business days
- Follow-up reviews: Within 1-2 business days
- Merge: After approval and passing all checks

## Code Style Guidelines

### Go Code Style

We follow standard Go conventions with some additional guidelines:

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package connector

import (
    "context"
    "fmt"
)

// MongoDBConnector implements the Connector interface for MongoDB.
// All exported types and functions must have documentation comments.
type MongoDBConnector struct {
    client *mongo.Client
    config *Config
}

// Connect establishes a connection to MongoDB.
// Documentation should explain what the function does and any important details.
func (c *MongoDBConnector) Connect(ctx context.Context) error {
    // Use meaningful variable names
    // Keep functions focused and small
    // Handle errors explicitly
    client, err := mongo.Connect(ctx, c.config.Options())
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    c.client = client
    return nil
}
```

**Key Points:**
- Use `gofmt` and `goimports` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write clear documentation comments for exported items
- Handle errors explicitly, never ignore them
- Use meaningful variable and function names
- Keep functions small and focused

### Database Naming Conventions

**CRITICAL**: All database objects must use the `dictamesh_` prefix:

```go
// Table names
func (EntityCatalog) TableName() string {
    return "dictamesh_entity_catalog"
}

// Indexes
CREATE INDEX idx_dictamesh_entity_catalog_type ON dictamesh_entity_catalog(entity_type);

// Functions
CREATE FUNCTION dictamesh_update_timestamp() RETURNS trigger AS $$
```

See the AGENT.md file in the repository root for complete database naming requirements.

### Testing Requirements

All code contributions must include tests:

```go
// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package connector

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMongoDBConnector_Connect(t *testing.T) {
    // Arrange
    connector := &MongoDBConnector{
        config: &Config{
            Host: "localhost",
            Port: 27017,
        },
    }

    // Act
    err := connector.Connect(context.Background())

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, connector.client)
}
```

**Testing Guidelines:**
- Write unit tests for all public functions
- Use table-driven tests for multiple test cases
- Mock external dependencies
- Aim for 80%+ code coverage
- Include both positive and negative test cases

### Documentation Style

When writing documentation:

- Use clear, concise language
- Include code examples
- Add copyright headers to documentation files
- Use proper markdown formatting
- Keep line length reasonable (100 characters)
- Include links to related documentation

## Reporting Bugs

### Before Reporting

1. Check if the bug has already been reported
2. Test with the latest version
3. Gather relevant information

### Bug Report Template

Create a new issue with:

```markdown
**Description**
A clear description of the bug.

**To Reproduce**
Steps to reproduce the behavior:
1. Configure '...'
2. Run '...'
3. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- DictaMesh Version: [e.g., 1.0.0]
- Go Version: [e.g., 1.21]
- OS: [e.g., Ubuntu 22.04]
- Database: [e.g., PostgreSQL 16]

**Logs**
```
Include relevant logs here
```

**Additional Context**
Any other information that might be helpful.
```

## Suggesting Features

We welcome feature suggestions! To propose a new feature:

1. **Check Existing Issues**: Ensure it hasn't been suggested before
2. **Create a Feature Request**: Open a new issue with the `enhancement` label
3. **Describe the Feature**: Explain what you want and why
4. **Provide Context**: Include use cases and examples

### Feature Request Template

```markdown
**Problem Statement**
What problem does this feature solve?

**Proposed Solution**
Describe your proposed solution.

**Alternatives Considered**
What other approaches did you consider?

**Use Cases**
How would this feature be used?

**Additional Context**
Any other relevant information.
```

## Community Guidelines

Please read and follow our [Code of Conduct](./code-of-conduct.md). We are committed to providing a welcoming and inclusive environment for all contributors.

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions, ideas, and general discussion
- **Pull Requests**: Code review and collaboration

### Getting Help

If you need help:

1. Check the [documentation](../getting-started/introduction.md)
2. Search [GitHub Discussions](https://github.com/Click2-Run/dictamesh/discussions)
3. Ask in GitHub Discussions
4. Review existing issues

## Recognition

Contributors are recognized in several ways:

- Listed in the project's contributors
- Mentioned in release notes for significant contributions
- Invited to become maintainers for sustained contributions

## License

By contributing to DictaMesh, you agree that your contributions will be licensed under the [GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later)](https://github.com/Click2-Run/dictamesh/blob/main/LICENSE).

Key points about the license:
- Your contributions can be used commercially
- If someone modifies and provides DictaMesh as a network service, they must make the source code available
- All derivative works must use the same license

## Questions?

If you have questions about contributing:

- Read the [Development Setup Guide](./development-setup.md)
- Check [GitHub Discussions](https://github.com/Click2-Run/dictamesh/discussions)
- Review existing pull requests for examples
- Ask in a GitHub Discussion

## Thank You!

Your contributions, no matter how small, make DictaMesh better for everyone. We appreciate your time and effort in helping improve this project.

---

**Next**: [Code of Conduct →](./code-of-conduct.md)
