import React from 'react';
import './App.css';

import {Route, Switch, BrowserRouter as Router, Link} from "react-router-dom";
import Navbar from "./components/navbar/navbar";
import Content from "./components/content/content";

function App() {


  return (
    <div className="App">
      <Router>
          <Navbar/>
          <Content/>
      </Router>
    </div>
  );
}

type AppState = {
    route: string,
}


export default App;
