import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faQuestion as faUnknown } from '@fortawesome/free-solid-svg-icons'
import {
    NavLink, 
    useMatch,
    useResolvedPath,
} from 'react-router-dom'
import { Route } from "../../../Service/router"

type StyleProps = {
    itemTitle: string,
    itemIcon: string,
    itemLinkWrapper: string,
    activeLink: string,
}

interface Props {
    hideIcon?: boolean,
    childMenuIsOpen?: boolean,
    route: Route,
    pathOverride?: string,
    menuIsOpen: boolean,
    displayTitle: boolean,
    toggle?: React.ReactElement,
    styles: StyleProps ,
}

function Link({route, menuIsOpen, childMenuIsOpen, displayTitle, toggle, pathOverride, hideIcon, styles}: Props) {
    const finalPath = pathOverride !== undefined ? pathOverride : route.path
    const resolved = useResolvedPath(finalPath);
    const match = useMatch({ path: resolved.pathname, end: route.path === "/" || (menuIsOpen && childMenuIsOpen)});

    const titleElement = displayTitle ?
        (<span className={styles.itemTitle}>{route.display ?? 'unknown'}</span>) :
        (<></>)

    const icon = !hideIcon ? (
        <span className={styles.itemIcon}>
            <FontAwesomeIcon icon={route.icon ?? faUnknown} />
        </span>) : <></>

    return (
        <>
            <NavLink to={finalPath}>
                <div className={`${styles.itemLinkWrapper} ` + (match ? styles.activeLink : '')}>
                        {icon}
                        {titleElement}
                </div>
            </NavLink>
            {toggle}
        </>
    )
}

export type {
    StyleProps,
}

export {
    Link,
}
