# cargohold

A lightweight CLI for managing environment-specific secret bundles without a full secrets manager.

---

## Installation

```bash
go install github.com/yourname/cargohold@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cargohold.git && cd cargohold && go build -o cargohold .
```

---

## Usage

Initialize a new secret bundle for an environment:

```bash
cargohold init --env production
```

Add a secret:

```bash
cargohold set --env production DB_PASSWORD=supersecret
```

Retrieve a secret:

```bash
cargohold get --env production DB_PASSWORD
```

Export all secrets for an environment as shell exports:

```bash
eval $(cargohold export --env production)
```

Bundles are encrypted at rest using a local key. Keep your keyfile safe — cargohold does not manage key recovery.

---

## Configuration

By default, cargohold stores bundles in `~/.cargohold/`. Override with the `CARGOHOLD_DIR` environment variable.

---

## License

MIT © yourname