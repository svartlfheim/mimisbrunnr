import styles from './Content.module.css'

type Props = {
    children: React.ReactNode,
}

function Content({children}: Props) {
    return (
        <div className={styles.contentWrapper}>
            <div className={styles.content}>
                {children}
            </div>
        </div>
    )
}

export default Content