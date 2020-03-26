import React, {Component} from 'react';
import {Col, Container, Row} from "react-bootstrap";
import SidebarProfile from "./SidebarProfile";

class Dashboard extends Component {
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
    document.title = "Dashboard | cast";
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

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col md={2} sm={12} style={{marginBottom: 32}}>
              <SidebarProfile/>
            </Col>
            <Col md={10} sm={12}>
              <h1 style={style.h1}>Dashboard</h1>
              Dashboard page
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
    padding: "36px 0 0 0",
    },
};

export default Dashboard
