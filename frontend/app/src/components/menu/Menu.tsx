import styles from './Menu.module.css'
import { useLocation } from 'react-router-dom'
import type { Location } from "history";
import { Route } from "../../Service/router"
import { useState } from "react"
import { Item } from './priv'
import type {StyleProps as ItemStyles} from './priv'

const MENU_ID = "sidebar-menu";

function activeTopLevelMenuItem(routes: Route[], p: Location): string | null {
    const trimmedPath = p.pathname.startsWith("/") ? p.pathname.substring(1) : p.pathname;
    const parts = trimmedPath.split("/")

    if (parts.length === 0) {
        return null
    }

    for (const r of routes) {
        if (`/${parts[0]}` === r.path) {
            return r.name
        }
    }


    return null
}

function extractItemStyles({ 
    itemChildrenToggle, 
    itemChildren, 
    itemChildrenOpen, 
    itemTitle, 
    itemIcon, 
    itemLinkWrapper, 
    activeLink,
    itemChildrenHoverList
}: { readonly [key: string]: string }): ItemStyles {
    return { 
        itemChildrenToggle,
        itemChildren,
        itemChildrenOpen,
        itemTitle,
        itemIcon,
        itemLinkWrapper,
        activeLink,
        itemChildrenHoverList
    }
}

type Props = {
    isOpen: boolean;
    routes: Route[]
}

function Menu({ isOpen, routes }: Props) {
    const menuStyles = `${styles.menu} ${isOpen ? styles.open : styles.closed}`
    const itemStyles: ItemStyles = extractItemStyles(styles)
    const currentPath = useLocation();
    const [childMenu, setChildMenu] = useState<string | null>(activeTopLevelMenuItem(routes, currentPath));
    const toggleChildMenu = (name: string) => {
        if (name === childMenu) {
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
                        <Item
                            styles={itemStyles}
                            key={i}
                            menuIsOpen={isOpen}
                            route={route}
                            displayTitle={isOpen}
                            childMenuIsOpen={childMenu === route.name}
                            childToggleCallback={(e: React.MouseEvent) => toggleChildMenu(route.name)}
                        />
                    )
                })}
            </ul>
        </div>
    )
}

export {
    Menu,
    MENU_ID,
}

export default Menu;