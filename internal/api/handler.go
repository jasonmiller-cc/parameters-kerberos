// Package api wires the HTTP routes for the parameters-kerberos service.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jasonmiller-cc/parameters-core/pkg/response"
	"github.com/jasonmiller-cc/parameters-kerberos/internal/service"
)

// Handler exposes Kerberos KDC operations over HTTP.
type Handler struct {
	svc *service.KerberosService
}

// New creates an API Handler backed by svc.
func New(svc *service.KerberosService) *Handler {
	return &Handler{svc: svc}
}

// Register mounts all routes onto mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/principals", h.principals)
	mux.HandleFunc("/api/v1/principals/", h.principalByName)
	mux.HandleFunc("/api/v1/policies", h.policies)
	mux.HandleFunc("/api/v1/policies/", h.policyByName)
	mux.HandleFunc("/api/v1/realm", h.realm)
}

// ---------------------------------------------------------------------------
// Principal handlers
// ---------------------------------------------------------------------------

func (h *Handler) principals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := h.svc.ListPrincipals(r.Context())
		if err != nil {
			response.Err(w, err)
			return
		}
		response.OK(w, list)

	case http.MethodPost:
		var req service.AddPrincipalRequest
		if err := decodeJSON(r, &req); err != nil {
			response.ErrStatus(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if err := h.svc.AddPrincipal(r.Context(), req); err != nil {
			response.Err(w, err)
			return
		}
		response.Created(w, map[string]string{"principal": req.Principal})

	default:
		methodNotAllowed(w)
	}
}

// principalByName handles:
//
//	GET    /api/v1/principals/{principal}
//	DELETE /api/v1/principals/{principal}
//	POST   /api/v1/principals/{principal}/password
//	POST   /api/v1/principals/{principal}/keytab
func (h *Handler) principalByName(w http.ResponseWriter, r *http.Request) {
	// Strip prefix and split path segments.
	tail := strings.TrimPrefix(r.URL.Path, "/api/v1/principals/")
	parts := strings.SplitN(tail, "/", 2)
	name := parts[0]

	if name == "" {
		response.ErrStatus(w, http.StatusBadRequest, "bad_request", "principal name is required")
		return
	}

	// Sub-resources: /password, /keytab
	if len(parts) == 2 {
		switch parts[1] {
		case "password":
			if r.Method != http.MethodPost {
				methodNotAllowed(w)
				return
			}
			var req service.ChangePWRequest
			if err := decodeJSON(r, &req); err != nil {
				response.ErrStatus(w, http.StatusBadRequest, "bad_request", err.Error())
				return
			}
			if err := h.svc.ChangePassword(r.Context(), name, req.Password); err != nil {
				response.Err(w, err)
				return
			}
			response.NoContent(w)

		case "keytab":
			if r.Method != http.MethodPost {
				methodNotAllowed(w)
				return
			}
			data, err := h.svc.GenerateKeytab(r.Context(), name)
			if err != nil {
				response.Err(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", `attachment; filename="`+name+`.keytab"`)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)

		default:
			response.ErrStatus(w, http.StatusNotFound, "not_found", "unknown sub-resource")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := h.svc.GetPrincipal(r.Context(), name)
		if err != nil {
			response.Err(w, err)
			return
		}
		response.OK(w, p)

	case http.MethodDelete:
		if err := h.svc.DeletePrincipal(r.Context(), name); err != nil {
			response.Err(w, err)
			return
		}
		response.NoContent(w)

	default:
		methodNotAllowed(w)
	}
}

// ---------------------------------------------------------------------------
// Policy handlers
// ---------------------------------------------------------------------------

func (h *Handler) policies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list, err := h.svc.ListPolicies(r.Context())
		if err != nil {
			response.Err(w, err)
			return
		}
		response.OK(w, list)

	case http.MethodPost:
		var p service.Policy
		if err := decodeJSON(r, &p); err != nil {
			response.ErrStatus(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if err := h.svc.CreatePolicy(r.Context(), p); err != nil {
			response.Err(w, err)
			return
		}
		response.Created(w, p)

	default:
		methodNotAllowed(w)
	}
}

// policyByName handles:
//
//	GET    /api/v1/policies/{name}
//	DELETE /api/v1/policies/{name}
func (h *Handler) policyByName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/v1/policies/")
	if name == "" {
		response.ErrStatus(w, http.StatusBadRequest, "bad_request", "policy name is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := h.svc.GetPolicy(r.Context(), name)
		if err != nil {
			response.Err(w, err)
			return
		}
		response.OK(w, p)

	case http.MethodDelete:
		if err := h.svc.DeletePolicy(r.Context(), name); err != nil {
			response.Err(w, err)
			return
		}
		response.NoContent(w)

	default:
		methodNotAllowed(w)
	}
}

// ---------------------------------------------------------------------------
// Realm handler
// ---------------------------------------------------------------------------

func (h *Handler) realm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	info, err := h.svc.RealmInfo(r.Context())
	if err != nil {
		response.Err(w, err)
		return
	}
	response.OK(w, info)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func methodNotAllowed(w http.ResponseWriter) {
	response.ErrStatus(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
