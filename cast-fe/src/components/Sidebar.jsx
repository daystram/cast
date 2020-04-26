import React, {Component} from 'react';
import {Nav} from "react-bootstrap";
import {NavLink} from "react-router-dom";

class Sidebar extends Component {
  render() {
    return (
      <Nav className="flex-column" onClick={this.props.onSelect}>
        <NavLink to="/" exact style={style.nav_item} activeStyle={style.nav_item_current}>
          <i className="material-icons">home</i> Home</NavLink>
        <NavLink to="/live" style={style.nav_item} activeStyle={style.nav_item_current}>
          <i className="material-icons">rss_feed</i> Live Casts</NavLink>
        <NavLink to="/subscription" style={style.nav_item} activeStyle={style.nav_item_current}>
          <i className="material-icons">video_library</i> Subscription</NavLink>
        <NavLink to="/bookmarked" style={style.nav_item} activeStyle={style.nav_item_current}>
          <i className="material-icons">playlist_play</i> Bookmarked</NavLink>
      </Nav>
    )
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
  }
};

export default Sidebar
