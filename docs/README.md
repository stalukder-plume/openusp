# OpenUSP Documentation

This directory contains comprehensive documentation for the OpenUSP platform, covering both TR-069 (CWMP) and TR-369 (USP) implementations.

## Documentation Files

### 📋 [Use Cases](use-cases.md)
Comprehensive use case documentation covering both TR-069 and TR-369 protocols:
- **TR-069 (CWMP) Use Cases**: Legacy device management, service provisioning, monitoring
- **TR-369 (USP) Use Cases**: Next-generation IoT, multi-controller scenarios, advanced automation  
- **Integration Scenarios**: Hybrid deployments, migration strategies, best practices
- **Implementation Examples**: API usage, monitoring, event-driven automation

### 🏗️ [Controller Architecture](controller-architecture.md) 
Technical deep-dive into the OpenUSP Controller Subsystem:
- **System Architecture**: Component overview, message flows, state management
- **Protocol Controllers**: TR-069 ACS and TR-369 USP controller implementations
- **Message Processing**: SOAP and Protocol Buffer message handling
- **Advanced Scenarios**: Multi-service orchestration, compliance monitoring, QoS management
- **Integration Patterns**: Event-driven architecture, workflow engine, external system integration

### 📡 [TR-069 CWMP ACS](TR069_CWMP_ACS.md)
Specific documentation for TR-069 CWMP ACS implementation:
- CWMP protocol specifics
- ACS server configuration
- Device management workflows
- SOAP message handling

## Quick Navigation

| Topic | Document | Description |
|-------|----------|-------------|
| **Business Use Cases** | [use-cases.md](use-cases.md) | When and how to use TR-069 vs TR-369 |
| **Technical Architecture** | [controller-architecture.md](controller-architecture.md) | System design and implementation |
| **CWMP ACS Details** | [TR069_CWMP_ACS.md](TR069_CWMP_ACS.md) | TR-069 specific implementation |
| **API Reference** | [../api/README.md](../api/README.md) | REST API documentation |
| **Getting Started** | [../README.md](../README.md) | Project overview and quick start |

## Use Case Categories

### 🏠 Residential/Consumer
- **TR-069**: Home gateways, modems, set-top boxes, VoIP devices
- **TR-369**: Smart home devices, IoT sensors, next-gen gateways

### 🏢 Enterprise/Business  
- **TR-069**: Network equipment, managed services, bulk configurations
- **TR-369**: SD-WAN devices, edge computing, multi-tenant management

### 🌐 Service Provider
- **TR-069**: Legacy infrastructure management, compliance, service assurance
- **TR-369**: 5G edge, network slicing, multi-controller federation

### 🏭 Industrial/IoT
- **TR-069**: Legacy industrial equipment with CWMP support
- **TR-369**: Modern IIoT, predictive maintenance, real-time automation

## Protocol Comparison

| Aspect | TR-069 (CWMP) | TR-369 (USP) | Use Case Recommendation |
|--------|---------------|--------------|-------------------------|
| **Maturity** | Mature (2004) | Modern (2018) | TR-069 for existing, USP for new |
| **Transport** | HTTP/SOAP | Multi-transport | USP for flexibility |
| **Security** | Basic TLS | Enhanced E2E | USP for security-critical |
| **Multi-Controller** | Single ACS | Multi-controller | USP for complex orchestration |
| **Real-time** | Limited | Enhanced | USP for low-latency requirements |
| **Ecosystem** | Extensive | Growing | TR-069 for immediate compatibility |

## Implementation Guidance

### When to Choose TR-069
- ✅ Managing existing/legacy device deployments
- ✅ Simple configuration and monitoring requirements  
- ✅ Regulatory compliance with established standards
- ✅ Working with mature vendor ecosystems
- ✅ Cost-sensitive deployments leveraging existing infrastructure

### When to Choose TR-369 (USP)
- ✅ New device deployments and greenfield projects
- ✅ Complex automation and orchestration requirements
- ✅ Multi-vendor, heterogeneous device environments
- ✅ Scalable, distributed management needs
- ✅ Advanced security and privacy requirements
- ✅ Real-time service adaptation and event-driven workflows

### Hybrid Approach
- 🔄 **Gradual Migration**: Start with USP for new services while maintaining TR-069 for existing
- 🔄 **Service Segmentation**: Use different protocols for different service tiers
- 🔄 **Geographic Rollout**: Deploy USP in new markets while supporting TR-069 in established areas

## Getting Started

1. **Read the Use Cases**: Start with [use-cases.md](use-cases.md) to understand business scenarios
2. **Understand Architecture**: Review [controller-architecture.md](controller-architecture.md) for technical details  
3. **Try the APIs**: Use the [API documentation](../api/README.md) and Swagger UI for hands-on testing
4. **Deploy with Docker**: Follow the [main README](../README.md) for Docker Compose deployment
5. **Explore Examples**: Check the implementation examples in the use cases document

## Contributing to Documentation

We welcome contributions to improve and extend this documentation:

1. **Use Cases**: Add new scenarios, implementation patterns, or industry-specific examples
2. **Technical Details**: Contribute architecture improvements, code examples, or best practices  
3. **Tutorials**: Create step-by-step guides for specific deployment scenarios
4. **Troubleshooting**: Document common issues and solutions

Please see the main project [contributing guidelines](../CONTRIBUTING.md) for submission process.

## Additional Resources

- **OpenUSP Project**: [Main Repository](/)
- **USP Specification**: [Broadband Forum USP](https://usp.technology)  
- **TR-069 Specification**: [Broadband Forum TR-069](https://www.broadband-forum.org/technical/download/TR-069.pdf)
- **API Documentation**: [OpenAPI Specification](../api/openusp.yaml)
- **Docker Deployment**: [Docker Compose Setup](../deployments/README.md)
