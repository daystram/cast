import React from "react";
import {Switch, Redirect, Route} from 'react-router-dom'
import {Home, Scene} from "./components";

const Routes = () => {
  return (
    <Switch>
      <Route path={"/"} exact component={Home}/>
      <Route path="/w/:_id" exact component={Scene}/>
      <Route><Redirect to="/"/></Route>
    </Switch>
  )
};

export default Routes;
