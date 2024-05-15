package domain

func AppendLastResourceToResourceMap(mapping map[string][]Resource, key string) (map[string][]Resource, Resource) {
	resource := LastResource(mapping, key)
	if resource == (Resource{}) || resource.IsLastPicture() {
		resource = NewResource()
		mapping[key] = append(mapping[key], resource)
	}
	return mapping, resource
}

func AppendResourceMap(mapping map[string][]Resource, key string, value Resource) map[string][]Resource {
	if len(mapping[key]) == 0 {
		return mapping
	}
	lastIndex := len(mapping[key]) - 1
	mapping[key][lastIndex] = value
	return mapping
}

func LastResource(mapping map[string][]Resource, key string) Resource {
	if len(mapping[key]) == 0 {
		return Resource{}
	}
	lastIndex := len(mapping[key]) - 1
	return mapping[key][lastIndex]
}
