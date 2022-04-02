import Header from '../components/header/Header'
import Content from '../components/content/Content';
import Menu from '../components/menu/Menu';
import React, { useState } from 'react';
import {routes} from '../router'

type Props = {
    children: React.ReactNode
}


function Base({ children }: Props) {
    const [menuOpen, setMenuIsOpen] = useState(false);
    return (
        <>
            <Header menuIsOpen={menuOpen} toggleMenuCallback={() => { setMenuIsOpen(!menuOpen) }} />
            <Menu isOpen={menuOpen} routes={routes} />
            <Content>
                {children}
            </Content>
        </>
    )
}

export default Base;