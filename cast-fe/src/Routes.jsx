import React from "react";
import {Redirect, Route, Switch} from 'react-router-dom'
import {Home, LogOut, Scene} from "./components";
import Login from "./components/LogIn";
import SignUp from "./components/SignUp";

const Routes = () => {
  return (
    <Switch>
      <Route path={"/"} exact component={Home}/>
      <Route path="/w/:hash" exact component={Scene}/>
      <Route path="/logout" component={LogOut}/>
      {/*<PrivateRoute path="/profile" exact component={Profile}/>*/}
      <PublicRoute path="/login" exact component={Login}/>
      <PublicRoute path="/signup" exact component={SignUp}/>
      <Route><Redirect to="/"/></Route>
    </Switch>
  )
};

const PrivateRoute = ({component: Component, ...rest}) => (
  <Route {...rest} render={(props) => (
    localStorage.getItem("username")
      ? <Component {...props} />
      : <Redirect to={{pathname: '/login', state: {from: props.location}}}/>
  )}/>
);

const PublicRoute = ({component: Component, ...rest}) => (
  <Route {...rest} render={(props) => (
    localStorage.getItem("username")
      ? <Redirect to={{pathname: '/', state: {from: props.location}}}/>
      : <Component {...props} />
  )}/>
);

export default Routes;
