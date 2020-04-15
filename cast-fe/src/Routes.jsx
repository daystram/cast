import React from "react";
import {Redirect, Route, Switch} from 'react-router-dom'
import {Dashboard, Home, LogIn, LogOut, Profile, Scene, Search, SignUp, Verify} from "./components";
import Manage from "./components/Manage";
import auth from "./helper/auth";

const Routes = () => {
  return (
    <Switch>
      <Route path={"/"} exact component={Home}/>
      <Route path="/w/:hash" exact component={Scene}/>
      <Route path="/s" exact component={Search}/>
      <Route path="/logout" component={LogOut}/>
      <PrivateRoute path="/profile" exact component={Profile}/>
      <PrivateRoute path="/dashboard" exact component={Dashboard}/>
      <PrivateRoute path="/manage" exact component={Manage}/>
      <PublicRoute path="/verify" exact component={Verify}/>
      <PublicRoute path="/login" exact component={LogIn}/>
      <PublicRoute path="/signup" exact component={SignUp}/>
      <Route><Redirect to="/"/></Route>
    </Switch>
  )
};

const PrivateRoute = ({component: Component, ...rest}) => (
  <Route {...rest} render={(props) => (
    auth().is_authenticated()
      ? <Component {...props} />
      : <Redirect to={{pathname: '/login', state: {from: props.location}}}/>
  )}/>
);

const PublicRoute = ({component: Component, ...rest}) => (
  <Route {...rest} render={(props) => (
    auth().is_authenticated()
      ? <Redirect to={{pathname: '/', state: {from: props.location}}}/>
      : <Component {...props} />
  )}/>
);

export default Routes;
