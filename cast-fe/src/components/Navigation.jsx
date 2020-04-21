import React, {useRef, useState} from 'react';
import {Route, useHistory, useLocation} from 'react-router-dom';
import {Button, Col, Form, FormControl, Image, InputGroup, Nav, Navbar, Row} from "react-bootstrap";
import MediaQuery from 'react-responsive'
import logo from './logo.svg'
import urls from "../helper/url";
import auth from "../helper/auth";
import {MOBILE_BP} from "../constants/breakpoint";
import Sidebar from "./Sidebar";
import SidebarProfile from "./SidebarProfile";

function Navigation() {
  const [query, setQuery] = useState(new URLSearchParams(useLocation().search).get("query") || "");
  const [expanded, setExpanded] = useState(false);
  const inputRef = useRef();
  const history = useHistory();
  let profileButton = (
    auth().username() ?
      <Image src={urls().profile(auth().username())}
             height={38} width={38}
             style={style.profile_image} onClick={() => {
        history.push("/profile")
      }} roundedCircle/> :
      <>
        <Button variant="outline-primary" style={{marginRight: 8}}
                onClick={() => {
                  setExpanded(false);
                  history.push("/login")
                }}>
          Log In</Button>
        <Button variant="primary" onClick={() => {
          setExpanded(false);
          history.push("/signup")
        }}>Sign Up</Button>
      </>
  );
  let searchBar = (
    <Form noValidate onSubmit={(e) => {
      e.preventDefault();
      if (inputRef) inputRef.current.blur();
      setExpanded(false);
      if (query.trim()) history.push(`/s?query=${query}`);
    }}>
      <InputGroup style={{flexWrap: "nowrap"}}>
        <FormControl type="text" placeholder="Search" value={query}
                     onChange={e => setQuery(e.target.value)}
                     ref={inputRef}/>
        <InputGroup.Append>
          <Button variant="outline-primary" type="submit"
                  style={{background: "#8B280344", border: "none", marginLeft: 4}}>
            <i className="material-icons">search</i></Button>
        </InputGroup.Append>
      </InputGroup>
    </Form>
  );
  return (
    <Navbar style={style.navbar} sticky="top" expanded={expanded} expand="xl">
      <Row style={{width: "100%"}} noGutters>
        <Col xs={3}>
          <Navbar.Brand onClick={() => {
            setExpanded(false);
            history.push("/")
          }}>
            <img
              src={logo}
              height="30"
              className="d-inline-block align-top cast-logo"
              alt="cast"
            /></Navbar.Brand>
        </Col>
        <Col style={{display: "flex", justifyContent: "center"}}>
          <MediaQuery minDeviceWidth={MOBILE_BP}>
            <Nav className="m-auto">
              {searchBar}
            </Nav>
          </MediaQuery>
        </Col>
        <Col xs={3} style={{display: "flex", justifyContent: "flex-end", alignItems: "center"}}>
          <MediaQuery minDeviceWidth={MOBILE_BP}>
            {profileButton}
          </MediaQuery>
          <MediaQuery maxDeviceWidth={MOBILE_BP}>
            <Navbar.Toggle aria-controls="responsive-navbar-nav" style={{border: "none", paddingRight: 0}}
                           onClick={() => setExpanded(!expanded)}>
              <span class="material-icons" style={{color: "#E84409", fontSize: 28, lineHeight: 1}}>menu</span>
            </Navbar.Toggle>
          </MediaQuery>
        </Col>
        <MediaQuery maxDeviceWidth={MOBILE_BP}>
          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="m-auto">
              <Row style={{marginTop: 16}}>
                <Col style={{marginBottom: 8}}>
                  {searchBar}
                </Col>
                <Col xs={"auto"}>
                  {profileButton}
                </Col>
              </Row>
            </Nav>
            <Route path={['/', '/w/:hash', '/s', '/verify', '/login', '/signup']} exact
                   render={props => <Sidebar {...props} onSelect={() => setExpanded(false)}/>}/>
            <Route path={['/profile', '/dashboard', '/manage']} exact
                   render={props => <SidebarProfile {...props} onSelect={() => setExpanded(false)}/>}/>
          </Navbar.Collapse>
        </MediaQuery>
      </Row>
    </Navbar>
  )
}

let style = {
  navbar: {
    padding: "8px 36px",
    background: "#55180287",
    backdropFilter: "blur(8px)",
    marginBottom: 16
  },
  profile_image: {
    cursor: "pointer"
  }
};

export default Navigation
