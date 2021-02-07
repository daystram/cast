import React, { Component } from "react";
import { Nav } from "react-bootstrap";
import { NavLink } from "react-router-dom";
import { authManager } from "../helper/auth";

class Sidebar extends Component {
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
          to="/trending"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">trending_up</i> Trending Casts
        </NavLink>
        <NavLink
          to="/live"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">rss_feed</i> Live Casts
        </NavLink>
        <NavLink
          to="/fresh"
          style={style.nav_item}
          activeStyle={style.nav_item_current}
        >
          <i className="material-icons">whatshot</i> Fresh Casts
        </NavLink>
        {authManager.isAuthenticated() && (
          <>
            <hr style={{ width: "100%", margin: "6px 0" }} />
            <NavLink
              to="/liked"
              style={style.nav_item}
              activeStyle={style.nav_item_current}
            >
              <i className="material-icons">thumb_up</i> Liked Casts
            </NavLink>
            <NavLink
              to="/subscribed"
              style={style.nav_item}
              activeStyle={style.nav_item_current}
            >
              <i className="material-icons">video_library</i> Subscribed
            </NavLink>
          </>
        )}
      </Nav>
    );
  }
}

let style = {
  nav_item: {
    padding: "8px 6px",
    fontFamily: "Open Sans",
    color: "#FFFFFFBB",
    textDecoration: "none",
    whiteSpace: "nowrap",
  },
  nav_item_current: {
    color: "#E84409",
    textDecoration: "none",
    whiteSpace: "nowrap",
  },
};

export default Sidebar;
