# music.jelly.ninja 🌀

AI-generated Strudel live coding patterns via SSH.

## Connect

```bash
ssh music.jelly.ninja | strudel
```

That's it! Opens strudel.cc in your browser with a fresh pattern.

## Self-Host with Cloudflare Tunnel

Run locally and expose to the world:

```bash
# 1. Install cloudflared
brew install cloudflared  # macOS
# or: curl -L https://github.com/cloudflare/cloudflared/releases/download/2024.1.5/cloudflared-darwin-arm64 -o /usr/local/bin/cloudflared

# 2. Create tunnel
cloudflared tunnel create music-jelly-ninja

# 3. Configure DNS
cloudflared tunnel route dns music-jelly-ninja music.jelly.ninja

# 4. Run server + tunnel
cd ssh && go build -o strudel-wish
./strudel-wish --port 23234 &

cloudflared tunnel run --config ../cloudflared.yml music-jelly-ninja
```

## Deploy to VPS

```bash
cd ssh
go build -o strudel-wish
./strudel-wish --port 22
```
