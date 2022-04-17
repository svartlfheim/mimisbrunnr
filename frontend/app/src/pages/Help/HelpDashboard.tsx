import styles from './HelpDashboard.module.css'
import { Base } from '../Common'
import '../Common/grid.css'
import { Panel } from '../../Components/Layout'
import {IconButton, Type} from '../../Components/Button'
import { Icon } from '@iconify/react';
import { Link } from 'react-router-dom';

interface ActionParams {
    to: string,
    type: Type,
    icon: React.ReactElement,
    title: string,
}

function Action({to, type, icon, title}: ActionParams) {
    return (
        <Panel className={styles.action}>
            <Link to={to}>
                <IconButton
                    type={type}
                    icon={icon}
                    title={title}
                />
            </Link>
        </Panel>
    )
}

function HelpDashboard() {
    return (
        <Base>
            <Action title="Open API Spec" type={Type.Neutral} to="/help/openapi" icon={<Icon width="100px" icon="logos:swagger" />} />
        </Base>
    )
}

export {
    HelpDashboard,
}