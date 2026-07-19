# Rencana Pengembangan NAEOS

Dokumen ini menjabarkan rencana pengembangan teknis jangka pendek dan menengah,
melengkapi [ROADMAP.md](ROADMAP.md) dan [NAEOS-GOV-007](governance/NAEOS-GOV-007.md)
dengan item kerja konkret.

## Status Saat Ini

- **Versi:** v1.5.0
- **Fase Roadmap:** Fase 1-2 ✅ | Fase 3 (Referensi Implementasi) ⏳ | Fase 4 (Ekosistem) ⏳
- **Milestone Berikutnya:** v2.0.0 — Ecosystem Platform

---

## Sprint 1: Stabilitas & Kebersihan Kode

**Tujuan:** Membersihkan technical debt, bug, dan memperkuat fondasi sebelum fitur baru.

### 1.1 Perbaikan Konteks & Timeout
| Item | Prioritas | Est. |
|------|-----------|------|
| Ganti `context.Background()` dengan timeout context di `internal/ai/llm.go:234,290` | Medium | 1 hr |
| Ganti `context.Background()` dengan timeout context di `internal/telemetry/telemetry.go:175` | Medium | 0.5 hr |
| Ganti `context.Background()` dengan timeout context di `internal/marketplace/remote.go:53,202` | Medium | 1 hr |
| Ganti `context.Background()` dengan `exec.CommandContext` + timeout di `cmd/naeos/doctor_cmd.go` (8 lokasi) | Medium | 1 hr |
| Ganti `context.Background()` dengan timeout context di database adapter (`mysql/sqlite/postgres/redis`) | Medium | 2 hr |

### 1.2 Error Handling
| Item | Prioritas | Est. |
|------|-----------|------|
| Perbaiki `save()` di `internal/workflow/workflow.go:576-591` — jangan silent ignore errors | High | 1 hr |
| Perbaiki `load()` di `internal/workflow/workflow.go:593-616` — jangan silent ignore read errors | High | 0.5 hr |
| Perbaiki hanya first error yang ditangkap dari parallel goroutines di `internal/workflow/workflow.go:351-356` | High | 1 hr |
| HTTP response body di-drain sebelum close di `cmd/naeos/doctor_cmd.go:314` | Medium | 0.5 hr |

### 1.3 Locking & Concurrency
| Item | Prioritas | Est. |
|------|-----------|------|
| Ganti `Lock()` → `RLock()` untuk operasi read-only di `internal/database/database.go:139-149` | Low | 0.5 hr |
| Fix gap lock di `internal/cache/cache.go:367-386` (RLock dilepas lalu Lock lagi) | Medium | 1 hr |
| Tambah nil check untuk `ConnectionPool.Next()` caller di `broker/broker.go:677` | Medium | 0.5 hr |

### 1.4 Security
| Item | Prioritas | Est. |
|------|-----------|------|
| Perbaiki `rand.Read` error diabaikan di `internal/auth/auth.go:510` | Medium | 0.5 hr |
| Validasi path traversal di file operations (gosec G304 audit) | Low | 1 hr |

### 1.5 Testing
| Item | Prioritas | Est. |
|------|-----------|------|
| Tambah test untuk `internal/rollback` edge cases (empty tar, corrupt file) | Medium | 2 hr |
| Tambah test untuk `internal/broker` connection pool edge cases | Medium | 2 hr |
| Isolasi test `cmd/naeos` yang flaky karena shared state | High | 3 hr |
| Benchmark regresi untuk pipeline parsing | Low | 2 hr |

**Total Estimasi Sprint 1:** ~17 jam

---

## Sprint 2: Observability & Operability

**Tujuan:** Meningkatkan kemampuan observasi, debugging, dan operasional.

| Item | Prioritas | Est. |
|------|-----------|------|
| Dashboard: real-time pipeline progress via WebSocket | High | 4 hr |
| Metrics: tambah custom metrics untuk pipeline stages (histogram latency) | Medium | 3 hr |
| Logging: structured logging untuk semua error path (current: banyak error pakai `fmt.Errorf` saja) | Medium | 3 hr |
| CLI: `naeos doctor` — tambah check untuk database connectivity, broker health | High | 2 hr |
| CLI: `naeos status` — tampilkan versi, uptime, pipeline cache stats | Low | 2 hr |
| Tracing: propagasi trace context ke subprocess (exec.CommandContext) | Low | 3 hr |

