import React, {Component} from 'react';
import {Col, Container, Row} from "react-bootstrap";
import axios from "axios";
import SidebarProfile from "./SidebarProfile";

class Profile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: "",
      error_username: "",
      error_password: "",
      error_login: "",
      loading: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  handleChange(e) {
    this.setState({error_login: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    switch (field) {
      case "username":
        if (!value.trim()) {
          this.setState({error_username: "Please enter your username"});
          return false;
        }
        this.setState({error_username: ""});
        return true;
      case "password":
        if (!value) {
          this.setState({error_password: "Please enter your password"});
          return false;
        }
        this.setState({error_password: ""});
        return true;
      default:
        return false;
    }
  }

  submitForm(e) {
    e.preventDefault();
    let ok = true;
    ok &= this.validate("username", this.state.username);
    ok &= this.validate("password", this.state.password);
    if (!ok) return;
    this.setState({loading: true});
    axios.post('/auth/login', {
      username: this.state.username.trim(),
      password: this.state.password.trim(),
    }).then((response) => {
      this.setState({loading: false});
      switch (response.data.code) {
        case 200:
          localStorage.setItem("username", this.state.username.trim());
          this.props.history.push("/");
          return;
        case 404:
          this.setState({error_username: "Username not registered"});
          return;
        case 403:
          this.setState({error_password: "Incorrect password"});
          return;
        default:
          this.setState({error_login: "An error has occurred!"});
          return;
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col md={2} sm={12} style={{marginBottom: 32}}>
              <SidebarProfile/>
            </Col>
            <Col md={10} sm={12}>
              <h1 style={style.h1}>Profile</h1>
              Profile page
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
  },
  content_container: {
    padding: "36px 0 0 0"
  },
};

export default Profile