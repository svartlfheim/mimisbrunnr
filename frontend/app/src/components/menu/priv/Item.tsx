import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import { Route } from "../../../Service/router"
import {Link} from './Link'
import type {StyleProps as LinkStyles} from './Link'

type StyleProps = LinkStyles & {
    itemChildrenToggle: string,
    itemChildren: string,
    itemChildrenOpen: string,
}

interface ItemProps {
    title?: string,
    displayTitle: boolean,
    icon?: IconDefinition,
    route: Route,
    menuIsOpen: boolean,
    childToggleCallback: React.MouseEventHandler<HTMLButtonElement>,
    childMenuIsOpen: boolean,
    styles: StyleProps
}

function Item({ displayTitle, route, menuIsOpen, childToggleCallback, childMenuIsOpen, styles }: ItemProps) {
    const routesToAdd = route.children.filter((v: Route) => v.showInMenu)
    const hasChildrenToShow = routesToAdd.length > 0
    const childrenToggle = menuIsOpen && hasChildrenToShow ? (
        <span onClick={childToggleCallback} className={styles.itemChildrenToggle}>
            <FontAwesomeIcon icon={childMenuIsOpen ? faChevronUp : faChevronDown} />
        </span>) : <></>
    const linkStyles: LinkStyles = (({itemTitle, itemIcon, itemLinkWrapper, activeLink}: LinkStyles) => ({itemTitle, itemIcon, itemLinkWrapper, activeLink}))(styles)
    return (
        <li>
            <Link styles={linkStyles} displayTitle={displayTitle} menuIsOpen={menuIsOpen} childMenuIsOpen={childMenuIsOpen} route={route} toggle={childrenToggle} />
            
            {hasChildrenToShow && menuIsOpen ? (
                <ul className={`${styles.itemChildren} ` + (childMenuIsOpen ? `${styles.itemChildrenOpen}` : '')}>
                    {routesToAdd.map((v: Route, i: number) =>
                        <li key={i} >
                            <Link 
                                styles={linkStyles}
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

export {
    Item,
}

export type {
    StyleProps,
}