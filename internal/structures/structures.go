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
	HomePage, EmployeeLogin, EmployeeRegister, HirerLogin, HirerRegister TemplateDataStructure
}

type TemplateDataStructure struct {
	PageName string
	UserData UserDataStructure
}

type UserDataStructure struct {
	Username string
	Authentificated bool
}

type TokenStruct struct {
	Name string
	LifeTime int
}