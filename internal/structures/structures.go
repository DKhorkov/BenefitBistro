package structures

type RouteHandlersInfoStructure struct {
	StaticFiles, 
	HomePage, 
	EmployeeLogin, 
	EmployeeRegister, 
	HirerLogin, 
	HirerRegister, 
	SaveEmployee RouteInfoStructure
}

type RouteInfoStructure struct {
	TemplateName, URLPath string
}

type ServerParametersStructure struct {
	Host string
	Port string
}

type TemplateParamsStructure struct {
	HomePage, EmployeeLogin, EmployeeRegister, HirerLogin, HirerRegister TemplateData
}

type TemplateData struct {
	PageName string
}