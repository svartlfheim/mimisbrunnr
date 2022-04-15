import * as v1Models from "./models/v1"

interface ProjectsApiV1 {
    
}


interface SCMIntegrationsApiV1 {
    List(params: v1Models.ListRequestParameters): Promise<v1Models.ListSCMIntegrationsResponse>;
    Create(params: v1Models.CreateSCMIntegrationRequest): Promise<v1Models.SCMIntegration>;
}

interface ApiV1 {
    Projects: ProjectsApiV1,
    SCMIntegrations: SCMIntegrationsApiV1,
}

interface Api {
    V1: ApiV1,
}

