import React, {Component} from 'react';
import {Alert, Button, Container, Form, Spinner} from "react-bootstrap";
import axios from "axios";
import urls from "../helper/url";

class LogIn extends Component {
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
    axios.post(urls().login(), {
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
          <h1 style={style.h1}>Log In</h1>
          {this.state.error_login && <Alert variant={"danger"}>{this.state.error_login}</Alert>}
          <Form noValidate onSubmit={this.submitForm}>
            <Form.Group>
              <Form.Label>Username</Form.Label>
              <Form.Control name={"username"} value={this.state.username} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"name"}
                            isInvalid={!!this.state.error_username}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_username}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group>
              <Form.Label>Password</Form.Label>
              <Form.Control name={"password"} value={this.state.password} onBlur={this.handleChange}
                            onChange={this.handleChange}
                            type={"password"} isInvalid={!!this.state.error_password}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_password}</Form.Control.Feedback>
            </Form.Group>
            <Button variant="primary" type="submit" block disabled={this.state.loading}>
              Log In{" "}
              {this.state.loading &&
              <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                       size="sm" role="status" aria-hidden="true"/>}
            </Button>
          </Form>
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
    maxWidth: 480
  },
};

export default LogIn
