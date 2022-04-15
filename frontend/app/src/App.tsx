import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route as ReactRoute,
  Link
} from "react-router-dom";
import {routes, Route} from './router'
import NotFound from './pages/NotFound'



import Home from './pages/Home'
import Projects from './pages/Projects';
import Settings from './pages/Settings';
import OpenApiSpec from './pages/OpenApiSpec';
import Integrations from './pages/SCMIntegrations';

function buildRoutes(route: Route, key: string) {
  return (
    <ReactRoute key={key} path={route.path} element={route.element}>
      {route.children.map((r: Route, j: number) => buildRoutes(r, `${key}-${j}`))}
    </ReactRoute>
  )
}

function _buildRoutes(r: Route, key: string, pathPrefix: string, carry: Array<React.ReactElement>) {
  const path = `${pathPrefix}/${r.path}`
  carry.push(<ReactRoute key={key} path={path} element={r.element}  />)

  r.children.forEach((cR: Route, i: number) => {
    _buildRoutes(cR, `${key}-${i}`, path, carry)
  })
}

function App() {
  const routeElements: Array<React.ReactElement> = []
  routes.forEach((r: Route, i: number) => _buildRoutes(r, `${i}`, "", routeElements))

  return (
    <Router>
      <Routes>
        {routeElements}
        <ReactRoute path="*" element={<NotFound/>} />
      </Routes>
  </Router>
  );
}

export default App;
