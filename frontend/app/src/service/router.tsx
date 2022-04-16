import React from "react";
import * as pages from '../Pages'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import { 
    /* eslint-disable-next-line no-useless-rename */
    faHome as faHome,
    faCodeBranch as faIntegrations,
    faDiagramProject as faProjects,
    faCog as faSettings,
    faQuestion as faHelp
 } from '@fortawesome/free-solid-svg-icons'
import type { Location } from "history";


type Route = {
    path: string;
    name: string;
    display?: string;
    element: React.ReactElement;
    icon?: IconDefinition;
    children: Route[];
    showInMenu?: boolean;
    buildBreadcrumbs: (self: Route, part: string, previousCrumbs: Breadcrumb) => Breadcrumb[] | null
}

type breadcrumbBuilder = (routes: Route[], parts: string[], carry: Breadcrumb[]) => Breadcrumb[]

const breadcrumbsFromRoutes: breadcrumbBuilder  = (routes: Route[], parts: string[], carry: Breadcrumb[]) => {
    if (parts.length === 0) {
        return carry;
    }

    const firstPart = parts.shift();

    // Shouldn't really be needed as we bailed out if it was an empty array, but belt & braces
    if (firstPart === undefined) {
        return carry;
    }

    for(const route of routes) {
        if (route.path === "*") {
            const allParts = [firstPart].concat(parts).join("/")

            const bcs = route.buildBreadcrumbs(route, allParts, carry[carry.length - 1]);

            if (bcs === null) {
                return carry
            }

            return carry.concat(bcs)
        } 
       
        const bcs = route.buildBreadcrumbs(route, firstPart, carry[carry.length - 1])

        if (bcs === null) {
            continue;
        }

        carry = carry.concat(bcs)

        return breadcrumbsFromRoutes(route.children, parts, carry)
    }

    return carry
}

const routes: Route[] = [
    {
        path: "/",
        name: "home",
        element: <pages.Home />,
        display: "Home",
        icon: faHome,
        children: [],
        buildBreadcrumbs: (self, part, previous) => null,
    },
    {
        path: "/projects",
        name: "projects",
        element: <pages.ProjectList />,
        display: "Projects",
        icon: faProjects,
        buildBreadcrumbs: (self, part, previous) => {
            if (part === "projects") {
                return [
                    {
                        path: "/projects",
                        title: "Projects",
                    }
                ];
            }

            return null
        },
        children: [
            {
                path: "add",
                name: "projects.add",
                element: <pages.AddProject />,
                display: "+ Add",
                showInMenu: true,
                buildBreadcrumbs: (self, part, previous) => {
                    if (part === "add") {
                        return [
                            {
                                path: `${previous.path}/${part}`,
                                title: "+ Add",
                            }
                        ];
                    }
        
                    return null
                },
                children: [],
            },
            {
                path: ":projectID",
                name: "projects.view",
                element: <pages.ViewProject />,
                buildBreadcrumbs: (self, part, previous) => {
                    if (part !== "") {
                        return [
                            {
                                path: `${previous.path}/${part}`,
                                title: part,
                            }
                        ];
                    }

                    return null
                },
                children: [
                    {
                        path: "docs",
                        name: "projects.view.docs",
                        element: <pages.ViewProjectDocumentation />,
                        buildBreadcrumbs: (self, part, previous) => {
                            if (part === "docs") {
                                return [
                                    {
                                        path: `${previous.path}/${part}`,
                                        title: "Docs",
                                    }
                                ];
                            }
        
                            return null
                        },
                        children: [
                            {
                                path: "*",
                                name: "projects.view.docs.page",
                                element: <pages.ViewProjectDocumentation />,
                                buildBreadcrumbs: (self, part, previous) => {
                                    console.log(part)
                                    const bcs: Breadcrumb[] = []
                                    const parts = part.split("/").filter((s) => s !== '')

                                    if (parts.length === 0) {
                                        return null;
                                    }

                                    let prevPath = previous.path

                                    for (const p of parts) {
                                        const newPath = `${prevPath}/${p}`
                                        bcs.push(
                                            {
                                                path: newPath,
                                                title: p,
                                            }
                                        )
                                        prevPath = newPath;
                                    }

                                    return bcs
                                },
                                children: [],
                            }
                        ],
                    }
                ],
            },
        ],
    },
    {
        path: "/integrations",
        name: "integrations",
        element: <pages.SCMIntegrationList />,
        display: "Integrations",
        icon: faIntegrations,
        children: [
            {
                path: "add",
                name: "integrations.add",
                element: <pages.AddSCMIntegration />,
                display: "+ Add",
                showInMenu: true,
                buildBreadcrumbs: (self, part, previous) => {
                    if (part === "add") {
                        return [
                            {
                                path: `${previous.path}/${part}`,
                                title: "+ Add",
                            }
                        ];
                    }
        
                    return null
                },
                children: [],
            },
        ],
        buildBreadcrumbs: (self, part, previous) => {
            if (part === "integrations") {
                return [
                    {
                        path: "/integrations",
                        title: "Integrations",
                    }
                ]
            }

            return null
        },
    },
    {
        path: "/settings",
        name: "settings",
        element: <pages.SettingsDashboard />,
        display: "Settings",
        icon: faSettings,
        children: [],
        buildBreadcrumbs: (self, part, previous) => {
            if (part === "settings") {
                return [
                    {
                        path: "/settings",
                        title: "Settings",
                    }
                ]
            }

            return null
        },
    },
    {
        path: "/help",
        name: "help",
        element: <pages.HelpDashboard />,
        display: "Help?",
        icon: faHelp,
        children: [
            {
                path: "openapi",
                name: "help.openapi",
                display: "OpenAPI Spec",
                element: <pages.OpenApiSpec />,
                showInMenu: true,
                buildBreadcrumbs: (self, part, previous) => {
                    if (part === "openapi") {
                        return [
                            {
                                path: `${previous.path}/${part}`,
                                title: part,
                            }
                        ];
                    }

                    return null
                },
                children: [],
            }
        ],
        buildBreadcrumbs: (self, part, previous) => {
            if (part === "help") {
                return [
                    {
                        path: "/help",
                        title: "Help",
                    }
                ]
            }

            return null
        },
    },
];

const HomeRoute: Route = routes[0];

interface Breadcrumb {
    path: string,
    title: string,
    icon?: React.ReactElement,
}


const BuildBreadcrumbLinks = (path: Location): Breadcrumb[] => {
    const parts: string[] = path.pathname.split("/")
        .filter((s: string) => s !== '' )
    
    if (parts.length === 0) {
        return [];
    }

    return breadcrumbsFromRoutes(routes, parts, [])
}

export type { 
    Route,
    Breadcrumb
};
export {
    BuildBreadcrumbLinks,
    routes,
    HomeRoute,
}