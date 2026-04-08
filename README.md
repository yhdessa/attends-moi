# attends-moi

A lightweight Kanban board application for development teams, designed to showcase modern cloud-native practices.

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Backend** | Go, Chi Router, pgx, slog |
| **Frontend** | Vue 3, TypeScript, TailwindCSS, daisyUI |
| **Database** | PostgreSQL 16 |
| **Infrastructure** | Kubernetes, Minikube, Kustomize, Docker |

## ✨ Features

- **Kanban Boards**: Create and manage multiple project boards.
- **Task Management**: CRUD operations for cards with status, priority, labels, and assignee.
- **Comments**: Discuss tasks directly on cards.
- **Real-time UI**: Responsive interface with a dark theme (Dracula).
- **Cloud-Native Ready**: Health checks, graceful shutdown, resource limits, and Kustomize overlays.

## 📂 Project Structure

```
attends-moi/
├── apps/
│   ├── backend/           # Go REST API
│   └── frontend/          # Vue 3 SPA
├── k8s/
│   ├── base/              # Base Kubernetes manifests
│   └── overlays/minikube/ # Environment-specific patches
├── scripts/
│   ├── setup-minikube.sh  # Cluster initialization
│   └── deploy.sh          # Build and deploy pipeline
└── docker-compose.yml     # Local development stack
```

## 🚀 Getting Started

### Prerequisites

- **Docker** & **Docker Compose**
- **Minikube** (for K8s deployment)
- **kubectl**
- **Go 1.22+** (for local backend dev)
- **Node.js 20+** (for local frontend dev)

### Local Development

1. **Start the database:**
   ```bash
   docker run -d --name postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=attends_moi \
     -p 5432:5432 postgres:16-alpine
   ```

2. **Run migrations:**
   ```bash
   for f in apps/backend/migrations/*.sql; do
     psql -U postgres -h localhost -d attends_moi -f "$f"
   done
   ```

3. **Start backend:**
   ```bash
   cd apps/backend
   cp .env.example .env
   go mod tidy
   go run cmd/server/main.go
   ```

4. **Start frontend:**
   ```bash
   cd apps/frontend
   npm install
   npm run dev
   ```

### Kubernetes Deployment (Minikube)

1. **Setup cluster:**
   ```bash
   ./scripts/setup-minikube.sh
   ```

2. **Build and deploy:**
   ```bash
   ./scripts/deploy.sh
   ```

3. **Access the application:**
   ```bash
   minikube service frontend -n attends-moi --url
   ```

## 📡 API Endpoints

### Boards
- `GET /api/boards` — List all boards
- `POST /api/boards` — Create a new board
- `GET /api/boards/{id}` — Get board details

### Cards
- `GET /api/boards/{id}/cards` — List cards on a board
- `POST /api/boards/{id}/cards` — Create a card
- `GET /api/cards/{id}` — Get card details
- `PATCH /api/cards/{id}` — Update card (status, priority, etc.)
- `DELETE /api/cards/{id}` — Delete a card

### Comments
- `GET /api/cards/{id}/comments` — List comments
- `POST /api/cards/{id}/comments` — Add a comment

### Health
- `GET /health` — Liveness probe
- `GET /ready` — Readiness probe (checks DB connection)

## 🏗️ Kubernetes Architecture

```
[Ingress] → [frontend:80]
         ↘ [backend:8080] → [postgres:5432]
```

- **Kustomize**: Used for environment-specific configurations (`base` vs `minikube`).
- **Probes**: HTTP health checks for backend, `pg_isready` for Postgres.
- **Resources**: CPU/memory requests and limits defined for all pods.
- **Persistence**: PVC for Postgres data durability.

## 📝 Future Improvements

- [ ] WebSocket support for real-time updates
- [ ] Authentication & RBAC
- [ ] CI/CD pipeline with GitHub Actions
- [ ] Helm chart for simplified deployment
- [ ] Prometheus metrics & Grafana dashboards
