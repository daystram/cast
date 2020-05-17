import React, {Component} from 'react';
import {Redirect} from 'react-router-dom';
import auth from "../helper/auth";
import notification from "../helper/notification";

class Logout extends Component {
  constructor(props) {
    super(props);
    notification().destroy();
    auth().deauthenticate();
  }
  render() {
    return (
      <Redirect to={"/"}/>
    )
  }
}

export default Logout
