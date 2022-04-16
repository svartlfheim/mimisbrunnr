import Header, {MENU_TOGGLE_ID} from '../../Components/Header'
import Content from '../../Components/Content';
import Menu, {MENU_ID} from '../../Components/Menu';
import React, { useState, useEffect, useCallback } from 'react';
import {routes} from '../../Service/router'
import {Mode} from '../../Components/Content'
import {Footer} from '../../Components/Footer'

type Props = {
    children: React.ReactNode,
    gridMode?: Mode
}

const menuIsOpenKey = "menu-is-open" 

function Base({ children, gridMode }: Props) {
    gridMode = gridMode === undefined ? Mode.ResponsiveColumns : gridMode
    const [menuOpen, setMenuIsOpen] = useState(localStorage.getItem(menuIsOpenKey) === "true");
    const toggleMenu = useCallback(() => {
        const newValue = !menuOpen
        setMenuIsOpen(!menuOpen)
        localStorage.setItem(menuIsOpenKey, newValue ? "true" : "false");
    }, [menuOpen])

    useEffect(() => {
        const hideMenuOnContentClick = (e: any) => {
            /*
            Oh this is gross...
            But doing something clever doesn't work. Checking whether the elements 
            with the given ids contain (https://www.w3schools.com/jsref/met_node_contains.asp)
            the event target, fails when the user clicks one of the links. I'm 
            guessing that as react is triggering a re-render when they're clicked it
            interferes with the DOM too much for this too work.
    
            However doing it this way, by traversing up the dom, it works...
    
            I mean it's kind of an infinite loop, but there's always the HTML root, 
            it just seems fucking nasty man...
            */
            let element = e.target
            while (element !== null) {
                if ([MENU_ID, MENU_TOGGLE_ID].includes(element.id)) {
                    return true
                }
    
                element = element.parentElement
            }
    
            if (menuOpen) {
                toggleMenu()
            }
        }

        window.addEventListener("click", hideMenuOnContentClick);
        return () => {
            window.removeEventListener("click", hideMenuOnContentClick);
        };
    }, [toggleMenu, menuOpen]);

    return (
        <>
            <Header menuIsOpen={menuOpen} toggleMenuCallback={() => { toggleMenu() }} />
            <Menu isOpen={menuOpen} routes={routes} />
            <Content mode={gridMode}>
                {children}
            </Content>
            <Footer />
        </>
    )
}

export {
    Base
};