import React, {Component} from 'react';
import {Card, Nav} from "react-bootstrap";
import {NavLink} from "react-router-dom";

class Sidebar extends Component {
  render() {
    return (
      <Card body style={style.sidebar}>
        <Nav className="flex-column">
          <NavLink to="/" exact style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">home</i> Home</NavLink>
          <NavLink to="/live" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">rss_feed</i> Live Casts</NavLink>
          <NavLink to="/subscription" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">video_library</i> Subscription</NavLink>
          <NavLink to="/bookmarked" style={style.nav_item} activeStyle={style.nav_item_current}>
            <i className="material-icons">playlist_play</i> Bookmarked</NavLink>
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
    color: "blue",
    fontSize: 18,
    textDecoration: "none"
  }
};

export default Sidebar
