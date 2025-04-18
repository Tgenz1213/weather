# Contributing to This Project

Thank you for your interest in contributing! We welcome pull requests and suggestions that help make this project better. Please review the following guidelines to ensure a smooth contribution process.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Workflow](#workflow)
- [Branching & Commit Messages](#branching--commit-messages)
- [Pull Requests](#pull-requests)
- [Code Quality](#code-quality)
- [Testing](#testing)
- [Documentation](#documentation)
- [Changelog](#changelog)
- [Code of Conduct](#code-of-conduct)
- [Contact](#contact)

---

## Getting Started

- Fork the repository and clone your fork locally.
- Set up your development environment according to the [README](./README.md).
- Make sure your code builds and passes all tests before submitting a PR.

---

## Workflow

- We use **trunk-based development**: all changes are merged into the `main` branch via pull requests.
- Keep your changes focused and minimal; avoid large, unrelated changes in a single PR.
- Sync your branch with the latest `main` before opening a PR.

---

## Branching & Commit Messages

- **Branch naming:** Use the format `<type>/<short-description>`, e.g., `feature/add-auth`, `fix/login-bug`.
- **Commit messages:** Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/): `<type>(<scope>): <description>`

Example: `fix(auth): correct password hashing logic (#123)`

---

## Pull Requests

- Ensure there is an approved issue for your work. Reference it in your PR (e.g., `Closes #123`).
- Complete the PR template, describing your changes and linking related issues.
- Squash your commits before merging.
- Ensure your PR passes all CI checks (tests, linting, etc.).
- Assign reviewers as appropriate.

---

## Code Quality

- Follow the project’s code style and linting rules.
- Write clear, maintainable, and well-documented code.
- Prefer small, focused changes.
- Address all automated review comments and requested changes.

---

## Testing

- All new code must have appropriate test coverage (unit/integration/E2E as relevant).
- Run all tests locally before pushing.
- If you fix a bug, add a test that covers the regression.

---

## Documentation

- Update documentation (README, code comments, API docs) for any new features, changes, or breaking changes.
- Ensure all public interfaces are documented.

---

## Changelog

- The changelog is generated automatically from commit messages on `main` at release time.
- Use clear, conventional commit messages to ensure your changes are included in the changelog.

---

## Code of Conduct

- Be respectful and inclusive in all interactions.
- Review our [Code of Conduct](./CODE_OF_CONDUCT.md) for details.

---

## Contact

For questions, reach out by opening an issue or contacting the maintainers.

---

Thank you for helping make this project better!



