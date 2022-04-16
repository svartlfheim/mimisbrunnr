import styles from './Panel.module.css'

interface Props {
    children: React.ReactElement[]|React.ReactElement,
    fullSize?: boolean,
    className?: string,
    theme?: string,
    noPadding?: boolean,
}

function Panel({children, fullSize, className, theme, noPadding}: Props) {
    let classes = `${styles.panel}`
    classes += className !== undefined ? ` ${className}` : ''

    if (fullSize) {
        classes += ` ${styles.fullSize}`
    }

    if (theme !== undefined) {
        classes += ` ${styles[theme]}`
    }

    if (! noPadding) {
        classes += ` ${styles.padded}`
    }

    return <div className={classes}>
        {children}
    </div>
}

export {
    Panel,
}