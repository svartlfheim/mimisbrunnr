import {Base} from '../Common'
import SwaggerUI from "swagger-ui-react"
import "swagger-ui-react/swagger-ui.css"
import {Mode} from '../../Components/Content'
import { Panel } from '../../Components/Layout'

function OpenApiSpec() {
    return (
        <Base gridMode={Mode.SingleColumn}>
          <Panel fullSize theme="light">
            <SwaggerUI url="/static/openapi.json" />
          </Panel>
        </Base>
    )
}


export {
  OpenApiSpec,
}