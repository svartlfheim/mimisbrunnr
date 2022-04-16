import styles from './Content.module.css'
import {Panel} from '../Layout'

enum Mode {
    ResponsiveColumns,
    SingleColumn,
}

type Props = {
    children: React.ReactNode,
    mode: Mode,
}

function getGridModeModifier(m: Mode): string {
    switch(m) {
        case Mode.ResponsiveColumns:
            return styles.responsiveColumns
        case Mode.SingleColumn:
            return styles.singleColumn
        default:
            console.warn("unknown grid modifier")
            return ""
    }
}

function Content({children, mode}: Props) {
    const modifier = getGridModeModifier(mode)
    return (
        <div className={`${styles.content} ${modifier}`}>
            {children}
        </div>
    )
}

export {
    Content,
    Mode,
}


export default Content;