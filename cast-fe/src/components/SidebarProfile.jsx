import React, {Component} from 'react';
import {Card, Nav} from "react-bootstrap";
import {NavLink} from "react-router-dom";

class SidebarProfile extends Component {
  render() {
    return (
      <Card body style={style.sidebar}>
        <Nav className="flex-column">
          <NavLink to="/" exact style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">home</i> Home</NavLink>
          <NavLink to="/profile" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">person</i> Profile</NavLink>
          <NavLink to="/dashboard" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">dashboard</i> Dashboard</NavLink>
          <NavLink to="/manage" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">settings</i> Manage</NavLink>
          <NavLink to="/logout" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">exit_to_app</i> Logout</NavLink>
        </Nav>
      </Card>
    )
  }
}

let style = {
  sidebar: {
    borderRadius: "8px 48px 8px 8px"
  },
  nav_item: {
    padding: "8px 6px",
    fontFamily: "Open Sans",
    fontSize: 18,
    color: "black",
    textDecoration: "none"
  },
  nav_item_current: {
    color: "#E84409",
    fontSize: 18,
    textDecoration: "none"
  }
};

export default SidebarProfile
