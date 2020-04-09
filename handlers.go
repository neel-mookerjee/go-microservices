package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\nAPI instructions here")
}

func (h *HandlerWrapper) GetMetricsDownsamplingForPreview(w http.ResponseWriter, r *http.Request) {
	res, err := h.Usecase.GetNextPreview()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if res.QueryId == "" {
		respondJSON(w, http.StatusNoContent, "no object found")
		return
	}

	respondJSON(w, http.StatusOK, res)
}

func (h *HandlerWrapper) GetMetricsDownsamplingPendingQuery(w http.ResponseWriter, r *http.Request) {
	res, err := h.Usecase.GetNextQuery()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if res.QueryId == "" {
		respondJSON(w, http.StatusNoContent, "no object found")
		return
	}

	respondJSON(w, http.StatusOK, res)
}

func (h *HandlerWrapper) GetMetricsDownsamplingItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["queryid"]

	res, err := h.Usecase.GetDownsamplingItem(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if res.QueryId == "" {
		respondJSON(w, http.StatusNoContent, "no object found")
		return
	}

	respondJSON(w, http.StatusOK, res)
}

func (h *HandlerWrapper) DeployMetricsDownsamplingPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["queryid"]

	err := h.Usecase.DeployPreviewItem(id)
	if err != nil {
		if err.Error() == "object not found" {
			respondError(w, http.StatusNotFound, err)
		} else if err.Error() == "status changed" {
			respondError(w, http.StatusConflict, err)
		} else {
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respondJSON(w, http.StatusAccepted, "Success")
}

func (h *HandlerWrapper) GetMetricsStacks(w http.ResponseWriter, r *http.Request) {
	res, err := h.Usecase.GetMetricsStacks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if len(res) == 0 {
		respondJSON(w, http.StatusNoContent, "no object found")
		return
	}

	respondJSON(w, http.StatusOK, res)
}
