**BlockIRC**:

---

# BlockIRC

BlockIRC is a modern and secure IRC (Internet Relay Chat) server/daemon written in Go. It is designed to support blockchain and distributed systems projects, offering enhanced security, privacy, and scalability. BlockIRC aims to integrate seamlessly with blockchain ecosystems, providing a reliable communication platform for decentralized applications and services.

## Key Features

- **Secure Communication**: Supports TLS for encrypted communication.
- **Decentralized Integration**: Compatible with various blockchain protocols for secure and verifiable messaging.
- **Scalable Architecture**: Designed for high availability and scalability, suitable for large-scale deployments.
- **Extensible and Customizable**: Easily extendable with custom modules and configurations to fit specific needs.
- **Audit Logging**: Comprehensive logging features for monitoring and auditing chat activities.

## Use Cases

1. **Decentralized Autonomous Organizations (DAOs)**: Secure and verifiable communication channels for DAO members.
2. **Blockchain Development Communities**: Real-time collaboration and discussion for developers working on blockchain projects.
3. **Decentralized Finance (DeFi)**: Coordination and support for DeFi projects and users.
4. **Supply Chain Management**: Secure communication channels for supply chain participants in blockchain-based systems.

## Installation

### Docker

To build and run BlockIRC using Docker, follow these steps:

1. **Clone the repository:**
   ```sh
   git clone https://github.com/EncrypteDL/BlockIRC.git
   cd BlockIRC
   ```

2. **Build the Docker image:**
   ```sh
   docker build -t blockirc:latest .
   ```

3. **Run the container:**
   ```sh
   docker run -d -p 6667:6667 -p 6697:6697 blockirc:latest
   ```

### Docker Compose

To run BlockIRC using Docker Compose:

1. **Create a `docker-compose.yml` file:**

```yaml
version: "3.3"

services:
  blockirc:
    image: blockirc:latest
    configs:
      - source: ircd_yml
        target: /ircd.yml
      - source: ircd_motd
        target: /ircd.motd
      - source: cert_pem
        target: /cert.pem
      - source: key_pem
        target: /key.pem
    ports:
      - target: 6667
        published: 6667
        protocol: tcp
        mode: host
      - target: 6697
        published: 6697
        protocol: tcp
        mode: host
    deploy:
      endpoint_mode: dnsrr
      restart_policy:
        condition: on-failure
      replicas: 1

configs:
  ircd_yml:
    file: ./ircd.yml
  ircd_motd:
    file: ./ircd.motd
  cert_pem:
    file: ./cert.pem
  key_pem:
    file: ./key.pem
```

2. **Deploy the service:**
   ```sh
   docker-compose up -d
   ```

## Contributing

We welcome contributions from the community. Please read our [contributing guidelines](CONTRIBUTING.md) for more details.

## License

BlockIRC is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.

---