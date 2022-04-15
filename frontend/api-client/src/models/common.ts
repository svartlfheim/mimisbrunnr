type KVMap = { [key: string]: string; }

class KVStore {
    readonly kv: KVMap

    constructor(kv: KVMap) {
        this.kv = kv
    }

    Get(k: string): string | null {
        if (! this.Has(k)) {
            return null;
        }

        return this.kv[k];
    }

    Has(k: string): boolean {
        return this.kv.hasOwnProperty(k) && this.kv[k] !== undefined && this.kv[k] !== null
    }
}

export {
    KVMap,
    KVStore,
}