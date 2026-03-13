# Multi-file error scenarios — by genre

Goal: trigger panics where the **stack trace spans multiple files**, with **different genres** (structural patterns), not just one pattern repeated.

---

## Genre 1: Layered sync chain

Linear call chain across 2–3 files; panic in the deepest layer.

### 1a — Order pipeline (nil deref in formatter)

**Route:** `GET /error/multi-file/order`

| Layer | File | Role |
|-------|------|------|
| 1 | `handlers/multi_file_order.go` | Builds order, calls pipeline |
| 2 | `pipeline/order.go` | `ProcessOrder()` → `formatter.BuildInvoice(order)` |
| 3 | `pipeline/formatter.go` | **Panic:** `order.ShippingAddress.Street` (ShippingAddress nil) |

**Stack:** formatter.go → order.go → multi_file_order.go

---

### 1b — Config / DSN chain (type assertion in env)

**Route:** `GET /error/multi-file/config`

| Layer | File | Role |
|-------|------|------|
| 1 | `handlers/multi_file_config.go` | Calls config service |
| 2 | `configsvc/service.go` | `GetDatabaseDSN()` → `env.Expand(config["dsn"])` |
| 3 | `configsvc/env.go` | **Panic:** `Expand(v)` does `v.(string)`, v is map |

**Stack:** env.go → service.go → multi_file_config.go

---

### 1c — Cache + repo (nil DB in repo)

**Route:** `GET /error/multi-file/cache`

| Layer | File | Role |
|-------|------|------|
| 1 | `handlers/multi_file_cache.go` | Calls cache service |
| 2 | `cachesvc/cache.go` | Cache miss → `repo.FindByID(id)` |
| 3 | `cachesvc/repo.go` | **Panic:** `r.db.QueryRow(...)`, r.db nil |

**Stack:** repo.go → cache.go → multi_file_cache.go

---

## Genre 2: Interface boundary

Handler calls an **interface**; the **implementation** lives in another package and panics there. Stack crosses the abstraction boundary.

### 2 — User fetcher (impl panic)

**Route:** `GET /error/multi-file/interface`

| Layer | File | Role |
|-------|------|------|
| 1 | `handlers/multi_file_interface.go` | Calls `usersvc.GetUser(id)` |
| 2 | `usersvc/svc.go` | `Service.GetUser` delegates to `Fetcher.FetchUser(id)` |
| 3 | `userfetcher/impl.go` | **Panic:** `Impl.FetchUser` derefs nil `*User` on cache miss |

**Stack:** userfetcher/impl.go → usersvc/svc.go → handlers/multi_file_interface.go

**Genre:** Panic is inside the *implementation* of an interface (different package), not in a “layer” that the handler called directly.

---

## Genre 3: Callback / visitor

Handler passes a **callback** (or visitor) into another package; the callback runs in the **caller’s** code path and panics there. Stack: handler (callback) → processor.

### 3 — Processor invokes callback that panics

**Route:** `GET /error/multi-file/callback`

| Layer | File | Role |
|-------|------|------|
| 1 | `handlers/multi_file_callback.go` | Defines callback, calls `processor.Process(items, callback)` |
| 2 | `processor/process.go` | `Process` loops, calls `invoke(it, fn)` |
| 3 | `processor/invoke.go` | `invoke(it, fn)` → `fn(it)` |
| (panic) | `handlers/multi_file_callback.go` | Callback does `it.Child.Name`; **Child is nil** → panic in handler’s closure |

**Stack:** multi_file_callback.go (callback) → processor/invoke.go → processor/process.go

**Genre:** Panic is in *caller-provided code* (the callback), invoked by another package — different from “panic in callee” or “panic in impl”.

---

## Summary

| Genre | Route(s) | What’s different |
|-------|----------|-------------------|
| **Layered** | `/error/multi-file/order`, `/config`, `/cache` | Straight A→B→C; panic in deepest layer. |
| **Interface** | `/error/multi-file/interface` | Panic in interface *implementation* (other package). |
| **Callback** | `/error/multi-file/callback` | Panic in *caller’s callback* invoked by another package. |

---

## Routes

- `GET /error/multi-file/order` — layered (formatter nil)
- `GET /error/multi-file/config` — layered (env type assertion)
- `GET /error/multi-file/cache` — layered (repo nil db)
- `GET /error/multi-file/interface` — interface (userfetcher impl)
- `GET /error/multi-file/callback` — callback (handler callback nil Child)

## Layout

```
handlers/
  multi_file_order.go    # layered
  multi_file_config.go   # layered
  multi_file_cache.go    # layered
  multi_file_interface.go # genre: interface
  multi_file_callback.go  # genre: callback
pipeline/     # layered
configsvc/    # layered
cachesvc/     # layered
usersvc/      # interface (defines Fetcher)
userfetcher/  # interface (impl, panic here)
processor/    # callback (process.go, invoke.go, item.go)
```
