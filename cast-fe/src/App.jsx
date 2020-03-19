import React from 'react';
import './App.css';
import {BrowserRouter} from 'react-router-dom'
import Routes from './Routes'
import Navigation from "./components/Navigation";
import axios from 'axios';

axios.defaults.baseURL = "http://localhost:8080/api/v1";
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem("auth_token");
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  error => Promise.reject(error)
);
axios.interceptors.response.use(
  response => {
    const token = response.headers.authorization; // must be lowercase
    if (token) {
      localStorage.setItem("auth_token", token);
    }
    return response;
  },
  error => Promise.reject(error)
);

function App() {
  return (
    <BrowserRouter>
      <Navigation/>
      <Routes/>
    </BrowserRouter>
  );
}

export default App;
