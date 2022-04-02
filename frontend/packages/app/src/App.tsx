import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route as ReactRoute,
  Link
} from "react-router-dom";
import {routes, Route} from './router'
import NotFound from './pages/NotFound'


function App() {
  return (
    <Router>
      <Routes>
        {routes.map((route: Route, i: number) => {
          return (
            <ReactRoute key={i} path={route.path} element={route.element} />
          );
        })}
        <ReactRoute path="*" element={<NotFound/>} />
        {/* <Route path="/projects" element={<Projects />} />
        <Route path="/scm-integrations" element={<SCMIntegrations />} />
        <Route path="/settings" element={<Settings />} />
        <Route path="/" element={<Home />} /> */}
      </Routes>
  </Router>
      // <Home />

  );
}

export default App;
