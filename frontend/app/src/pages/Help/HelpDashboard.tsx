import { Base } from '../Common'
import '../Common/grid.css'
import { Panel } from '../../Components/Layout'
import {IconButton, Type} from '../../Components/Button'
import { Icon } from '@iconify/react';
import { Link } from 'react-router-dom';

function HelpDashboard() {
    return (
        <Base>
            <Panel>
                <Link to="/help/openapi">
                    <IconButton
                        type={Type.Neutral}
                        icon={<Icon width="100px" icon="logos:swagger" />}
                        title="Open API Spec" 
                    />
                </Link>
            </Panel>
        </Base>
    )
}

export {
    HelpDashboard,
}