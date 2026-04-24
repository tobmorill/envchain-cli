# envchain-cli

A CLI tool for managing per-project environment variable sets with encrypted local storage and shell integration.

---

## Installation

```bash
go install github.com/yourname/envchain-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain-cli.git && cd envchain-cli && go build -o envchain .
```

---

## Usage

**Add a variable to a project chain:**
```bash
envchain set myproject AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
```

**Load variables into your current shell session:**
```bash
eval $(envchain load myproject)
```

**List all chains and their keys:**
```bash
envchain list
```

**Remove a variable:**
```bash
envchain unset myproject AWS_ACCESS_KEY_ID
```

Variables are encrypted at rest using AES-256 and stored locally in `~/.envchain/`. Each project maintains its own isolated chain, keeping credentials scoped and organized.

---

## Shell Integration

Add the following to your `.bashrc` or `.zshrc` to auto-load a chain when entering a project directory:

```bash
export ENVCHAIN_AUTO=1
```

---

## License

MIT © [yourname](https://github.com/yourname)