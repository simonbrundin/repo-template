# Repository Template

Detta repo innehåller applikationskod och GitOps-manifest för pr-flow.

---

## Översikt

```mermaid
flowchart LR
    subgraph APP_REPO["📁 pr-flow (detta repo)"]
        CODE["frontend/"]
        MANIFESTS["environments/"]
    end

    subgraph INFRA_REPO["📁 infrastructure"]
        FLUX_CONFIG["infrastructure-flux/apps/pr-flow/"]
    end

    subgraph K8S["☸ Kubernetes"]
        FLUX["Flux"]
        CLUSTER["Kluster"]
    end

    CODE -->|bygger| IMAGE
    IMAGE["GHCR Image"]
    MANIFESTS -->|pekar på| FLUX_CONFIG
    FLUX_CONFIG -->|GitRepository + Kustomization| FLUX
    FLUX -->|deployar| CLUSTER
    IMAGE -->|används av| CLUSTER
```

---

## Repositories

### pr-flow (detta repo)

Innehåller:
- **Applikationskod** i `frontend/`
- **GitOps-manifest** i `environments/`
- **CI/CD** i `.github/workflows/`

### infrastructure

Innehåller:
- **Flux-konfiguration** i `infrastructure-flux/apps/pr-flow/`
- **Kubernetes-komponenter** i `infrastructure-flux/components/`

---

## Mappstruktur

```
pr-flow/
├── .github/
│   └── workflows/
│       ├── ci-for-js-app.yaml        # Kör tester + bygger image
│       └── pr-preview.yaml            # Skapar PR preview overlays
│
├── frontend/                          # Applikationskod
│   └── Dockerfile                    # Byggs till GHCR image
│
└── environments/
    ├── base/                         # Gemensamma manifests (för alla miljöer)
    │   ├── kustomization.yaml
    │   ├── namespace.yaml
    │   ├── deployment.yaml
    │   ├── service.yaml
    │   └── httproute.yaml
    │
    ├── production/                   # Produktions-overlay
    │   ├── kustomization.yaml       # resources: [../base]
    │   └── namespace.yaml           # (unik: environment: production)
    │
    └── pr/                          # PR Preview overlays (dynamiska)
        ├── kustomization.yaml       # [./157, ./158, ...] ← uppdateras av CI
        └── <PR_NUM>/
            └── kustomization.yaml   # resources: [../../base] + patches
```

---

## Kustomize-arkitektur

```
┌─────────────────────────────────────────────────────────────┐
│                      environments/                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────┐                                           │
│   │    base/    │  ← Gemensamma manifests                   │
│   │             │    (namespace, deployment, service, etc.) │
│   │  deployment │                                           │
│   │  service    │                                           │
│   │  httproute  │                                           │
│   └──────┬──────┘                                           │
│          │                                                   │
│          │  refereras av                                     │
│          ▼                                                   │
│   ┌─────────────────────────────────────────────────────┐  │
│   │                 Overlays                              │  │
│   ├─────────────────────────────────────────────────────┤  │
│   │                                                       │  │
│   │  production/        pr/157/        pr/158/           │  │
│   │  ├─ kustomize      ├─ kustomize    ├─ kustomize    │  │
│   │  └─ namespace      └─ patches       └─ patches       │  │
│   │                                                       │  │
│   └─────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Base (`environments/base/`)

Innehåller standardmanifest för appen. Används av alla miljöer.

### Production (`environments/production/`)

Produktions-specifika overrides:
```yaml
# environments/production/kustomization.yaml
resources:
  - ../base
```
```yaml
# environments/production/namespace.yaml (override)
apiVersion: v1
kind: Namespace
metadata:
  name: pr-flow
  labels:
    environment: production  # ← Unik för prod
```

### PR Previews (`environments/pr/<NUM>/`)

Dynamiskt skapade av GitHub Actions:
```yaml
# environments/pr/157/kustomization.yaml
resources:
  - ../../base

namespace: pr-157  # ← Unik per PR

patches:
  # Image-patch: använder PR-specifik image
  # Hostname-patch: pr-157.example.com
  # Labels: preview: "true"
```

---

## Hur Flux hittar manifesten

### I `infrastructure/infrastructure-flux/apps/pr-flow/`

```
infrastructure-flux/apps/pr-flow/
├── gitrepository.yaml        # ← Flux läser pr-flow repo
├── flux-kustomization.yaml  # ← Två Kustomizations
└── kustomization.yaml
```

```yaml
# gitrepository.yaml
apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: pr-flow
  namespace: flux-system
spec:
  url: ssh://git@github.com/simonbrundin/pr-flow
  ref:
    branch: main
```

```yaml
# flux-kustomization.yaml
---
# Produktion
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: pr-flow
spec:
  path: environments/production  # ← Läser pr-flow/environments/production
  sourceRef:
    kind: GitRepository
    name: pr-flow

