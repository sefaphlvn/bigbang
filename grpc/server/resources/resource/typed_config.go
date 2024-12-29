package resource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/sjson"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
)

func (ar *AllResources) GetTypedConfigs(ctx context.Context, paths []models.TypedConfigPath, jsonData interface{}, context *db.AppContext) (interface{}, error) {
	jsonStringStr, err := helper.MarshalJSON(jsonData, context.Logger)
	if err != nil {
		return jsonData, err
	}

	for _, pathd := range paths {
		if err := ar.processTypedConfigPath(ctx, pathd, &jsonStringStr, context); err != nil {
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

func (ar *AllResources) processTypedConfigPath(ctx context.Context, pathd models.TypedConfigPath, jsonStringStr *string, context *db.AppContext) error {
	_, typedConfigsMap := resources.ProcessTypedConfigs(*jsonStringStr, pathd, context.Logger)

	for path, tempTypedConfig := range typedConfigsMap {
		conf, err := resources.GetResourceNGeneral(ctx, context, tempTypedConfig.Collection, tempTypedConfig.Name, ar.Project)
		if err != nil {
			context.Logger.Warnf("Error getting resource from DB: %v", err)
			return err
		}

		resource := conf.GetResource()

		typedConfigJSON, err := json.Marshal(resource)
		if err != nil {
			context.Logger.Warnf("Error marshaling typed config: %v", err)
			return err
		}
		typedConfigStr := string(typedConfigJSON)

		ar.processUpstreamPaths(ctx, tempTypedConfig.Gtype.UpstreamPaths(), &typedConfigStr, tempTypedConfig.ParentName, context, context.Logger)

		typedConfig, err := decodeTypedConfig([]byte(typedConfigStr), tempTypedConfig.Gtype)
		if err != nil {
			context.Logger.Warnf("Error decoding typed config: %v", err)
			return err
		}

		if err := ar.updateJSONConfig(jsonStringStr, path, typedConfig, pathd.IsPerTypedConfig, tempTypedConfig); err != nil {
			context.Logger.Warnf("Error updating JSON config: %v", err)
			return err
		}
	}

	return nil
}

func (ar *AllResources) updateJSONConfig(jsonStringStr *string, path string, typedConfig *anypb.Any, isPerTypedConfig bool, tempTypedConfig *models.TypedConfig) error {
	var config interface{}
	var err error

	if isPerTypedConfig && tempTypedConfig.Disabled {
		config = map[string]interface{}{
			"@type":    "type.googleapis.com/envoy.config.route.v3.FilterConfig",
			"disabled": true,
		}
	} else {
		anyJSON, err := protojson.Marshal(typedConfig)
		if err != nil {
			return fmt.Errorf("error marshaling any typed config: %w", err)
		}

		if err := json.Unmarshal(anyJSON, &config); err != nil {
			return fmt.Errorf("error marshaling any typed config: %w", err)
		}
	}

	if *jsonStringStr, err = sjson.Set(*jsonStringStr, path, config); err != nil {
		return fmt.Errorf("error setting new config value with sjson.Set: %w", err)
	}

	return nil
}

func decodeTypedConfig(typedConfigJSON []byte, gtype models.GTypes) (*anypb.Any, error) {
	msg := gtype.ProtoMessage()
	if err := helper.Unmarshaler.Unmarshal(typedConfigJSON, msg); err != nil {
		return nil, fmt.Errorf("typed_config not resolved: %w", err)
	}

	return anypb.New(msg)
}
