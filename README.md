# 📦 hexttp – Standardized HTTP Response & Metadata Handling for Go APIs

<span style="color: pink"> Experimental Version</span>
---

## 🛠 Installation
```bash
go get github.com/tamvanum/go-hexttp
```
---
`hexttp` is a minimalistic HTTP response helper for Go web applications. It standardizes the way responses are returned from handlers, enriching them with structured metadata and logging. The library focuses on pragmatic conventions and separation of concerns, allowing you to build **clean, consistent, and traceable APIs**.

---

## 🚀 Purpose

The goal of `hexttp` is to **centralize the structure of HTTP responses** while also handling **logging levels** and **metadata propagation** (e.g. request ID, duration) with minimal boilerplate.

This package helps you:

- Return API responses in a **consistent format**
- Attach and propagate **metadata** like `X-Request-ID` and `duration`
- Handle **logging** in a clean and declarative way (`warn`, `error`)
- Make your handlers return clean responses

---

## 📦 Features

- Unified response wrapper: `*HTTPResponse`
- Built-in helpers for common HTTP responses
- Structured log level tagging
- Context-based request metadata
- Chi-compatible middleware: `MetaDataCollector`
- Wrapper function `Make()` to use your handlers cleanly

---

## ⚙️ Usage

### 1. Register middleware in your router

Apply `MetaDataCollector` to enrich each request with:

- A generated or forwarded `X-Request-ID`
- A timestamp for calculating duration

```go
import (
    "github.com/go-chi/chi/v5"
	"github.com/tamvanum/go-hexttp"
)

rt := chi.NewRouter()

rt.Use(middleware.Logger)      // Optional: Chi's logger
rt.Use(middleware.Recoverer)   // Optional: Panic recovery
rt.Use(hexttp.MetaDataCollector)
```




## 2. Wrap your routes with `hexttp.Make`

All route handlers should return a `*hexttp.HTTPResponse` instead of directly writing to the `http.ResponseWriter`.

```go
rt.Get("/", hexttp.Make(hd.GetAll))
rt.Get("/{id}", hexttp.Make(hd.GetByID))
rt.Post("/", hexttp.Make(hd.Create))
rt.Patch("/{id}", hexttp.Make(hd.Update))
rt.Delete("/{id}", hexttp.Make(hd.Delete))
```


## 3. Define your handlers to return `*HTTPResponse`

This allows `hexttp` to:

- Automatically set response headers  
- Encode the response body  
- Log based on severity  

```go
func (h *SectionHTTPHandler) Create(w http.ResponseWriter, r *http.Request) *hexttp.HTTPResponse {
    var req sectionreq.SectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return hexttp.InvalidJSON()
	}

	if errors := req.Validate(); len(errors) != 0 {
        return hexttp.InvalidRequestData(errors)
	}

	section, err := h.sv.Create(&req)
	if err != nil {
        if errors.Is(err, sectiondom.ErrSectionNumberAlreadyExists) {
            return hexttp.AlreadyExist(err)
		}
		return hexttp.InternalError("unexpected error")
	}

	res := sectionres.SectionModelToDTO(section)
	return hexttp.Created(res)
}
```


## 🧱 Available Response Helpers

| Function                    | Status Code | Log Level | Description                     |
|-----------------------------|-------------|-----------|---------------------------------|
| `OK(data)`                  | 200         | None      | Standard success with data      |
| `Created(data)`             | 201         | None      | Resource created                |
| `Updated(data)`             | 202         | None      | Resource updated                |
| `NoContent()`               | 204         | None      | Successful with no body         |
| `InvalidJSON()`             | 400         | Warn      | Request body not valid JSON     |
| `InvalidRequestData(errs)`  | 400         | Warn      | Domain or validation error      |
| `InvalidID()`               | 400         | Warn      | Invalid path or ID format       |
| `Unauthorized()`            | 401         | Warn      | Authentication required         |
| `Forbidden()`               | 403         | Warn      | Access not allowed              |
| `NotFound(err)`             | 404         | Warn      | Entity not found                |
| `AlreadyExist(err)`         | 409         | Warn      | Conflict: already exists        |
| `InternalError(msg)`        | 500         | Error     | Unexpected internal failure     |

---

## 🧩 Metadata Collected per Request

The `MetaDataCollector` middleware injects into context:

| Metadata Key   | Header/Context Value | Description                        |
|----------------|----------------------|------------------------------------|
| `requestID`    | `X-Request-ID`       | Propagated or auto-generated UUID  |
| `requestStart` | context value        | Timestamp for duration calculation |

You can retrieve them manually:

```go
reqID := r.Context().Value(hexttp.RequestIDKey).(string)
start := r.Context().Value(hexttp.RequestStartKey).(time.Time)
```

Or use helper functions (recommended):

```go
reqID := hexttp.GetRequestID(r.Context())
start, ok := hexttp.GetRequestStart(r.Context())
```

---

## ✨ Example Log Output

```text
level=ERROR msg="HTTP API error" status=400 path="/api/sections"
request_id="71f0...c934" duration="34ms" body="invalid id"
```