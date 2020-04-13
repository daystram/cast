import React from 'react';
import './App.css';
import {BrowserRouter} from 'react-router-dom'
import Routes from './Routes'
import Navigation from "./components/Navigation";
import axios from 'axios';

// axios.defaults.baseURL = "http://localhost:8080";
// axios.defaults.baseURL = "https://dev.cast.daystram.com";
axios.defaults.withCredentials = true;

function App() {
  return (
    <BrowserRouter>
      <Navigation/>
      <Routes/>
    </BrowserRouter>
  );
}

export default App;
