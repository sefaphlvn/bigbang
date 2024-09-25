package dependency

import (
	"fmt"

	"github.com/sefaphlvn/bigbang/pkg/filters"
	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/tidwall/gjson"
)

func GenericUpstreamCollector(ctx *AppHandler, activeResource Depend) (Node, []Depend) {
	var dependencies []Depend
	id, jsonData := ctx.getResourceData(activeResource.Collection, activeResource.Name, activeResource.Project)
	rootResult := gjson.Parse(jsonData)

	node := Node{
		ID:         id,
		Name:       activeResource.Name,
		Gtype:      activeResource.Gtype,
		Collection: activeResource.Gtype.CollectionString(),
		Link:       activeResource.Gtype.URL(),
		First:      activeResource.First,
		Direction:  "upstream",
	}

	jsonPaths := getDynamicJsonPaths(activeResource.Gtype)
	for path, gtype := range jsonPaths {
		resourcePath := fmt.Sprintf("%s.%s", "resource.resource", path)
		collectDependenciesFromPath(ctx, rootResult, resourcePath, gtype, activeResource, &dependencies)
	}

	dependencies = append(dependencies, parseTypedConfig(ctx, rootResult, activeResource)...)
	dependencies = append(dependencies, parseConfigDiscovery(ctx, rootResult, activeResource)...)

	if len(dependencies) == 0 {
		ctx.Context.Logger.Debugf("No dependencies found for resource: %s of type %s", activeResource.Name, activeResource.Gtype)
	}

	return node, dependencies
}

func collectDependenciesFromPath(ctx *AppHandler, rootResult gjson.Result, path string, gtype models.GTypes, activeResource Depend, dependencies *[]Depend) {
	results := rootResult.Get(path)

	if !results.Exists() {
		ctx.Context.Logger.Debugf("Result does not exist at path: %s", path)
		return
	}

	results.ForEach(func(_, item gjson.Result) bool {
		if item.IsArray() {
			item.ForEach(func(_, subItem gjson.Result) bool {
				processItem(ctx, subItem, path, gtype, activeResource, dependencies)
				return true
			})
		} else {
			processItem(ctx, item, path, gtype, activeResource, dependencies)
		}
		return true
	})
}

func processItem(ctx *AppHandler, item gjson.Result, path string, gtype models.GTypes, activeResource Depend, dependencies *[]Depend) {
	if item.IsArray() {
		item.ForEach(func(_, subItem gjson.Result) bool {
			addDependency(ctx, subItem.String(), path, gtype, activeResource, dependencies)
			return true
		})
	} else {
		addDependency(ctx, item.String(), path, gtype, activeResource, dependencies)
	}
}

func addDependency(ctx *AppHandler, name, path string, gtype models.GTypes, activeResource Depend, dependencies *[]Depend) {
	if name == "" {
		ctx.Context.Logger.Debugf("Name not found at path: %s for gtype: %s", path, gtype)
		return
	}

	itemID, _ := ctx.getResourceData(gtype.CollectionString(), name, activeResource.Project)
	if itemID == "" {
		ctx.Context.Logger.Debugf("ID not found for %s of type %s, skipping... Path: %s", name, gtype, path)
		return
	}

	dependency := Depend{
		Name:       name,
		Gtype:      gtype,
		Collection: gtype.CollectionString(),
		Project:    activeResource.Project,
		ID:         itemID,
		Direction:  "upstream",
	}

	*dependencies = append(*dependencies, dependency)
	ctx.Context.Logger.Debugf("Added dependency: %s of type %s with ID: %s", dependency.Name, dependency.Gtype, dependency.ID)
}

func GenericDownstreamCollector(ctx *AppHandler, activeResource Depend, visited map[string]bool) (Node, []Depend) {
	var dependencies []Depend
	if activeResource.ID == "" {
		id, _ := ctx.getResourceData(activeResource.Collection, activeResource.Name, activeResource.Project)
		activeResource.ID = id
	}

	uniqueKey := fmt.Sprintf("%s_%s_%s_%s", activeResource.Name, activeResource.Gtype, activeResource.Collection, activeResource.Project)
	if visited[uniqueKey] {
		return Node{}, nil
	}
	visited[uniqueKey] = true

	node := Node{
		ID:         activeResource.ID,
		Name:       activeResource.Name,
		Gtype:      activeResource.Gtype,
		Collection: activeResource.Gtype.CollectionString(),
		Link:       activeResource.Gtype.URL(),
		First:      activeResource.First,
		Direction:  "downstream",
	}

	downstreamFilters := activeResource.Gtype.DownstreamFilters(activeResource.Name)
	for _, filter := range downstreamFilters {
		collectDependenciesFromFilter(ctx, filter, activeResource, &dependencies)
	}

	for _, dep := range dependencies {
		if dep.Direction == "downstream" {
			_, downstreamDeps := GenericDownstreamCollector(ctx, dep, visited)
			dependencies = append(dependencies, downstreamDeps...)
		}
	}

	return node, dependencies
}

func collectDependenciesFromFilter(ctx *AppHandler, filter filters.MongoFilters, activeResource Depend, dependencies *[]Depend) {
	collection := filter.Collection
	query := filter.Filter

	cursor, err := ctx.Context.Client.Collection(collection).Find(ctx.Context.Ctx, query)
	if err != nil {
		ctx.Context.Logger.Debugf("Error fetching downstream dependencies: %v", err)
		return
	}

	for cursor.Next(ctx.Context.Ctx) {
		var resource models.DBResource
		if err := cursor.Decode(&resource); err != nil {
			ctx.Context.Logger.Debugf("Error decoding downstream resource: %v", err)
			continue
		}

		if resource.General.Name == activeResource.Name && resource.General.GType == activeResource.Gtype {
			continue
		}

		dependency := Depend{
			Name:       resource.General.Name,
			Gtype:      resource.General.GType,
			Collection: collection,
			Project:    activeResource.Project,
			ID:         resource.ID.Hex(),
			Direction:  "downstream",
			Source:     activeResource.ID,
		}

		*dependencies = append(*dependencies, dependency)
		ctx.Context.Logger.Debugf("Added downstream dependency: %s of type %s with ID: %s", dependency.Name, dependency.Gtype, dependency.ID)
	}
}
