package xds

import (
	"context"
	"errors"
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sefaphlvn/bigbang/rest/crud/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectorFunc func(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails, version string) (models.DBResourceClass, error)

type BootstrapCollector struct {
	collectors map[string]CollectorFunc
}

func (xds *AppHandler) NewBootstrapCollector() *BootstrapCollector {
	return &BootstrapCollector{
		collectors: map[string]CollectorFunc{
			"clusters":    xds.collectBootstrapClusters,
			"access_logs": xds.collectAccessLoggers,
		},
	}
}

func (xds *AppHandler) DownloadBootstrap(ctx context.Context, requestDetails models.RequestDetails) (interface{}, error) {
	resource := &models.DBResource{}
	collection := xds.Context.Client.Collection(requestDetails.Collection)
	filter, err := common.AddResourceIDFilter(requestDetails, bson.M{"general.name": requestDetails.Name})
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	filterWithRestriction := common.AddUserFilter(requestDetails, filter)
	result := collection.FindOne(ctx, filterWithRestriction)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + requestDetails.Name + ")")
		}
		return nil, errstr.ErrUnknownDBError
	}

	if err := result.Decode(resource); err != nil {
		return nil, err
	}

	collector := xds.NewBootstrapCollector()
	return collector.CollectAll(ctx, resource, requestDetails)
}

func (bc *BootstrapCollector) CollectAll(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails) (models.DBResourceClass, error) {
	var err error
	bootstrap := resource
	version := resource.GetGeneral().Version

	for name, collector := range bc.collectors {
		if shouldSkip, err := bc.shouldSkipCollection(bootstrap, name); err != nil {
			return nil, fmt.Errorf("checking %s collection: %w", name, err)
		} else if shouldSkip {
			continue
		}

		bootstrap, err = collector(ctx, bootstrap, requestDetails, version)
		if err != nil {
			return nil, fmt.Errorf("collecting %s: %w", name, err)
		}
	}

	return bootstrap, nil
}

func (bc *BootstrapCollector) shouldSkipCollection(resource models.DBResourceClass, collectorName string) (bool, error) {
	bootstrapMap, ok := resource.GetResource().(primitive.M)
	if !ok {
		return false, fmt.Errorf("invalid bootstrap format")
	}

	switch collectorName {
	case "access_logs":
		admin, ok := bootstrapMap["admin"].(primitive.M)
		if !ok {
			return true, nil
		}
		_, hasAccessLog := admin["access_log"].(primitive.A)
		return !hasAccessLog, nil

	case "clusters":
		staticResources, ok := bootstrapMap["static_resources"].(primitive.M)
		if !ok {
			return true, nil
		}
		_, hasClusters := staticResources["clusters"].(primitive.A)
		return !hasClusters, nil

	default:
		return false, fmt.Errorf("unknown collector: %s", collectorName)
	}
}

func (xds *AppHandler) collectBootstrapClusters(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails, version string) (models.DBResourceClass, error) {
	bootstrap := resource.GetResource()
	bootstrapMap, ok := bootstrap.(primitive.M)
	if !ok {
		return nil, fmt.Errorf("failed to parse bootstrap as primitive.M, got type: %T", bootstrap)
	}

	staticResources, ok := bootstrapMap["static_resources"].(primitive.M)
	if !ok {
		return nil, fmt.Errorf("'static_resources' key not found or invalid")
	}

	clusters, ok := staticResources["clusters"].(primitive.A)
	if !ok {
		return nil, fmt.Errorf("'clusters' key not found or invalid")
	}

	var clusterNames []string
	for _, cluster := range clusters {
		clusterMap, ok := cluster.(primitive.M)
		if !ok {
			continue
		}

		if name, ok := clusterMap["name"].(string); ok {
			clusterNames = append(clusterNames, name)
		}
	}

	clusters, err := xds.GetNonEdsClusters(ctx, clusterNames, requestDetails, version)
	if err != nil {
		return nil, err
	}

	resource.SetBootstrapClusters(clusters)

	return resource, nil
}

