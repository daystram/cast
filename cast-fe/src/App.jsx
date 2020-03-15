import React from 'react';
import './App.css';
import {BrowserRouter} from 'react-router-dom'
import Routes from './Routes'
import Navigation from "./components/Navigation";

function App() {
  return (
    <BrowserRouter>
      <Navigation/>
      <Routes/>
    </BrowserRouter>
  );
}

export default App;
