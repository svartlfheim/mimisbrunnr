import Header, {MENU_TOGGLE_ID} from '../components/header/Header'
import Content from '../components/content/Content';
import Menu, {MENU_ID} from '../components/menu/Menu';
import React, { useState, useEffect } from 'react';
import {routes} from '../service/router'

type Props = {
    children: React.ReactNode
}

const menuIsOpenKey = "menu-is-open" 

function Base({ children }: Props) {

    const [menuOpen, setMenuIsOpen] = useState(localStorage.getItem(menuIsOpenKey) === "true");
    const toggleMenu = () => {
        console.log('toggling menu to ', !menuOpen)
        const newValue = !menuOpen
        setMenuIsOpen(!menuOpen)
        localStorage.setItem(menuIsOpenKey, newValue ? "true" : "false");
    }
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


    useEffect(() => {
        window.addEventListener("click", hideMenuOnContentClick);
        return () => {
            window.removeEventListener("click", hideMenuOnContentClick);
        };
    }, [hideMenuOnContentClick]);

    return (
        <>
            <Header menuIsOpen={menuOpen} toggleMenuCallback={() => { toggleMenu() }} />
            <Menu isOpen={menuOpen} routes={routes} />
            <Content>
                {children}
            </Content>
        </>
    )
}

export default Base;