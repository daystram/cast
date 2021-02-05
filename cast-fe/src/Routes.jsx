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
} from "./components";
import auth from "./helper/auth";

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
      <PublicRoute path={"/login"} exact component={auth().login} />
      <PrivateRoute path={"/logout"} exact component={auth().logout} />
      <Route path={"/callback"} exact component={auth().callback} />
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
      auth().is_authenticated() ? (
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
      auth().is_authenticated() ? (
        <Redirect to={{ pathname: "/", state: { from: props.location } }} />
      ) : (
        <Component {...props} />
      )
    }
  />
);

export default Routes;
