import React, {Component} from 'react';
import {Redirect} from 'react-router-dom';
import auth from "../helper/auth";

class Logout extends Component {
  constructor(props) {
    super(props);
    auth().deauthenticate();
  }
  render() {
    return (
      <Redirect to={"/"}/>
    )
  }
}

export default Logout