---
# PR Previews
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: pr-flow-previews
spec:
  path: ./environments/pr       # ← Läser pr-flow/environments/pr
  sourceRef:
    kind: GitRepository
    name: pr-flow
  prune: true                  # Städar bort resurser vid borttagning
```

---

## CI/CD Flöde

### Pull Request

```mermaid
flowchart TB
    subgraph DEVELOPER["👤 Developer"]
        PUSH[git push]
        PR[Skapa PR]
    end

    subgraph CI["⚡ GitHub Actions"]
        TEST[Test & Lint]
        BUILD[docker build]
        PUSH_IMAGE[Push to GHCR]
        CREATE_OVERLAY[Skapar pr/<NUM>/]
        UPDATE_KUSTOMIZE[Uppdaterar pr/kustomization.yaml]
        COMMIT[git commit & push]
    end

    subgraph GIT["📁 Git Repository"]
        ENV_PR["environments/pr/<NUM>/"]
        ENV_BASE["environments/base/"]
        ENV_PR_KUSTOMIZE["environments/pr/kustomization.yaml"]
    end

    subgraph CLUSTER["☸ Kubernetes"]
        FLUX[Flux]
        NS["Namespace: pr-<NUM>"]
        DEPLOY["Deployment"]
        ROUTE["HTTPRoute"]
    end

    PUSH --> PR
    PR --> TEST
    TEST --> BUILD
    BUILD --> PUSH_IMAGE
    PUSH_IMAGE --> CREATE_OVERLAY
    CREATE_OVERLAY --> UPDATE_KUSTOMIZE
    UPDATE_KUSTOMIZE --> COMMIT
    COMMIT --> GIT
    
    ENV_PR_KUSTOMIZE --> FLUX
    ENV_PR --> FLUX
    ENV_BASE --> ENV_PR
    FLUX --> NS
    NS --> DEPLOY
    DEPLOY --> ROUTE
```

### Steg-för-steg

| Steg | Action | Var |
|------|--------|-----|
| 1 | Developer skapar PR | GitHub |
| 2 | GitHub Actions triggas | `.github/workflows/pr-preview.yaml` |
| 3 | Kör tester | CI |
| 4 | Bygger image | `frontend/Dockerfile` |
| 5 | Pushar till GHCR | Tag: `pr-<NUM>-<sha>` |
| 6 | Skapar `environments/pr/<NUM>/` | med base + patches |
| 7 | Uppdaterar `environments/pr/kustomization.yaml` | Lägger till `./<NUM>` |
| 8 | Commit & push | Till PR-branchen |
| 9 | Flux ser ändringarna | `infrastructure-flux/apps/pr-flow/` |
| 10 | Flux deployar | `pr-<NUM>` namespace |

### PR Stängning

| Steg | Action |
|------|--------|
| 1 | Developer stänger PR |
| 2 | GitHub Actions triggas |
| 3 | Tar bort `environments/pr/<NUM>/` |
| 4 | Uppdaterar `environments/pr/kustomization.yaml` |
| 5 | Flux tar bort resurser (prune: true) |

---

## GitOps-principen

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   GitHub Actions                   Flux                      │
│                                                             │
│   ┌───────────────┐              ┌───────────────┐         │
│   │  Bygger image │   push       │  Läser Git    │         │
│   │  Commit:ar    │ ──────────▶  │  Deployar     │         │
│   │  manifests    │              │  Kubernetes   │         │
│   └───────────────┘              └───────────────┘         │
│                                                             │
│   ÄNDRAR INTE KLASTRET DIREKT    LÄSER ENBART FRÅN GIT    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Regler:**
1. GitHub Actions **aldrig** applicerar direkt till Kubernetes
2. Flux **aldrig** ändrar Git
3. Allt går via Git

---

## Miljöer

| Miljö | Källa | Styrning | URL |
|-------|-------|----------|-----|
| **production** | `environments/production/` | Flux | `prflow.example.com` |
| **preview** | `environments/pr/<NUM>/` | Flux + GitHub Actions | `pr-<NUM>.example.com` |
| **dev** | lokalt | Tilt | `localhost:3000` |

---

## Infrastructure-koppling

```
infrastructure/
└── infrastructure-flux/
    └── apps/
        └── pr-flow/                    # ← Konfiguration för denna app
            ├── gitrepository.yaml       # Läser pr-flow repo
            └── flux-kustomization.yaml  # Två Kustomizations:
                                         #   - pr-flow (production)
                                         #   - pr-flow-previews (PRs)
