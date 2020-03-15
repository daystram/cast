import React from "react";
import {Switch, Redirect, Route} from 'react-router-dom'
import {Home, Scene} from "./components";
import Login from "./components/Login";
import SignUp from "./components/SignUp";

const Routes = () => {
  return (
    <Switch>
      <Route path={"/"} exact component={Home}/>
      <Route path="/w/:_id" exact component={Scene}/>
      <Route path="/login" exact component={Login}/>
      <Route path="/signup" exact component={SignUp}/>
      <Route><Redirect to="/"/></Route>
    </Switch>
  )
};

export default Routes;
