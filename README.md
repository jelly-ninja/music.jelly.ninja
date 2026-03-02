# music.jelly.ninja 🌀

AI-generated Strudel live coding patterns via SSH.

## Connect

```bash
ssh music.jelly.ninja | strudel
```

That's it! Opens strudel.cc in your browser with a fresh pattern.

**Port:** Default SSH (22)

## Requirements

- SSH client
- A browser

No install needed.

## How It Works

```
ssh music.jelly.ninja → Returns Strudel pattern
                    → strudel CLI opens browser with pattern
                    → Listen in strudel.cc!
```

## Development

### SSH Server (Go)
```bash
cd ssh
go build -o strudel-wish
./strudel-wish --port 23234
```

### CLI (Go)
```bash
cd cli-go
go build -o strudel
./strudel
```

### With AI Generation
```bash
export AI_GATEWAY_API_KEY="your-key"
# Daily refresh (default)
./strudel-wish
```

### Docker
```bash
docker build -t music-jelly-ninja ./ssh
docker run -p 23234:23234 -e AI_GATEWAY_API_KEY=your-key music-jelly-ninja
```

## License

AGPL-3.0
