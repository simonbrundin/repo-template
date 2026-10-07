# Development Environments

Detta repo stödjer flera verktyg för lokal Kubernetes-utveckling.

## Alternativ 1: Tilt (nuvarande)

Webbaserat UI för att övervaka och hantera resurser.

```bash
# Installera Tilt
brew install tilt

# Starta
simon dev

# Eller direkt
tilt up
```

**Fördelar:**
- Bra webb-GUI
- Live update
- Loggvisning i realtid

**Nackdelar:**
- Starlark-konfiguration (Python-liknande)
- Komplext att debugga

---

## Alternativ 2: Skaffold

Declarativ YAML-konfiguration med CLI.

```bash
# Installera Skaffold
brew install skaffold

# Starta med hot-reload
skaffold dev

# Eller utan cleanup vid avslut
skaffold dev --cleanup=false
```

**Fördelar:**
- YAML-konfiguration (enkel att läsa)
- Bra CI/CD-integration
- Stöd för Jib, Buildpacks, Bazel

**Nackdelar:**
- CLI utan GUI (kan vara minus för vissa)

**Konfiguration:** `skaffold.yaml`

---

## Alternativ 3: DevSpace

CLI + GUI med fokus på team-samarbete.

```bash
# Installera DevSpace
brew install devspace

# Initiera (om inte redan gjort)
devspace generate

# Starta dev-läge
devspace dev

# Eller med specifik konfiguration
devspace dev -c devspace.yaml
```

**Fördelar:**
- GUI finns (localhost)
- Bra för team-workflows
- Blue/green deployments

**Nackdelar:**
- Mer komplext setup
- Kan teara down resurser vid avslut

**Konfiguration:** `devspace.yaml`

---

## Alternativ 4: Mirrord

Kör lokalkod direkt i cluster-kontext (INGEN image-building!).

```bash
# Installera Mirrord
brew install metalbear-co/mirrord/mirrord

# I valfri Go/Node/Python-app, kör med:
mirrord run -- <your-app-command>

# Eller som VS Code extension:
# Sök efter "Mirrord" i Extensions
```

**Fördelar:**
- ✅ INGEN image-building!
- ✅ INGEN push till registry!
- ✅ Lokal kod kör i cluster-kontext
- ✅ Network-traffic interceptas
- ✅ Extremt snabbt feedback-loop

**Nackdelar:**
- Endast en tjänst åt gången
- Kräver att cluster finns
- Mindre kontroll över hela miljön

**Användningsfall:**
```bash
# Kör din lokala Go-api i cluster
cd src/go
mirrord run -- go run cmd/api/main.go

# Kör din lokala Nuxt i cluster
cd src/nuxt  
mirrord run -- npm run dev
```

---

## Jämförelse

| Verktyg | GUI | Konfiguration | Hot-reload | CI/CD |
|---------|-----|--------------|-------------|-------|
| Tilt | ✅ | Starlark | ✅ | ⚠️ |
| Skaffold | ❌ | YAML | ✅ | ✅ |
| DevSpace | ✅ | YAML | ✅ | ⚠️ |
| Mirrord | ❌ | - | ✅ (via språk) | ❌ |

---

## Rekommendation

- **Bästa för enskild utveckling:** Mirrord (snabbast, inga images)
- **Bästa för team:** Skaffold (enkel YAML, CI/CD)
- **Bästa för GUI-älskare:** Tilt/DevSpace

Testa gärna alla och se vilket som passar din workflow bäst!
