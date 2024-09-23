package resource

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/pkg/resources"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
)

// DecodeListener processes a raw listener resource and collects extensions.
func (ar *AllResources) DecodeListener(rawListenerResource *models.DBResource, context *db.AppContext, logger *logrus.Logger) {
	ar.UniqueResources = make(map[string]struct{})
	if err := ar.initializeListener(rawListenerResource, context, logger); err != nil {
		logger.Fatalf("Error initializing listener: %v", err)
	}

	ar.processConfigDiscoveries(rawListenerResource.General.ConfigDiscovery, context, logger)
	fmt.Println("Listener decoded")
}

// Initialize listener by decoding and setting up listener resources.
func (ar *AllResources) initializeListener(rawListenerResource *models.DBResource, context *db.AppContext, logger *logrus.Logger) error {
	resArray, ok := rawListenerResource.Resource.Resource.(primitive.A)
	if !ok {
		return fmt.Errorf("unexpected resource format")
	}

	newVersion, err := resources.IncrementResourceVersion(context, rawListenerResource.General.Name, rawListenerResource.General.Project)
	if err != nil {
		return err
	}
	ar.SetVersion(newVersion)
	ar.SetProject(rawListenerResource.General.Project)

	var listeners []types.Resource
	for _, lstnr := range resArray {
		listenerWithTransportSocket, _ := ar.GetTypedConfigs(models.ConfigGetters[rawListenerResource.GetGtype()], lstnr, context)

		singleListener := &listener.Listener{}
		if err := resources.MarshalUnmarshalWithType(listenerWithTransportSocket, singleListener); err != nil {
			logger.Errorf("Listener Unmarshal error: %s", err)
			continue
		}

		listeners = append(listeners, singleListener)
	}

	ar.SetListener(listeners)
	return nil
}

// Process config discoveries and collect resources.
func (ar *AllResources) processConfigDiscoveries(configDiscoveries []*models.ConfigDiscovery, context *db.AppContext, logger *logrus.Logger) {
	for _, configDiscovery := range configDiscoveries {
		ar.processExtension(configDiscovery, configDiscovery.ParentName, context, logger)
	}
}

// Process a single extension, collect resources, and add to AllResources if not duplicate.
func (ar *AllResources) processExtension(extension *models.ConfigDiscovery, parentName string, context *db.AppContext, logger *logrus.Logger) {
	uniqKey := fmt.Sprintf("%s__%s__%s", extension.Name, parentName, extension.GType.String())
	if ar.checkAndMarkDuplicate(uniqKey) {
		return
	}

	extConfigs, additionalExtResources, err := ar.CollectAllResourcesWithParent(extension.GType, extension.Name, parentName, context, logger)
	if err != nil {
		logger.Errorf("Error collecting resources: %v", err)
		return
	}

	for _, extConfig := range extConfigs {
		if extension.GType == models.VirtualHost {
			uniqKey := fmt.Sprintf("%s__%s", extension.Name, extension.GType.String())

			ar.AddToCollection(extConfig, extension.GType, uniqKey, &parentName)

		} else {
			typedConfigAsAny, err := anypb.New(extConfig)
			if err != nil {
				logger.Errorf("Error converting extTypedConfig to *anypb.Any for %s: %v", parentName, err)
				continue
			}

			ar.addTypedExtensionConfig(typedConfigAsAny, parentName)
		}
	}

	if additionalExtResources != nil {
		ar.processConfigDiscoveries(additionalExtResources, context, logger)
	}
}

// Add a typed extension configuration to the Extensions slice.
func (ar *AllResources) addTypedExtensionConfig(typedConfig *anypb.Any, parentName string) {

	typedExtensionConfig := &core.TypedExtensionConfig{
		Name:        parentName,
		TypedConfig: typedConfig,
	}
	ar.Extensions = append(ar.Extensions, typedExtensionConfig)
}

// Check if a resource is a duplicate and mark it as processed.
func (ar *AllResources) checkAndMarkDuplicate(name string) bool {
	if _, exists := ar.UniqueResources[name]; exists {
		return true
	}

	ar.UniqueResources[name] = struct{}{}
	return false
}

// CollectAllResourcesWithParent processes and collects resources with a parent.
func (ar *AllResources) CollectAllResourcesWithParent(gtype models.GTypes, resourceName, parentName string, context *db.AppContext, logger *logrus.Logger) ([]proto.Message, []*models.ConfigDiscovery, error) {
	resource, err := resources.GetResourceNGeneral(context, gtype.CollectionString(), resourceName, ar.Project)
	if err != nil {
		logger.Errorf("Error getting resource %s: %v", resourceName, err)
		return nil, nil, err
	}

	// resource.GetResource() array olabilir, bu yüzden array olup olmadığını kontrol ediyoruz.
	resourceData := resource.GetResource()

	var protoMessages []proto.Message
	var finalConfigDiscoveries []*models.ConfigDiscovery

	switch res := resourceData.(type) {
	case primitive.A: // Eğer array ise
		for _, item := range res {
			// Array'in her elemanını işleyerek JSON string'e çeviriyoruz
			jsonStringStr, err := helper.MarshalJSON(item, context.Logger)
			if err != nil {
				logger.Errorf("Error marshalling array item: %v", err)
				return nil, nil, err
			}

			typedProtoMsg := gtype.ProtoMessage()
			if err := ar.processTypedConfigsAndUpstream(typedProtoMsg, &jsonStringStr, gtype, parentName, context, logger); err != nil {
				logger.Errorf("Error processing typed configs and upstream resources: %v", err)
				return nil, nil, err
			}

			// Sonucu biriktiriyoruz
			protoMessages = append(protoMessages, typedProtoMsg)
			finalConfigDiscoveries = resource.General.ConfigDiscovery
			ar.processConfigDiscoveries(resource.General.ConfigDiscovery, context, logger)
		}
	default: // Eğer array değilse, tek bir resource olarak işlem yapıyoruz
		jsonStringStr, err := helper.MarshalJSON(resourceData, context.Logger)
		if err != nil {
			return nil, nil, err
		}

		typedProtoMsg := gtype.ProtoMessage()
		if err := ar.processTypedConfigsAndUpstream(typedProtoMsg, &jsonStringStr, gtype, parentName, context, logger); err != nil {
			logger.Errorf("Error processing typed configs and upstream resources: %v", err)
			return nil, nil, err
		}

		// Tek resource için ConfigDiscovery'leri işliyoruz
		protoMessages = append(protoMessages, typedProtoMsg)
		finalConfigDiscoveries = resource.General.ConfigDiscovery
		ar.processConfigDiscoveries(resource.General.ConfigDiscovery, context, logger)
	}

	// Döngü tamamlandıktan sonra sonuçları döndürüyoruz
	return protoMessages, finalConfigDiscoveries, nil
}