**Total Estimasi Sprint 2:** ~17 jam

---

## Sprint 3: Fitur v2.0 — Ecosystem Platform

**Tujuan:** Mewujudkan target v2.0 sebagai Ecosystem Platform.

### 3.1 Dashboard UI
| Item | Prioritas | Est. |
|------|-----------|------|
| Dashboard: spec visualizer (tree view + dependency graph) | High | 8 hr |
| Dashboard: pipeline run history with filter/search | High | 4 hr |
| Dashboard: real-time logs stream dari WebSocket | Medium | 4 hr |
| Dashboard: plugin management UI (list, install, uninstall) | Medium | 3 hr |

### 3.2 Profile & Marketplace
| Item | Prioritas | Est. |
|------|-----------|------|
| Profile Registry: publish/subscribe profil via remote registry | High | 6 hr |
| Profile: tambah 5 profile industri (SaaS, FinTech, Health, Edu, AI Agent) | Medium | 5 hr |
| Marketplace: dependency resolution untuk plugin | Medium | 4 hr |
| Marketplace: versioning + rollback plugin | Low | 3 hr |

### 3.3 Distributed Builds
| Item | Prioritas | Est. |
|------|-----------|------|
| Distributed: agent registration + heartbeat | High | 4 hr |
| Distributed: task queuing dengan priority | High | 3 hr |
| Distributed: result streaming (partial results sebelum selesai) | Medium | 3 hr |
| Distributed: CLI `naeos build --distributed` | Medium | 2 hr |

### 3.4 AI & Compiler
| Item | Prioritas | Est. |
|------|-----------|------|
| AI: streaming response dari LLM (SSE) | High | 3 hr |
| AI: context window management (auto-truncate sesuai batas model) | Medium | 3 hr |
| Compiler: tambah output adapter untuk Windsurf | Medium | 2 hr |
| Compiler: custom prompt builder (user-defined templates) | Low | 3 hr |

### 3.5 NEIR & Specification
| Item | Prioritas | Est. |
|------|-----------|------|
| Spec: `$import{}` function untuk modular spec fragments | High | 4 hr |
| Spec: schema registry (versioned JSON Schema untuk spec) | Medium | 5 hr |
| NEIR: visual diff (graphical perbandingan 2 versi arsitektur) | Low | 4 hr |

**Total Estimasi Sprint 3:** ~70 jam

---

## Sprint 4: Enterprise Readiness

**Tujuan:** Mempersiapkan fondasi untuk adopsi enterprise.

| Item | Prioritas | Est. |
|------|-----------|------|
| Role-Based Access Control (RBAC) untuk API server | High | 6 hr |
| Multi-tenant workspace isolation | High | 8 hr |
| Audit trail lengkap (siapa, apa, kapan untuk semua operasi) | Medium | 4 hr |
| Rate limiting per tenant (bukan per API key saja) | Medium | 3 hr |
| Encryption at rest untuk sensitive config tersimpan | Low | 3 hr |
| Compliance reporting (export audit log ke JSON/CSV) | Low | 2 hr |

**Total Estimasi Sprint 4:** ~26 jam

---

## Ringkasan Timeline

| Sprint | Fokus | Estimated Hours | Target Rilis |
|--------|------|----------------|--------------|
| Sprint 1 | Stabilitas & Kebersihan | ~17 jam | v1.6.0 |
| Sprint 2 | Observability | ~17 jam | v1.7.0 |
| Sprint 3 | Ecosystem Platform | ~70 jam | v2.0.0 |
| Sprint 4 | Enterprise Readiness | ~26 jam | v2.1.0 |

---

## Prinsip Pengembangan

1. **Test first** — setiap perubahan harus memiliki test yang lulus dengan `-race`
2. **No regressions** — semua test existing harus tetap hijau
3. **Backward compatibility** — tidak ada breaking change tanpa major version
4. **Dokumentasi** — tiap fitur baru harus disertai NES spec update
5. **Code review** — minimal 1 reviewer untuk tiap PR

---

## Referensi

- [ROADMAP.md](ROADMAP.md) — Roadmap utama
- [NAEOS-GOV-007](governance/NAEOS-GOV-007.md) — Governance roadmap
- [NAEOS-GOV-008](governance/NAEOS-GOV-008.md) — Versioning policy
- [CONTRIBUTING.md](CONTRIBUTING.md) — Panduan kontribusi
