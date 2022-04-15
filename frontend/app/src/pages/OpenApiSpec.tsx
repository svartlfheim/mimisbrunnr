import Base from './Base'
import SwaggerUI from "swagger-ui-react"
import "swagger-ui-react/swagger-ui.css"

function OpenApiSpec() {
    return (
        <Base>
        {/* <h3>OpenAPI Spec</h3> */}
          <SwaggerUI url="/static/openapi.json" />
        </Base>
    )
}


export default OpenApiSpec;