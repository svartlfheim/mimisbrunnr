import {KVStore} from '../common'

interface ValidationError {
    path: string,
    message: string,
    params: KVStore,
    rule: string,
}

interface ValidationErrorList {
    ExistsForPath(p: string): boolean,
    ForPath(p: string): ValidationError[] | null,
    IsEmpty(): boolean,
}

class ErrorStore implements ValidationErrorList {
    readonly errors: ValidationError[]

    constructor(errors: ValidationError[]) {
        this.errors = errors
    }

    ExistsForPath(p: string): boolean {
        return this.ForPath(p).length > 0
    }

    ForPath(p: string): ValidationError[] {
        return this.errors.filter((v) => v.path === p)
    }

    IsEmpty(): boolean {
        return this.errors.length === 0
    }
}

export {
    ValidationError,
    ValidationErrorList,
    ErrorStore,
}