import React from 'react';
import './App.css';
import {BrowserRouter, Route} from 'react-router-dom'
import Routes from './Routes'
import Navigation from "./components/Navigation";
import axios from 'axios';

axios.defaults.baseURL = "http://localhost:8080";
// axios.defaults.baseURL = "https://dev.cast.daystram.com";
axios.defaults.withCredentials = true;

function App() {
  return (
    <BrowserRouter>
      <Route path={['/', '/w/:hash', '/s', '/profile', '/dashboard',
        '/manage', '/verify', '/login', '/signup']} exact component={Navigation}/>
      <Routes/>
    </BrowserRouter>
  );
}

export default App;
