import styles from './Menu.module.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faQuestion as faUnknown, faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import {
    NavLink, 
    useLocation,
    useMatch,
    useResolvedPath,
} from 'react-router-dom'
import type { Location } from "history";
import { Route } from "../../service/router"
import { useState } from "react"

const MENU_ID="sidebar-menu";

type ItemProps = {
    title?: string;
    displayTitle: boolean;
    icon?: IconDefinition;
    route: Route;
    menuIsOpen: boolean;
    childToggleCallback: React.MouseEventHandler<HTMLButtonElement>;
    childMenuIsOpen: boolean;
}

function activeTopLevelMenuItem(routes: Route[], p: Location): string | null {
    const trimmedPath = p.pathname.startsWith("/") ? p.pathname.substring(1) : p.pathname;
    const parts = trimmedPath.split("/")

    if (parts.length == 0) {
        return null
    }

    for (const r of routes) {
        if (`/${parts[0]}` === r.path) {
            return r.name
        }
    }


    return null
}

function MenuNavLink({route, menuIsOpen, childMenuIsOpen, displayTitle, toggle, pathOverride, hideIcon}: {hideIcon?: boolean, childMenuIsOpen?: boolean, route: Route, pathOverride?: string, menuIsOpen: boolean, displayTitle: boolean, toggle?: React.ReactElement}) {
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

function MenuItem({ displayTitle, route, menuIsOpen, childToggleCallback, childMenuIsOpen }: ItemProps) {
    const routesToAdd = route.children.filter((v: Route) => v.showInMenu)
    const hasChildrenToShow = routesToAdd.length > 0
    const childrenToggle = menuIsOpen && hasChildrenToShow ? (
        <span onClick={childToggleCallback} className={styles.itemChildrenToggle}>
            <FontAwesomeIcon icon={childMenuIsOpen ? faChevronUp : faChevronDown} />
        </span>) : <></>

    return (
        <li>
            <MenuNavLink displayTitle={displayTitle} menuIsOpen={menuIsOpen} childMenuIsOpen={childMenuIsOpen} route={route} toggle={childrenToggle} />
            
            {hasChildrenToShow && menuIsOpen ? (
                <ul className={`${styles.itemChildren} ` + (childMenuIsOpen ? `${styles.itemChildrenOpen}` : '')}>
                    {routesToAdd.map((v: Route, i: number) =>
                        <li key={i} >
                            <MenuNavLink 
                                displayTitle={displayTitle} 
                                menuIsOpen={menuIsOpen} 
                                route={v} 
                                pathOverride={`${route.path}/${v.path}`}
                                hideIcon
                            />
                        </li>
                    )}
                </ul>
            ): <></>}
        </li>
    )
}

type Props = {
    isOpen: boolean;
    routes: Route[]
}

function Menu({ isOpen, routes }: Props) {
    const menuStyles = `${styles.menu} ${isOpen ? styles.open : styles.closed}`
    const currentPath = useLocation();
    const [childMenu, setChildMenu] = useState<string | null>(activeTopLevelMenuItem(routes, currentPath));
    const toggleChildMenu = (name: string) => {
        if (name == childMenu) {
            setChildMenu(null)

            return
        }

        setChildMenu(name)
    }
    return (
        <div id={MENU_ID} className={menuStyles}>
            <ul>
                {routes.map((route: Route, i: number) => {
                    return (
                        <MenuItem
                            key={i}
                            menuIsOpen={isOpen}
                            route={route}
                            displayTitle={isOpen}
                            childMenuIsOpen={childMenu == route.name}
                            childToggleCallback={(e: React.MouseEvent) => toggleChildMenu(route.name)}
                        />
                    )
                })}
            </ul>
        </div>
    )
}

export {
    MENU_ID,
}

export default Menu;