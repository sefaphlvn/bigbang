package scenario

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/crud/scenario/scenarios"
)

func (sc *AppHandler) SetScenario(ctx context.Context, scenario models.ScenarioBody, reqDetails models.RequestDetails) (interface{}, error) {
	templateMap, exists := scenarios.Scenarios[scenarios.Scenario(reqDetails.Metadata["scenario_id"])]
	if !exists {
		return nil, fmt.Errorf("scenario not found")
	}

	listenerUniqs := map[string]string{
		"UniqListenerNameID":    helper.GenerateUniqueId(6),
		"UniqFilterChainNameID": helper.GenerateUniqueId(6),
		"UniqFilterNameID":      helper.GenerateUniqueId(6),
	}

	successfulResources := []models.DBResourceClass{}
	response := map[string]interface{}{}

	for key, templateStr := range templateMap {
		if data, ok := scenario[key]; ok {

			templateData := map[string]interface{}{
				"Data":     data,
				"Version":  reqDetails.Version,
				"Project":  reqDetails.Project,
				"Listener": listenerUniqs,
			}

			tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(templateStr)
			if err != nil {
				sc.rollback(ctx, successfulResources, reqDetails)
				return nil, fmt.Errorf("template parse error: %w", err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, templateData); err != nil {
				sc.rollback(ctx, successfulResources, reqDetails)
				return nil, fmt.Errorf("template execute error: %w", err)
			}

			var jsonData interface{}
			if err := json.Unmarshal(buf.Bytes(), &jsonData); err != nil {
				sc.rollback(ctx, successfulResources, reqDetails)
				return nil, fmt.Errorf("failed to parse template output as JSON: %w", err)
			}

			data, err := decodeXdsExtension(jsonData)
			if err != nil {
				sc.rollback(ctx, successfulResources, reqDetails)
				return nil, fmt.Errorf("failed to decode XDS extension: %w", err)
			}

			xdsResponse, err := sc.SetResource(ctx, data, reqDetails)
			if err != nil {
				sc.rollback(ctx, successfulResources, reqDetails)
				return nil, fmt.Errorf("failed to save resource: %w", err)
			}

			if key == "listener" {
				response = xdsResponse
			}

			successfulResources = append(successfulResources, data)
		}
	}

	return response, nil
}

func (sc *AppHandler) SetResource(ctx context.Context, data models.DBResourceClass, reqDetails models.RequestDetails) (map[string]interface{}, error) {
	Gtype := data.GetGeneral().GType
	result := map[string]interface{}{}
	if helper.Contains([]string{"filters", "extensions"}, Gtype.CollectionString()) {
		response, err := sc.Extension.SetExtension(ctx, data, reqDetails)
		if err == nil {
			result = response.(map[string]interface{})
		}
		return result, err
	}

	response, err := sc.XDS.SetResource(ctx, data, reqDetails)
	if err == nil {
		result = response.(map[string]interface{})
	}
	return result, err
}

func (sc *AppHandler) rollback(ctx context.Context, resources []models.DBResourceClass, reqDetails models.RequestDetails) {
	retryList := make([]models.DBResourceClass, len(resources))
	copy(retryList, resources)

	for len(retryList) > 0 {
		log.Printf("Rollback attempt, resources left: %d", len(retryList))
		var failedResources []models.DBResourceClass

		for _, resource := range retryList {
			err := sc.DeleteResource(ctx, resource, reqDetails)
			if err != nil {
				if isDependencyError(err) {
					log.Printf("Resource %v has dependencies, will retry: %v", resource.GetGeneral().Name, err)
					failedResources = append(failedResources, resource)
				} else {
					log.Printf("Failed to delete resource %v: %v", resource.GetGeneral().Name, err)
				}
			}
		}

		if len(failedResources) == len(retryList) {
			log.Printf("Rollback stuck, no progress made in removing resources. Exiting...")
			for _, resource := range failedResources {
				log.Printf("Could not delete resource: %v", resource.GetGeneral().Name)
			}
			break
		}

		retryList = failedResources
	}

	if len(retryList) > 0 {
		log.Printf("Rollback completed with unresolved resources: %d", len(retryList))
		for _, resource := range retryList {
			log.Printf("Unresolved resource: %v", resource.GetGeneral().Name)
		}
	} else {
		log.Println("Rollback successfully completed with no remaining resources.")
	}
}

func isDependencyError(err error) bool {
	return strings.Contains(err.Error(), "Resource has dependencies")
}

func (sc *AppHandler) DeleteResource(ctx context.Context, resource models.DBResourceClass, reqDetails models.RequestDetails) error {
	general := resource.GetGeneral()
	reqDetails.Name = general.Name
	reqDetails.Collection = general.Collection
	reqDetails.GType = general.GType

	if helper.Contains([]string{"filters", "extensions"}, general.GType.CollectionString()) {
		_, err := sc.Extension.DelExtension(ctx, resource, reqDetails)
		return err
	} else {
		_, err := sc.XDS.DelResource(ctx, resource, reqDetails)
		return err
	}
}

func decodeXdsExtension(data interface{}) (models.DBResourceClass, error) {
	var resource models.DBResource

	resourceBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := json.Unmarshal(resourceBytes, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to DBResource: %w", err)
	}

	return &resource, nil
}
