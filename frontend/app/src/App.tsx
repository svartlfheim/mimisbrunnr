import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route as ReactRoute,
  Link
} from "react-router-dom";
import {routes, Route} from './Service/router'
import {NotFound} from './Pages'


function buildRoutes(r: Route, key: string, pathPrefix: string, carry: Array<React.ReactElement>) {
  const path = `${pathPrefix}/${r.path}`
  carry.push(<ReactRoute key={key} path={path} element={r.element}  />)

  r.children.forEach((cR: Route, i: number) => {
    buildRoutes(cR, `${key}-${i}`, path, carry)
  })
}

function App() {
  const routeElements: Array<React.ReactElement> = []
  routes.forEach((r: Route, i: number) => buildRoutes(r, `${i}`, "", routeElements))

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
