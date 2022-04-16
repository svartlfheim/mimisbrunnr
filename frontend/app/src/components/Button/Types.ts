import styles from './Types.module.css'

enum Type {
    Neutral,
    Positive,
    Caution,
    Danger,
}

const typeClassMap: {[key in Type]: string} = {
    [Type.Neutral]: styles.neutral,
    [Type.Positive]: styles.positive,
    [Type.Caution]: styles.caution,
    [Type.Danger]: styles.danger,
}

function classForType(t: Type): string {
    const className: string | undefined = typeClassMap[t]

    if (className === undefined) {
        return typeClassMap[Type.Neutral]
    }

    return className
}

export {
    Type,
    classForType,
}