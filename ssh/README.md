# Strudel Wish

An SSH server that streams Strudel live coding patterns to clients.

## Quick Start

### Build

```bash
go build -o strudel-wish
```

### Run

```bash
./strudel-wish --port 23234
```

### Connect

```bash
ssh strudel@your-server.com -p 23234
```

### Play

```bash
ssh strudel@your-server.com -p 23234 | strudel
```

## Installation

### Option 1: From source

```bash
git clone https://github.com/yourusername/strudel-wish.git
cd strudel-wish/ssh
go build -o strudel-wish
```

### Option 2: Docker

See Dockerfile section below.

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `0.0.0.0` | Host to listen on |
| `--port` | `23234` | Port to listen on |
| `--seed` | random | Random seed for pattern selection |

## Client Setup

Install the Strudel CLI:

```bash
npm install -g strudel
# or
bunx install strudel
```

Or use npx:

```bash
npx strudel
```

## Deployment

### Docker

```bash
docker build -t strudel-wish ./ssh
docker run -p 23234:23234 strudel-wish
```

### Systemd

Create `/etc/systemd/system/strudel-wish.service`:

```ini
[Unit]
Description=Strudel Wish SSH Server
After=network.target

[Service]
Type=simple
User=strudel
WorkingDirectory=/opt/strudel-wish
ExecStart=/opt/strudel-wish/strudel-wish --port 23234
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## Usage

1. Start the server
2. Users connect via SSH
3. They receive a random Strudel pattern
4. Pipe to local strudel CLI to play

Example:
```bash
# On server (assuming strudel-wish is running on port 23234)
ssh user@server -p 23234 | strudel
```
