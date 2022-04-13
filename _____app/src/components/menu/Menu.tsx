import styles from './Menu.module.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHome as faHome } from '@fortawesome/free-solid-svg-icons'
import { faCodeBranch as faIntegrations } from '@fortawesome/free-solid-svg-icons'
import { faDiagramProject as faProjects } from '@fortawesome/free-solid-svg-icons'
import { faCog as faSettings } from '@fortawesome/free-solid-svg-icons'
import { faQuestion as faUnknown } from '@fortawesome/free-solid-svg-icons'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import { NavLink } from 'react-router-dom'
import { CSSProperties } from "react"



// import { faArrowsLeftRightToLine as faIntegrations } from '@fortawesome/free-solid-svg-icons'
// import { faBezierCurve as faIntegrations } from '@fortawesome/free-solid-svg-icons'
// import { faBookAtlas as faProjects } from '@fortawesome/free-solid-svg-icons'


type ItemProps = {
    title?: string;
    displayTitle: boolean;
    icon?: IconDefinition;
    route: string;
}


function MenuItem({ displayTitle, title, icon, route }: ItemProps) {
    const titleElement = displayTitle ?
        (<span className={styles.itemTitle}>{title ?? 'unknown'}</span>) :
        (<></>)

    return (
        <NavLink to={route} className={({ isActive }) =>
            isActive ? styles.activeLink : ''
        }>
            <li>
                <span className={styles.itemIcon}>
                    <FontAwesomeIcon icon={icon ?? faUnknown} />
                </span>
                {titleElement}
            </li>
        </NavLink>
    )
}

type ItemRoute = {
    display?: string;
    path: string;
    icon?: IconDefinition;
}

type Props = {
    isOpen: boolean;
    routes: ItemRoute[]
}

function Menu({ isOpen, routes }: Props) {
    const menuStyles = `${styles.menu} ${isOpen ? styles.open : styles.closed}`
    return (
        <div className={menuStyles}>
            <ul>
                {routes.map((route: ItemRoute, i: number) => {
                    return (
                        <MenuItem key={i} route={route.path} title={route.display} icon={route.icon} displayTitle={isOpen} />
                    )
                })}
                {/* <MenuItem route="/" title="Home" icon={faHome} displayTitle={isOpen} />
                <MenuItem route="/projects" title="Projects" icon={faProjects} displayTitle={isOpen} />
                <MenuItem route="/scm-integrations" title="SCM Integrations" icon={faIntegrations} displayTitle={isOpen} />
                <MenuItem route="/settings" title="Settings" icon={faSettings} displayTitle={isOpen} /> */}
            </ul>
        </div>
    )
}

export default Menu;