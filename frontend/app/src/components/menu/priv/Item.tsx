import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faChevronDown, faChevronUp } from '@fortawesome/free-solid-svg-icons'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import { Route } from "../../../Service/router"
import {Link} from './Link'
import type {StyleProps as LinkStyles} from './Link'
import {useCallback, useState, useRef} from 'react'

type StyleProps = LinkStyles & {
    itemChildrenToggle: string,
    itemChildren: string,
    itemChildrenOpen: string,
    itemChildrenHoverList: string,
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

interface HoverMenuState {
    x: number,
    y: number,
    visible: boolean,
}

function Item({ displayTitle, route, menuIsOpen, childToggleCallback, childMenuIsOpen, styles }: ItemProps) {
    const linkStyles: LinkStyles = (({itemTitle, itemIcon, itemLinkWrapper, activeLink}: LinkStyles) => ({itemTitle, itemIcon, itemLinkWrapper, activeLink}))(styles)
    const routesToAdd = route.children.filter((v: Route) => v.showInMenu)
    const hasChildrenToShow = routesToAdd.length > 0

    const hoverRef = useRef<HTMLLIElement>(null)
    const [hoverMenuState, setHoverMenuState] = useState<HoverMenuState>({x: 0, y: 0, visible: false})
    const showHoverMenu = useCallback(() => {
        if (! hoverMenuState.visible) {
            const pos =  hoverRef.current?.getClientRects()[0]
            const left =  (pos?.left || 0) + (pos?.width || 0)
            const newState = {
                x: left, 
                y: pos?.top || 0, 
                visible: true,
            }
            setHoverMenuState(newState)
        }
    }, [hoverMenuState, setHoverMenuState, hoverRef])

    const hideHoverMenu = useCallback(() => {
        if (hoverMenuState.visible) {
            const newState = {
                x: hoverMenuState.x, 
                y: hoverMenuState.y, 
                visible: false,
            }
            setHoverMenuState(newState)
        }
    }, [hoverMenuState, setHoverMenuState])

    const hoverMenu = (
        <ul className={styles.itemChildrenHoverList} style={{top: `${hoverMenuState.y}px`, left: `${hoverMenuState.x}px`}}>
            <li>
                <Link 
                    
                    styles={linkStyles} 
                    displayTitle={true} 
                    menuIsOpen={true} 
                    childMenuIsOpen={true} 
                    route={route} 
                    hideIcon
                />
            </li>
            {routesToAdd.map((v: Route, i: number) =>
                <li key={i} >
                    <Link 
                        styles={linkStyles}
                        displayTitle={true} 
                        menuIsOpen={menuIsOpen} 
                        route={v} 
                        pathOverride={`${route.path}/${v.path}`}
                        hideIcon
                    />
                </li>
            )}
        </ul>
    );

    const childrenToggle = menuIsOpen && hasChildrenToShow ? (
        <span onClick={childToggleCallback} className={styles.itemChildrenToggle}>
            <FontAwesomeIcon icon={childMenuIsOpen ? faChevronUp : faChevronDown} />
        </span>) : <></>
    

    return (
        <li 
            ref={hoverRef} 
            onMouseEnter={hasChildrenToShow && !menuIsOpen ? showHoverMenu : undefined} 
            onMouseLeave={hasChildrenToShow && !menuIsOpen ? hideHoverMenu : undefined} 
        >
            <Link 
                
                styles={linkStyles} 
                displayTitle={displayTitle} 
                menuIsOpen={menuIsOpen} 
                childMenuIsOpen={childMenuIsOpen} 
                route={route} 
                toggle={childrenToggle} 
            />
                {hoverMenuState.visible ? hoverMenu : <></>}
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