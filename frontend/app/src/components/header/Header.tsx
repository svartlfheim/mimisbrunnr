import styles from './Header.module.css'
import logo from './logo-white.png'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBars as faMenuOpen, faClose as faMenuClose } from '@fortawesome/free-solid-svg-icons'
import { BuildBreadcrumbLinks } from '../../Service/router';
import { Link, useLocation } from 'react-router-dom'

const MENU_TOGGLE_ID = "sidebar-menu-toggler"

type Props = {
  menuIsOpen: boolean;
  toggleMenuCallback: () => void
}

function Header({ menuIsOpen, toggleMenuCallback }: Props) {

  const menuIcon = menuIsOpen ?
    (<FontAwesomeIcon icon={faMenuClose} onClick={() => { toggleMenuCallback() }} />) :
    (<FontAwesomeIcon icon={faMenuOpen} onClick={() => { toggleMenuCallback() }} />);

  const currentPath = useLocation();
  const bcLinks = BuildBreadcrumbLinks(currentPath)

  return (
    <div className={styles.headerWrapper}>
      <div className={styles.header}>
        <div id={MENU_TOGGLE_ID} className={styles.menuIconWrapper}>
          <div className={styles.menuIcon}>

            {menuIcon}
          </div>
        </div>
        <div className={styles.logo}>
          <img className={styles.logoImg} src={logo} alt="Logo" />
        </div>
      </div>
      {/* <div className={styles.breadcrumbs}>
        <ul>
          {bcLinks.map((l: {
            path: string;
            title: string;
          }, index: number) => {
            const elements = [
              (<li key={index}><Link to={l.path}>{l.title}</Link></li>),
            ];
            if (index < bcLinks.length - 1) {
              elements.push(
                (<li key={`sep-${index}`} className={styles.breadcrumbSeparator}>/</li>)
              );
            }
            return elements;
          })}
        </ul>
      </div> */}
    </div>
  )
}

export {
  MENU_TOGGLE_ID,
  Header,
}

export default Header;