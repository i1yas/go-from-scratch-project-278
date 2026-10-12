package httpapi

var fieldsMapping = []struct {
	serviceName string
	httpName    string
}{
	{
		serviceName: "shortcode",
		httpName:    "short_name",
	},
	{
		serviceName: "referer",
		httpName:    "reffer",
	},
}

func mapFieldFromServiceToHTTP(name string) string {
	for _, m := range fieldsMapping {
		if m.serviceName == name {
			return m.httpName
		}
	}

	return name
}

func mapFieldFromHTTPtoService(name string) string {
	for _, m := range fieldsMapping {
		if m.httpName == name {
			return m.serviceName
		}
	}

	return name
}
