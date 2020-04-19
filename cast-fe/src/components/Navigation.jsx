import React, {useRef, useState} from 'react';
import {useHistory, useLocation} from 'react-router-dom';
import {Button, Form, FormControl, Image, InputGroup, Nav, Navbar} from "react-bootstrap";
import logo from './logo.svg'
import urls from "../helper/url";
import auth from "../helper/auth";

function Navigation() {
  const [query, setQuery] = useState(new URLSearchParams(useLocation().search).get("query") || "");
  const inputRef = useRef();
  const history = useHistory();
  return (
    <Navbar style={style.navbar} sticky="top">
      <Navbar.Brand onClick={() => history.push("/")}>
        <img
          src={logo}
          height="30"
          className="d-inline-block align-top cast-logo"
          alt="cast"
        /></Navbar.Brand>
      <Nav className="m-auto">
        <Form noValidate onSubmit={(e) => {
          e.preventDefault();
          inputRef.current.blur();
          if (query.trim()) history.push(`/s?query=${query}`);
        }}>
          <InputGroup>
            <FormControl type="text" placeholder="Search" value={query} onChange={e => setQuery(e.target.value)}
                         ref={inputRef}/>
            <InputGroup.Append>
              <Button variant="outline-primary" type="submit"><i className="material-icons">search</i></Button>
            </InputGroup.Append>
          </InputGroup>
        </Form>
      </Nav>
      {auth().username() ?
        <Image src={urls().profile(auth().username())}
               height={42} width={42}
               style={style.profile_image} onClick={() => history.push("/profile")} roundedCircle/> :
        <>
          <Button variant="outline-primary" size="sm" style={{marginRight: 8}} onClick={() => history.push("/login")}>
            Log In</Button>
          <Button variant="primary" size="sm" onClick={() => history.push("/signup")}>Sign Up</Button>
        </>
      }
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
    background: "white",
    marginBottom: 16
  },
  profile_image: {
    cursor: "pointer"
  }
};

export default Navigation