```

### Lägga till en ny app

1. Skapa mapp i `infrastructure-flux/apps/<app>/`
2. Lägg till `gitrepository.yaml`
3. Lägg till `flux-kustomization.yaml`
4. Uppdatera `infrastructure-flux/apps/kustomization.yaml`

---

## Preview URL

```
https://pr-<NUM>.example.com
```

**Förutsättningar:**
- `external-dns` konfigurerad i klustret
- Gateway `traefik-gateway` finns i `flux-system`
- DNS pekar mot klustret

---

## Architecture Overview

```mermaid
flowchart TB
    subgraph REPO["Git Repository / Source Code"]
        APP_SOURCE["Application Source Code"]

        subgraph GITOPS["GitOps"]
            BASE["environments/base"]
            DEV_GITOPS["environments/dev"]
            PR_GITOPS["environments/pr"]
            PROD_GITOPS["environments/prod"]
        end
    end

    subgraph LAPTOP[" "]
        LAPTOP_TITLE["💻 LAPTOP"]

        subgraph DEVENV[" "]
            DEV_TITLE["Developer Environment"]

            IDE["Neovim / VS Code / IDE"]
            TOOLS["Development Tools"]

            subgraph WT1["Worktree: main"]
                WORKTREE1["Source Code"]
            end

            subgraph WT2["Worktree: feature-login"]
                WORKTREE2["Source Code"]
            end

            IDE --> WORKTREE1
            IDE --> WORKTREE2
            TOOLS --> WORKTREE1
            TOOLS --> WORKTREE2
        end
    end

    APP_SOURCE --> WORKTREE1
    APP_SOURCE --> WORKTREE2

    WORKTREE1 --> TILT1["Tilt"]
    WORKTREE2 --> TILT2["Tilt"]

    subgraph K8S[" "]
        K8S_TITLE["☸ TALOS KUBERNETES CLUSTER"]

        subgraph DEVROW[" "]
            subgraph DEV1["🔵 Namespace: dev-main"]
                APP1["Application"]
                DB1["Database"]
                ROUTE1["HTTPRoute"]
            end

            subgraph DEV2["🔵 Namespace: dev-feature-login"]
                APP2["Application"]
                DB2["Database"]
                ROUTE2["HTTPRoute"]
            end
        end

        subgraph DELIVERYROW[" "]
            subgraph PR["🟡 Namespace: pr-123"]
                PRAPP["Application"]
                PRDB["Database"]
                PRROUTE["HTTPRoute"]
            end

            subgraph PROD["🟢 Namespace: production"]
                PRODAPP["Application"]
                PRODDB["Database"]
                PRODROUTE["HTTPRoute"]
            end
        end
    end

    TILT1 --> DEV1
    TILT2 --> DEV2

    WORKTREE2 --> PRREQ["Pull Request"]

    subgraph DELIVERY[" "]
        DELIVERY_TITLE["CI / GitOps"]

        CI["CI"]
        IMAGE["Container Image"]
        REGISTRY["Container Registry"]
        FLUX["Flux"]

        PRREQ --> CI
        CI --> IMAGE
        IMAGE --> REGISTRY
    end

    BASE -.-> DEV_GITOPS
    BASE -.-> PR_GITOPS
    BASE -.-> PROD_GITOPS
    DEV_GITOPS -.-> DEV1
    PR_GITOPS -.-> FLUX
    PROD_GITOPS -.-> FLUX

    REGISTRY --> FLUX
    FLUX --> PR
    PR --> REVIEW["Review / Merge"]
    REVIEW --> FLUX
    FLUX --> PROD

    classDef base fill:#8b5cf6,stroke:#7c3aed,color:#ffffff
    classDef dev fill:#3b82f6,stroke:#1d4ed8,color:#ffffff
    classDef pr fill:#facc15,stroke:#ca8a04,color:#111827
    classDef prod fill:#22c55e,stroke:#15803d,color:#ffffff

    class BASE base
    class APP1,DB1,ROUTE1,APP2,DB2,ROUTE2 dev
    class PRAPP,PRDB,PRROUTE pr
    class PRODAPP,PRODDB,PRODROUTE prod

    classDef title fill:none,stroke:none,color:#0f172a,font-size:18px,font-weight:bold
    class LAPTOP_TITLE,DEV_TITLE,K8S_TITLE,DELIVERY_TITLE title

    style REPO fill:#f1f5f9,stroke:#334155,stroke-width:3px,color:#0f172a
    style GITOPS fill:#ffffff,stroke:#64748b,stroke-width:2px,color:#0f172a
    style LAPTOP fill:#e2e8f0,stroke:#334155,stroke-width:3px
    style DEVENV fill:#ffffff,stroke:#64748b,stroke-width:2px
    style K8S fill:#f8fafc,stroke:#334155,stroke-width:3px
    style DELIVERY fill:#ffffff,stroke:#94a3b8
    style DEVROW fill:none,stroke:none
    style DELIVERYROW fill:none,stroke:none
```

## Databasschema

```mermaid
erDiagram
    %!generate db diagram from schema.sql
```
