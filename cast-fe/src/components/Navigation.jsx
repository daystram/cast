import React from 'react';
import {useHistory} from 'react-router-dom';
import {Button, FormControl, InputGroup, Nav, Navbar} from "react-bootstrap";
import logo from './logo.svg'

function Navigation() {
  const history = useHistory();
  return (
    <Navbar style={style.navbar} sticky="top">
      <Navbar.Brand onClick={() => history.push("/")}>
        <img
          src={logo}
          height="30"
          className="d-inline-block align-top"
          alt="cast"
        /></Navbar.Brand>
      <Nav className="m-auto">
        <InputGroup inline>
          <FormControl type="text" placeholder="Search"/>
          <InputGroup.Append>
            <Button variant="outline-info"><i className="material-icons">search</i></Button>
          </InputGroup.Append>
        </InputGroup>
      </Nav>
      <Button variant="outline-info">Sign Up</Button>
    </Navbar>
  )
}

let style = {
  navbar: {
    borderRadius: "0 0 48px 48px",
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    height: 64,
    padding: "8px 36px",
    background: "white"
  },
};

export default Navigation
