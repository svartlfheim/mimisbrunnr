# TSConfig

<!-- JSON doesn't allow comments, so here we are...notes about the choices in the tsconfig.json.

## `compilerOptions.typeRoots`

This is set as:

```json
{
    "typeRoots": [
        "node_modules/@types",
        "../node_modules/@types",
        "src/types"
    ]
}
```

Basically the generated code we use for the api client in `./src/fetch` relies on the `GlobalFetch` property, which used to come with the "dom" lib that's included. The type was renamed to `WindowOrWorkerGlobalScope` this. If we want our generated code to compile we need to essentially polyfil this type, which can be done with `declare type GlobalFetch = WindowOrWorkerGlobalScope`, just before the generated code uses the GlobalFetch. 

However this solution requires us to modify the generated code, every time we regenerate it; which kinds defeats the purpose of using the generated code as it will be a huge pain in the arse. So to workaround this problem we can declare a global type using: 

```ts
declare global {
    type GlobalFetch = WindowOrWorkerGlobalScope
}
```

Problem is we need to include this type somewhere, and the only way to really add additional types from outside the `node_modules` directory is to include them in the `typeRoots` config option. this comes with it's own problem which is that this overrides the default behaviour. The default behaviour will look for types in all parent folders `node_modules` directories as well, which is really useful as we're using yarn workspaces which hoist our dependencies up to the root. So until there is an option to append to the typeRoots directory we'll have to include all the paths we need manually, which is fortunately just the current and parent directories node_modules. -->

## `lib[*] "dom"`

The "dom" library is used to expose the types of the standard dom in web browsers. This includes things like fetch and window etc.

See: https://www.typescriptlang.org/tsconfig#lib
