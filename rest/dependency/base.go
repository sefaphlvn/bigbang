package dependency

import (
	"fmt"
)

var (
	visitedUpstream   = make(map[string]bool)
	visitedDownstream = make(map[string]bool)
)

func (h *AppHandler) ProcessResource(activeResource Depend) {
	visitedUpstream = make(map[string]bool)
	h.ProcessUpstream(activeResource)

	visitedDownstream = make(map[string]bool)
	h.ProcessDownstream(activeResource)
}

func generateUniqueKey(resource Depend) string {
	return fmt.Sprintf("%s_%s_%s_%s", resource.Name, resource.Gtype, resource.Collection, resource.Project)
}

func (h *AppHandler) ProcessUpstream(activeResource Depend) {
	uniqueKey := generateUniqueKey(activeResource)
	if visitedUpstream[uniqueKey] {
		return
	}

	visitedUpstream[uniqueKey] = true
	node, upstreams := h.CallUpstreamFunction(activeResource)
	if node.ID != "" && node.Name != "" && node.Gtype != "" {
		h.AddNode(node)
		activeResource.First = false
	} else {
		h.Context.Logger.Infof("Node is missing required fields, not adding: %+v\n", node)
	}

	for _, up := range upstreams {
		if up.ID != "" && up.Name != "" && up.Gtype != "" {
			h.AddNodeAndEdge(node, up, true)
			h.ProcessUpstream(up) // Sadece upstream keşfi yap
		} else {
			h.Context.Logger.Infof("Upstream is missing required fields, not adding: %+v\n", up)
		}
	}
}

func (h *AppHandler) ProcessDownstream(activeResource Depend) {
	uniqueKey := generateUniqueKey(activeResource)
	if visitedDownstream[uniqueKey] {
		return
	}

	visitedDownstream[uniqueKey] = true

	node, downstreams := h.CallDownstreamFunction(activeResource)
	if h.isNodeValid(node) {
		h.AddNode(node)
		activeResource.First = false
	} else {
		h.Context.Logger.Infof("Node is missing required fields, not adding: %+v\n", node)
	}

	for _, down := range downstreams {
		if h.isValidDownstream(node, down) {
			h.AddNodeAndEdge(node, down, false)
			h.ProcessDownstream(down)
		} else {
			h.Context.Logger.Infof("Downstream is missing required fields, not directly connected, or from incorrect source, not adding: %+v\n", down)
		}
	}
}

// Yardımcı fonksiyonlar.
func (h *AppHandler) isNodeValid(node Node) bool {
	return node.ID != "" && node.Name != "" && node.Gtype != ""
}

func (h *AppHandler) isValidDownstream(node Node, down Depend) bool {
	return down.ID != "" && down.Name != "" && down.Gtype != "" &&
		down.Direction == "downstream" && down.Source == node.ID
}
