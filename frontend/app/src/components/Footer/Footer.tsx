import styles from './Footer.module.css'
import { BuildBreadcrumbLinks, HomeRoute } from '../../Service/router';
import type {Breadcrumb} from '../../Service/router'
import { Link, useLocation } from 'react-router-dom'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';


function Footer() {
    const currentPath = useLocation();
    const bcLinks = BuildBreadcrumbLinks(currentPath)
    bcLinks.unshift({
        path: HomeRoute.path,
        title: HomeRoute.display ?? "home",
        icon: HomeRoute.icon !== undefined 
            ? <FontAwesomeIcon icon={HomeRoute.icon} />
            : <span>home</span>,
    })

    return (
        <div className={styles.footerWrapper}>
            <div className={styles.breadcrumbs}>
                <ul>
                    {bcLinks.map((l: Breadcrumb, index: number) => {
                        const elements = [
                            (<li key={index}><Link to={l.path}>{l.icon !== undefined ? l.icon : l.title}</Link></li>),
                        ];
                        if (index < bcLinks.length - 1) {
                            elements.push(
                                (<li key={`sep-${index}`} className={styles.breadcrumbSeparator}>/</li>)
                            );
                        }
                        return elements;
                    })}
                </ul>
            </div>
        </div>
    )
}

export {
    Footer,
}