import React from 'react';
import './App.css';
import {BrowserRouter, Route} from 'react-router-dom'
import Routes from './Routes'
import Navigation from "./components/Navigation";
import {Slide, ToastContainer} from 'react-toastify';
import axios from 'axios';
import notification from "./helper/notification";

// axios.defaults.baseURL = "http://localhost:8080";
// axios.defaults.baseURL = "https://dev.cast.daystram.com";
axios.defaults.withCredentials = true;

function App() {
  notification().init();
  return (
    <>
      <BrowserRouter>
        <Route path={['/', '/trending', '/live', '/fresh', '/liked', '/subscribed', '/w/:hash',
          '/s', '/profile', '/dashboard', '/manage', '/verify', '/login', '/signup', '/forget']} exact
               component={Navigation}/>
        <Routes/>
        <ToastContainer
          position={"top-right"}
          autoClose={20000}
          hideProgressBar={false}
          transition={Slide}
          closeOnClick
          closeButton={false}
          rtl={false}
          pauseOnFocusLoss
          draggable
          pauseOnHover/>
      </BrowserRouter>
    </>
  );
}

export default App;
