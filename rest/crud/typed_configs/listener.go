package typed_configs

/* func getListenerTypedConfigs(listener models.DBResourceClass, logger *logrus.Logger) []*models.TypedConfig {
	var typedConfigs []*models.TypedConfig
	resource := listener.GetResource()

	listeners, _ := resource.([]interface{})

	for _, lr := range listeners {
		jsonStringStr, err := helper.MarshalJSON(lr, logger)
		if err != nil {
			return typedConfigs
		}

		typedConfigsPart, _ := resources.ProcessFilterChains(jsonStringStr, models.TransportSocketPath.String(), logger)
		typedConfigs = append(typedConfigs, typedConfigsPart...)
	}

	return typedConfigs
} */
