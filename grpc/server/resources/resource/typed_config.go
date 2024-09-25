package resource

import (
	"encoding/json"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/tidwall/sjson"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func (ar *AllResources) GetTypedConfigs(paths []models.TypedConfigPath, jsonData interface{}, context *db.AppContext) (interface{}, error) {
	jsonStringStr, err := helper.MarshalJSON(jsonData, context.Logger)
	if err != nil {
		return jsonData, err
	}

	for _, pathd := range paths {
		if err := ar.processTypedConfigPath(pathd, &jsonStringStr, context); err != nil {
			context.Logger.Debugf("Error processing typed config path: %v", err)
		}
	}

	var updatedJSONData interface{}
	if err := json.Unmarshal([]byte(jsonStringStr), &updatedJSONData); err != nil {
		context.Logger.Errorf("Error unmarshalling updated JSON: %v", err)
		return nil, err
	}

	return updatedJSONData, nil
}

func (ar *AllResources) processTypedConfigPath(pathd models.TypedConfigPath, jsonStringStr *string, context *db.AppContext) error {
	_, typedConfigsMap := resources.ProcessTypedConfigs(*jsonStringStr, pathd, context.Logger)

	for path, tempTypedConfig := range typedConfigsMap {
		conf, err := resources.GetResourceNGeneral(context, tempTypedConfig.Collection, tempTypedConfig.Name, ar.Project)
		if err != nil {
			context.Logger.Warnf("Error getting resource from DB: %v", err)
			continue
		}

		resource := conf.GetResource()
		ar.DetectResources(pathd.Kind, conf, context)

		typedConfigJSON, err := json.Marshal(resource)
		if err != nil {
			context.Logger.Warnf("Error marshalling typed config: %v", err)
			continue
		}

		typedConfig, err := decodeTypedConfig(typedConfigJSON, tempTypedConfig.Gtype)
		if err != nil {
			context.Logger.Warnf("Error decoding typed config: %v", err)
			continue
		}

		if err := ar.updateJSONConfig(jsonStringStr, path, typedConfig); err != nil {
			context.Logger.Warnf("Error updating JSON config: %v", err)
		}
	}

	return nil
}

func (ar *AllResources) updateJSONConfig(jsonStringStr *string, path string, typedConfig *anypb.Any) error {
	anyJSON, err := protojson.Marshal(typedConfig)
	if err != nil {
		return fmt.Errorf("error marshalling any typed config: %w", err)
	}

	var typedConfigMap map[string]interface{}
	if err := json.Unmarshal(anyJSON, &typedConfigMap); err != nil {
		return fmt.Errorf("error unmarshalling any typed config: %w", err)
	}

	if *jsonStringStr, err = sjson.Set(*jsonStringStr, path, typedConfigMap); err != nil {
		return fmt.Errorf("error setting new config value with sjson.Set: %w", err)
	}

	return nil
}

func decodeTypedConfig(typedConfigJSON []byte, gtype models.GTypes) (*anypb.Any, error) {
	msg := gtype.ProtoMessage()

	if err := protojson.Unmarshal(typedConfigJSON, msg); err != nil {
		return nil, fmt.Errorf("typed_config not resolved: %w", err)
	}

	return anypb.New(msg)
}

func (ar *AllResources) DetectResources(pathKind string, conf *models.DBResource, context *db.AppContext) {
	if pathKind == "downstream_tls" {
		ar.DecodeDownstreamTLS(conf, context)
	}
}
