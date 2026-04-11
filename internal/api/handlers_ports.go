package api

import "net/http"

func (s *Server) handleTopPorts(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	limit := parseInt(r.URL.Query().Get("limit"), 10)
	data, err := s.store.TopPorts(r.Context(), f, limit)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}