func (xds *AppHandler) GetNonEdsClusters(ctx context.Context, clusterNames []string, requestDetails models.RequestDetails, version string) ([]interface{}, error) {
	resource := &models.DBResource{}
	collection := xds.Context.Client.Collection("clusters")
	results := []interface{}{}
	for _, clusterName := range clusterNames {
		filter := bson.M{"general.name": clusterName, "general.project": requestDetails.Project, "general.version": version}
		filter = common.AddUserFilter(requestDetails, filter)
		result := collection.FindOne(ctx, filter)

		if result.Err() != nil {
			if errors.Is(result.Err(), mongo.ErrNoDocuments) {
				return nil, errors.New("not found: (" + clusterName + ")")
			}
			return nil, errstr.ErrUnknownDBError
		}

		err := result.Decode(resource)
		if err != nil {
			return nil, err
		}

		general := resource.GetGeneral()
		if len(general.TypedConfig) != 0 {
			for _, typed := range general.TypedConfig {
				if typed.Gtype == "envoy.extensions.upstreams.http.v3.HttpProtocolOptions" {
					res := resource.GetResource()
					cluster, ok := res.(primitive.M)
					if !ok {
						return nil, fmt.Errorf("failed to parse cluster")
					}

					protocolOptions, err := xds.GetHttpProtocolOptions(ctx, typed.Collection, typed.Name, requestDetails, version)
					if err != nil {
						return nil, err
					}

					if len(protocolOptions) != 0 {
						cluster["typed_extension_protocol_options"] = protocolOptions
					}

					results = append(results, cluster)
				}
			}
		} else {
			results = append(results, resource.GetResource())
		}
	}

	return results, nil
}

func (xds *AppHandler) GetHttpProtocolOptions(ctx context.Context, collectionName, name string, requestDetails models.RequestDetails, version string) (primitive.M, error) {
	resource := &models.DBResource{}
	collection := xds.Context.Client.Collection(collectionName)
	filter := bson.M{"general.name": name, "general.project": requestDetails.Project, "general.version": version}
	result := collection.FindOne(ctx, filter)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("not found: (" + name + ")")
		}
		return nil, errstr.ErrUnknownDBError
	}

	err := result.Decode(resource)
	if err != nil {
		return nil, err
	}

	resourceData, ok := resource.GetResource().(primitive.M)
	if !ok {
		return nil, fmt.Errorf("failed to parse resource.GetResource() as primitive.M, got: %T", resource.GetResource())
	}

	httpProtocolOptions := primitive.M{
		"@type": "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
	}

	for key, value := range resourceData {
		httpProtocolOptions[key] = value
	}

	return httpProtocolOptions, nil
}

func (xds *AppHandler) collectAccessLoggers(ctx context.Context, resource models.DBResourceClass, requestDetails models.RequestDetails, version string) (models.DBResourceClass, error) {
	bootstrap := resource.GetResource()
	bootstrapMap, ok := bootstrap.(primitive.M)
	if !ok {
		return nil, fmt.Errorf("failed to parse bootstrap as primitive.M, got type: %T", bootstrap)
	}

	admin, ok := bootstrapMap["admin"].(primitive.M)
	if !ok {
		return nil, fmt.Errorf("'admin' key not found or invalid")
	}

	accessLog, ok := admin["access_log"].(primitive.A)
	if !ok {
		return nil, fmt.Errorf("'access_log' key not found or invalid")
	}

	var accessLogs []string
	for _, aclog := range accessLog {
		acLogMap, ok := aclog.(primitive.M)
		if !ok {
			continue
		}

		if typedConfig, ok := acLogMap["typed_config"].(primitive.M); ok {
			typedConf, err := resources.DecodeBase64Config(typedConfig["value"].(string))
			if err != nil {
				return nil, err
			}
			accessLogs = append(accessLogs, typedConf.Name)
		}
	}

	accessLoggers, err := xds.GetAccessLoggers(ctx, accessLogs, requestDetails, version)
	if err != nil {
		return nil, err
	}

	resource.SetBootstrapAccessLoggers(accessLoggers)

	return resource, nil
}

func (xds *AppHandler) GetAccessLoggers(ctx context.Context, alNames []string, requestDetails models.RequestDetails, version string) ([]interface{}, error) {
	resource := &models.DBResource{}
	collection := xds.Context.Client.Collection("extensions")
	results := []interface{}{}
	for _, alName := range alNames {
		filter := bson.M{"general.name": alName, "general.project": requestDetails.Project, "general.version": version}
		result := collection.FindOne(ctx, filter)

		if result.Err() != nil {
			if errors.Is(result.Err(), mongo.ErrNoDocuments) {
				return nil, errors.New("not found: (" + alName + ")")
			}
			return nil, errstr.ErrUnknownDBError
		}

		err := result.Decode(resource)
		if err != nil {
			return nil, err
		}

		resourceData, ok := resource.GetResource().(primitive.M)
		if !ok {
			return nil, fmt.Errorf("failed to parse resource.GetResource() as primitive.M, got: %T", resource.GetResource())
		}

		general := resource.GetGeneral()
		typedConfig := models.TC{
			Name: general.CanonicalName,
			TypedConfig: map[string]interface{}{
				"@type": "type.googleapis.com/" + general.GType,
			},
		}

		for key, value := range resourceData {
			typedConfig.TypedConfig[key] = value
		}

		results = append(results, typedConfig)
	}

	return results, nil
}
