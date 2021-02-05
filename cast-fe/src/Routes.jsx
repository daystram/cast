import React from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import {
  Chat,
  Dashboard,
  Fresh,
  Home,
  Liked,
  Live,
  Manage,
  Profile,
  Scene,
  Search,
  Subscribed,
  Trending,
} from "./views";
import { login, logout, callback, authManager } from "./helper/auth";

const Routes = () => {
  return (
    <Switch>
      <Route path={"/"} exact component={Home} />
      <Route path={"/trending"} exact component={Trending} />
      <Route path={"/live"} exact component={Live} />
      <Route path={"/fresh"} exact component={Fresh} />
      <PrivateRoute path={"/liked"} exact component={Liked} />
      <PrivateRoute path={"/subscribed"} exact component={Subscribed} />
      <Route path={"/w/:hash"} exact component={Scene} />
      <Route path={"/c/:hash"} exact component={Chat} />
      <Route path={"/s"} exact component={Search} />
      <PrivateRoute path={"/profile"} exact component={Profile} />
      <PrivateRoute path={"/dashboard"} exact component={Dashboard} />
      <PrivateRoute path={"/manage"} exact component={Manage} />
      <PublicRoute path={"/login"} exact component={login} />
      <PrivateRoute path={"/logout"} exact component={logout} />
      <Route path={"/callback"} exact component={callback} />
      <Route>
        <Redirect to={"/"} />
      </Route>
    </Switch>
  );
};

const PrivateRoute = ({ component: Component, ...rest }) => (
  <Route
    {...rest}
    render={(props) =>
      authManager.isAuthenticated() ? (
        <Component {...props} />
      ) : (
        <Redirect
          to={{ pathname: "/login", state: { from: props.location } }}
        />
      )
    }
  />
);

const PublicRoute = ({ component: Component, ...rest }) => (
  <Route
    {...rest}
    render={(props) =>
      authManager.isAuthenticated() ? (
        <Redirect to={{ pathname: "/", state: { from: props.location } }} />
      ) : (
        <Component {...props} />
      )
    }
  />
);

export default Routes;