// Process typed configs and upstream paths.
func (ar *AllResources) processTypedConfigsAndUpstream(protoMsg proto.Message, jsonStringStr *string, gtype models.GTypes, parentName string, context *db.AppContext, logger *logrus.Logger) error {
	ar.processTypedConfigPaths(models.ConfigGetters[gtype], jsonStringStr, context, logger)
	ar.processUpstreamPaths(gtype.GetUpstreamPaths(), jsonStringStr, parentName, context, logger)

	if err := protojson.Unmarshal([]byte(*jsonStringStr), protoMsg); err != nil {
		logger.Errorf("Error unmarshalling to proto message after processing nested configs: %v", err)
		return err
	}
	return nil
}

// Process typed config paths.
func (ar *AllResources) processTypedConfigPaths(configPaths []models.TypedConfigPath, jsonStringStr *string, context *db.AppContext, logger *logrus.Logger) {
	for _, path := range configPaths {
		if err := ar.processTypedConfigPath(path, jsonStringStr, context); err != nil {
			logger.Warnf("Error processing typed config path: %v", err)
		}
	}
}

// Process upstream paths recursively.
func (ar *AllResources) processUpstreamPaths(upstreamPaths map[string]models.GTypes, jsonStringStr *string, parentName string, context *db.AppContext, logger *logrus.Logger) {
	for jsonPath, upstreamType := range upstreamPaths {
		result := gjson.Get(*jsonStringStr, jsonPath)
		if result.Exists() {
			processUpstreamPaths(result, upstreamType, parentName, ar, context, logger)
		}
	}
}

// Process upstream paths recursively.
func processUpstreamPaths(result gjson.Result, upstreamType models.GTypes, parentName string, ar *AllResources, context *db.AppContext, logger *logrus.Logger) {

	if result.IsArray() {
		result.ForEach(func(_, item gjson.Result) bool {
			processUpstreamPaths(item, upstreamType, parentName, ar, context, logger)
			return true
		})
	} else {
		resourceName := result.String()
		upstreamResourceProtoMsgs, additionalExtResources, err := ar.CollectAllResourcesWithParent(upstreamType, resourceName, parentName, context, logger)
		if err != nil {
			logger.Errorf("Error collecting upstream resources: %v", err)
			return
		}

		for _, protoMsg := range upstreamResourceProtoMsgs {
			uniqKey := fmt.Sprintf("%s__%s", resourceName, upstreamType.String())
			if protoMsg != nil {
				ar.AddToCollection(protoMsg, upstreamType, uniqKey, nil)
			}
		}
		if additionalExtResources != nil {
			ar.processConfigDiscoveries(additionalExtResources, context, logger)
		}
	}
}

// Add a resource to the appropriate collection in AllResources.
func (ar *AllResources) AddToCollection(resource proto.Message, gtype models.GTypes, uniqName string, parentName *string) {
	if ar.checkAndMarkDuplicate(uniqName) {
		fmt.Printf("Skipping duplicate collection of resource: %s", uniqName)
		return
	}

	switch gtype {
	case models.Cluster:
		newCluster := proto.Clone(resource).(*cluster.Cluster)
		ar.Cluster = append(ar.Cluster, newCluster)
	case models.Route:
		newRoute := proto.Clone(resource).(*route.RouteConfiguration)
		ar.Route = append(ar.Route, newRoute)
	case models.Endpoint:
		newEndpoint := proto.Clone(resource).(*endpoint.ClusterLoadAssignment)
		ar.Endpoint = append(ar.Endpoint, newEndpoint)
	case models.VirtualHost:
		newVirtualHost := proto.Clone(resource).(*route.VirtualHost)
		newVirtualHost.Name = fmt.Sprintf("%s/%s", *parentName, newVirtualHost.Name)
		ar.VirtualHost = append(ar.VirtualHost, newVirtualHost)
	case models.DownstreamTlsContext, models.UpstreamTlsContext, models.CertificateValidationContext, models.TlsCertificate:
		newSecret := proto.Clone(resource).(*tls.Secret)
		ar.Secret = append(ar.Secret, newSecret)
	default:
		ar.Extensions = append(ar.Extensions, resource)
	}
}
