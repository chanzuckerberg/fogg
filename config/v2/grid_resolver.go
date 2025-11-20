package v2

func ResolveGrid(commons ...Common) *GridConfig {
	var enabled *bool
	var guid *string
	var endpoint *string

	for _, c := range commons {
		if c.Grid != nil {
			if c.Grid.Enabled != nil {
				enabled = c.Grid.Enabled
			}
			if c.Grid.GUID != nil {
				guid = c.Grid.GUID
			}
			if c.Grid.Endpoint != nil {
				endpoint = c.Grid.Endpoint
			}
		}
	}

	if enabled == nil && guid == nil && endpoint == nil {
		return nil
	}

	return &GridConfig{
		Enabled:  enabled,
		GUID:     guid,
		Endpoint: endpoint,
	}
}
