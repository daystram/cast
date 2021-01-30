import React, { Component } from "react";
import { Nav } from "react-bootstrap";
import { NavLink } from "react-router-dom";

class SidebarProfile extends Component {
  render() {
    return (
      <Nav className="flex-column" onClick={this.props.onSelect}>
        <NavLink
          to="/"
          exact
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">home</i> Home
        </NavLink>
        <hr style={{ width: "100%", margin: "6px 0" }} />
        <NavLink
          to="/profile"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">person</i> Profile
        </NavLink>
        <NavLink
          to="/dashboard"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">dashboard</i> Dashboard
        </NavLink>
        <NavLink
          to="/manage"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">settings</i> Manage
        </NavLink>
        <hr style={{ width: "100%", margin: "6px 0" }} />
        <NavLink
          to="/logout"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">exit_to_app</i> Logout
        </NavLink>
      </Nav>
    );
  }
}

let style = {
  nav_item: {
    padding: "8px 6px",
    fontFamily: "Open Sans",
    fontSize: 18,
    color: "#FFFFFFBB",
    textDecoration: "none",
    whiteSpace: "nowrap",
  },
  nav_item_current: {
    color: "#E84409",
    fontSize: 18,
    textDecoration: "none",
    whiteSpace: "nowrap",
  },
};

export default SidebarProfile;
