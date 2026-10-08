# Mirrord - Local Development with Remote Pod

Mirrord låter dig köra lokala processer med trafik från en fjärrpod.

## Installation

```bash
brew install metalbear-co/mirrord/mirrord
```

## Användning

### Grundläggande kommando

```bash
# Kör ett kommando med Mirrord (speglar trafik från pod)
mirrord exec --target deployment/nuxt-app -n default -- <kommando>

# Exempel: testa anslutning
mirrord exec --target deployment/nuxt-app -n default -- curl http://localhost
```

### Med lokalt dev-server

1. Starta din lokala dev-server:
```bash
cd src/nuxt
npm run dev
```

2. I en annan terminal, använd Mirrord för att spegla trafik:
```bash
mirrord exec --target deployment/nuxt-app -n default -- curl http://localhost:3000
```

## Fördelar

- ✅ **Ingen image-building** - kör lokala ändringar direkt
- ✅ **Fil-system-åtkomst** - läs/skriv filer i podden
- ✅ **Miljövariabler** - få samma env som i prod
- ✅ **Enkel setup** - bara `mirrord exec`

## Nackdelar

- ❌ Behöver en fungerande pod att spegla
- ❌ Trafikroutning kan vara förvirrande
- ❌ Inte lika "magiskt" som Tilt/DevSpace för full hot-reload

## Target-typer

```bash
# Pod
mirrord exec --target pod/min-pod -n default -- <cmd>

# Deployment
mirrord exec --target deployment/api -n default -- <cmd>

# Med container
mirrord exec --target deployment/api/container/main -n default -- <cmd>

# Label selector
mirrord exec --target label/app=api -n default -- <cmd>
```

## TUI (Interaktiv)

```bash
mirrord tui
```
