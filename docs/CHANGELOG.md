# Changelog

All notable changes to the dago-node-planner project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure
- Core planning service with LLM integration
- Task analysis capabilities
- Iterative graph generation with validation
- Schema repository and validation
- REST API with planning and validation endpoints
- Client SDK for Go
- Docker support
- GitHub Actions CI/CD workflows
- Comprehensive documentation
- Example tasks and graphs
- Prompt templates for LLM interaction

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)

## [0.1.0] - TBD

### Added
- Initial release of dago-node-planner
- LLM-powered graph generation from natural language
- Support for Anthropic Claude API
- JSON schema validation for graphs
- Iterative refinement with error feedback
- Task complexity analysis
- REST API endpoints:
  - POST /api/v1/plan - Generate graph from task
  - POST /api/v1/validate - Validate graph JSON
  - GET /health - Health check
  - GET /ready - Readiness check
- Configuration via YAML file and environment variables
- Docker image with multi-stage builds
- Makefile with common development tasks
- Development scripts for local testing
- Comprehensive test coverage
- API documentation
- Usage examples

### Dependencies
- dago-libs v0.2.0
- dago-adapters v0.1.0
- Go 1.25.5
- Anthropic Claude API

## Version History

- **v0.1.0** - Initial release (TBD)

## Upgrade Guide

### From Nothing to v0.1.0

This is the initial release. To get started:

1. Install dependencies
2. Configure LLM API key
3. Run the service
4. Send planning requests

See [README.md](../README.md) for detailed setup instructions.

## Breaking Changes

None (initial release)

## Deprecation Notices

None

## Future Plans

### v0.2.0 (Planned)
- OpenAI provider support
- Few-shot learning examples in prompts
- Confidence scoring for generated graphs
- Graph caching for similar tasks
- Streaming API for progress updates
- WebSocket support
- Enhanced metrics and observability

### v0.3.0 (Planned)
- Multi-model fallback
- Graph optimization suggestions
- Custom node type support
- Advanced routing strategies
- Integration with dago orchestrator

### v1.0.0 (Planned)
- Production-ready stability
- Full test coverage
- Performance optimizations
- Complete documentation
- Production deployment guides

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## Release Process

1. Update version in code
2. Update CHANGELOG.md
3. Create git tag
4. Push tag to trigger release workflow
5. GitHub Actions builds and publishes artifacts
6. Docker image pushed to Docker Hub

## Support

For issues and questions:
- GitHub Issues: https://github.com/aescanero/dago-node-planner/issues
- Documentation: https://github.com/aescanero/dago-node-planner/docs
