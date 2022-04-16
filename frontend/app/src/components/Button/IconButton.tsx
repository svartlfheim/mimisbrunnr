import styles from './IconButton.module.css'
import React from "react"
import {Type, classForType} from './Types'

interface Props {
    title?: string,
    icon: React.ReactElement,
    type: Type,
}

function IconButton({title, icon, type}: Props) {
    const titleElement = title !== undefined ? (
        <span className={styles.title}>{title}</span>
    ): <></>

    const typeClass = classForType(type)
    return (
        <button onClick={() => true} className={`${styles.button} ${typeClass}`}>
            <span className={styles.icon}>{icon}</span>
            {titleElement}
        </button>
    )
}

export {
    IconButton,
}