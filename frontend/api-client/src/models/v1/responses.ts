import {SCMIntegration} from './entities'
import {KVMap, KVStore} from '../common'
import {ValidationErrorList} from './errors'

type BaseResponseParams = {
    errors: ValidationErrorList, 
    statusCode: number, 
    expectedStatusCode: number
}

// class Response {
//     readonly errors: ValidationErrorList
//     readonly statusCode: number
//     readonly expectedStatusCode: number
    
//     constructor({errors, statusCode, internalErrors, expectedStatusCode}: BaseResponseParams) {
//         this.errors = errors
//         this.statusCode = statusCode
//         this.expectedStatusCode = expectedStatusCode
//     }

//     WasSuccessful(): boolean {
//         return this.expectedStatusCode === this.statusCode
//     }

//     ValidationErrors(): ValidationErrorList {
//         return this.errors
//     }

//     // Errors(): string[] {
//     //     return this.internalErrors
//     // }


// }

// type ListSCMIntegrationsResponseParams = BaseResponseParams & {
//     data: SCMIntegration[],
//     meta: ResponseMeta
// }

// class ListSCMIntegrationsResponseObj extends Response implements ListSCMIntegrationsResponse  {
//     readonly data: SCMIntegration[]
//     readonly meta: ResponseMeta

//     constructor({data, meta, errors, statusCode}: ListSCMIntegrationsResponseParams) {
//         super({errors, statusCode, expectedStatusCode: 200})

//         this.data = data
//         this.meta = meta
//     }

//     GetSCMIntegrations(): SCMIntegration[] {
//         return this.data
//     }
    
//     GetMeta(): ResponseMeta {
//         return this.meta
//     }
// }

interface ResponseMeta {
    Get(k: string): string | null,
    Has(k: string): boolean,
}

interface ListSCMIntegrationsResponse {
    GetSCMIntegrations(): SCMIntegration[],
    GetMeta(): ResponseMeta,
    DidFail?(): boolean,
    Errors(): string[]
    WasSuccessful(): boolean,
    ValidationErrors(): ValidationErrorList,
}

export {
    ListSCMIntegrationsResponse,
}