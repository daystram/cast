import React, { useRef, useState } from "react";
import { Route, useHistory, useLocation } from "react-router-dom";
import {
  Button,
  Col,
  Form,
  FormControl,
  InputGroup,
  Nav,
  Navbar,
  Row,
} from "react-bootstrap";
import MediaQuery from "react-responsive";
import logo from "./logo.svg";
import { authManager } from "../helper/auth";
import { MOBILE_BP } from "../constants/breakpoint";
import Sidebar from "./Sidebar";
import SidebarProfile from "./SidebarProfile";
import { ProfileImage } from "./index";
function Navigation() {
  const [query, setQuery] = useState(
    new URLSearchParams(useLocation().search).get("query") || ""
  );
  const [expanded, setExpanded] = useState(false);
  const inputRef = useRef();
  const history = useHistory();
  const user = authManager.getUser();
  let profileButton = authManager.isAuthenticated() ? (
    <ProfileImage
      size={38}
      name={`${user.given_name} ${user.family_name}`}
      style={style.profile_image}
      onClick={() => {
        history.push("/profile");
      }}
    />
  ) : (
    <>
      <Button
        variant="primary"
        onClick={() => {
          setExpanded(false);
          history.push("/login");
        }}
      >
        Login with{" "}
        <span>
          <i
            className="material-icons"
            style={{
              fontSize: 16,
              lineHeight: "22px",
              verticalAlign: "text-top",
            }}
          >
            lock
          </i>
          Ratify
        </span>
      </Button>
    </>
  );
  let searchBar = (
    <Form
      noValidate
      onSubmit={(e) => {
        e.preventDefault();
        if (inputRef) inputRef.current.blur();
        setExpanded(false);
        if (query.trim()) history.push(`/s?query=${query}`);
      }}
    >
      <InputGroup style={{ flexWrap: "nowrap" }}>
        <FormControl
          type="text"
          placeholder="Search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          ref={inputRef}
        />
        <InputGroup.Append>
          <Button
            variant="outline-primary"
            type="submit"
            style={{ background: "#8B280344", border: "none", marginLeft: 4 }}
          >
            <i className="material-icons">search</i>
          </Button>
        </InputGroup.Append>
      </InputGroup>
    </Form>
  );
  return (
    <Navbar style={style.navbar} sticky="top" expanded={expanded} expand="xl">
      <Row style={{ width: "100%" }} noGutters>
        <Col xs={3}>
          <Navbar.Brand
            onClick={() => {
              setExpanded(false);
              history.push("/");
            }}
          >
            <img
              src={logo}
              height="30"
              className="d-inline-block align-top cast-logo"
              alt="cast"
            />
          </Navbar.Brand>
        </Col>
        <Col style={{ display: "flex", justifyContent: "center" }}>
          <MediaQuery minDeviceWidth={MOBILE_BP}>
            <Nav className="m-auto">{searchBar}</Nav>
          </MediaQuery>
        </Col>
        <Col
          xs={3}
          style={{
            display: "flex",
            justifyContent: "flex-end",
            alignItems: "center",
          }}
        >
          <MediaQuery minDeviceWidth={MOBILE_BP}>{profileButton}</MediaQuery>
          <MediaQuery maxDeviceWidth={MOBILE_BP}>
            <Navbar.Toggle
              aria-controls="responsive-navbar-nav"
              style={{ border: "none", paddingRight: 0 }}
              onClick={() => setExpanded(!expanded)}
            >
              <span
                className="material-icons"
                style={{ color: "#E84409", fontSize: 28, lineHeight: 1 }}
              >
                menu
              </span>
            </Navbar.Toggle>
          </MediaQuery>
        </Col>
        <MediaQuery maxDeviceWidth={MOBILE_BP}>
          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="m-auto">
              <Row style={{ marginTop: 16 }}>
                <Col style={{ marginBottom: 8 }}>{searchBar}</Col>
                <Col xs={"auto"}>{profileButton}</Col>
              </Row>
            </Nav>
            <Route
              path={[
                "/",
                "/trending",
                "/live",
                "/fresh",
                "/liked",
                "/subscribed",
                "/w/:hash",
                "/s",
                "/login",
              ]}
              exact
              render={(props) => (
                <Sidebar {...props} onSelect={() => setExpanded(false)} />
              )}
            />
            <Route
              path={["/profile", "/dashboard", "/manage"]}
              exact
              render={(props) => (
                <SidebarProfile
                  {...props}
                  onSelect={() => setExpanded(false)}
                />
              )}
            />
          </Navbar.Collapse>
        </MediaQuery>
      </Row>
    </Navbar>
  );
}

let style = {
  navbar: {
    padding: "8px 36px",
    background: "#55180287",
    backdropFilter: "blur(8px)",
    marginBottom: 16,
  },
  profile_image: {
    cursor: "pointer",
  },
};

export default Navigation;
