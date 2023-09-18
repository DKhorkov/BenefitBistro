package structures

type RouteHandlersInfoStructure struct {
	HomePage, StaticFiles RouteInfoStructure
}

type RouteInfoStructure struct {
	TemplateName, URLPath string
}

type ServerParametersStructure struct {
	Host string
	Port string
}

type TemplateParamsStructure struct {
	HomePage TemplateData
}

type TemplateData struct {
	PageName string
}