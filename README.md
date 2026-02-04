# SRE GitOps Lab: The "OpenSource Dynatrace" Experiment 🚀

Welcome to my ultimate SRE playground. This repository documents my journey in building a fully-featured, enterprise-grade Observability and Reliability platform using strictly OpenSource technologies. 

Running locally on **Kind** (Kubernetes in Docker) hosted on a **Proxmox** home server, fully managed via **GitOps** methodology.

## 🎯 Architecture Goals

The goal is to replicate the capabilities of proprietary commercial tools (like Dynatrace) using a modern CNCF ecosystem stack:

*   **GitOps Core**: ArgoCD for continuous delivery.
*   **Progressive Delivery**: Blue/Green & Canary deployments using Argo Rollouts.
*   **Full-Stack Observability (The "LGTM" Stack)**:
    *   **Metrics**: Prometheus / VictoriaMetrics.
    *   **Logs**: Grafana Loki.
    *   **Traces**: Grafana Tempo + OpenTelemetry.
    *   **Visualization**: Grafana.
    *   **Auto-Instrumentation**: OpenTelemetry Operator (Zero-code instrumentation).
*   **Reliability & Chaos**: Chaos Mesh / Litmus to test system resilience.
*   **User Journeys**: Synthetic monitoring and SLI/SLO tracking.

## 🛠️ Tech Stack & Roadmap

| Component | Technology | Status |
| :--- | :--- | :--- |
| **Cluster** | Kind on Proxmox | ✅ Ready |
| **GitOps** | ArgoCD | ✅ Ready |
| **Deployment** | Argo Rollouts | 🚧 Planned |
| **Metrics** | Prometheus | 🚧 Planned |
| **Tracing** | OpenTelemetry + Tempo | 🚧 Planned |
| **Logs** | Loki | 🚧 Planned |
| **Chaos** | Chaos Mesh | 🚧 Planned |
| **Dashboards** | Grafana as Code | 🚧 Planned |

## 📂 Repository Structure

```
.
├── apps/               # Application charts and manifests
│   ├── web-demo/       # Sample generic app for testing release strategies
│   └── observability/  # The Monitoring Stack (Prometheus, Grafana, Otel) [Coming Soon]
├── bootstrap/          # App of Apps (ArgoCD root configurations)
└── infra/              # Core infrastructure definitions
```

## 🚀 Getting Started

1.  **Prerequisites**: A Kubernetes cluster (Kind, K3s, or standard).
2.  **Bootstrap**: Apply the ArgoCD Application manifest to kickstart the cluster state.

```bash
kubectl apply -f bootstrap/applications.yaml
```
