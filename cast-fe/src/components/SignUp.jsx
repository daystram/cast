import React, {Component} from 'react';
import {Alert, Button, Col, Container, Form, ProgressBar, Spinner} from "react-bootstrap";
import axios from 'axios';
import zxcvbn from 'zxcvbn';
import urls from "../helper/url";
import {Link} from "react-router-dom";

let timeout = {};

class SignUp extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      username: "",
      email: "",
      password: "",
      password2: "",
      error_name: "",
      error_username: "",
      error_email: "",
      error_password: "",
      error_password2: "",
      error_signup: "",
      loading: false,
      success: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  handleChange(e) {
    this.setState({error_signup: "", success: false});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    if (this.state.success) return;
    switch (field) {
      case "name":
        if (!value.trim()) {
          this.setState({error_name: "Please enter your name"});
          return false;
        }
        if (value.trim().length < 4) {
          this.setState({error_name: "Name too short"});
          return false;
        }
        this.setState({error_name: ""});
        this.checkAvailability(field, value.trim());
        return true;
      case "username":
        if (!value.trim()) {
          this.setState({error_username: "Please enter your username"});
          return false;
        }
        if (value.trim().length < 4) {
          this.setState({error_username: "Username must be at least 4 alphanumerics"});
          return false;
        }
        let usernameRe = /^[A-Za-z0-9]{4,}$/;
        if (!usernameRe.test(value.trim())) {
          this.setState({error_username: "Username must only contain alphanumerics"});
          return false;
        }
        this.setState({error_username: ""});
        this.checkAvailability(field, value.trim());
        return true;
      case "email":
        if (!value.trim()) {
          this.setState({error_email: "Please enter your email"});
          return false;
        }
        let emailRe = /.+@.+\..+/;
        if (!emailRe.test(value.trim())) {
          this.setState({error_email: "Invalid email address"});
          return false;
        }
        this.setState({error_email: ""});
        this.checkAvailability(field, value.trim());
        return true;
      case "password":
        if (!value) {
          this.setState({error_password: "Please enter your password"});
          return false;
        }
        this.setState({strength: zxcvbn(this.state.password).score});
        if (value.length < 8) {
          this.setState({error_password: "Password must be at least 8 characters"});
          return false;
        }
        this.setState({error_password: ""});
        return true;
      case "password2":
        if (!value) {
          this.setState({error_password2: "Please re-enter your password"});
          return false;
        }
        if (value !== this.state.password) {
          this.setState({error_password2: "Passwords do not match"});
          return false;
        }
        this.setState({error_password2: ""});
        return true;
      default:
        return false;
    }
  }

  checkAvailability(field, value) {
    clearTimeout(timeout[field]);
    timeout[field] = setTimeout(() => {
      axios.post(urls().auth_check(), {
        field: field.trim(),
        value: value.trim()
      }).then((response) => {
        if (this.state.success) return;
        if (response.data.code !== 200) {
          this.setState({[`error_${field}`]: response.data.error});
        } else {
          this.setState({[`error_${field}`]: ""});
        }
      }).catch((error) => {
        console.log(error);
        this.setState({error_signup: "An error has occurred!"});
      });
    }, 400)
  }

  submitForm(e) {
    e.preventDefault();
    let ok = true;
    if (!this.state.attempted) {
      this.setState({attempted: true});
      ok &= this.validate("name", this.state.name);
      ok &= this.validate("username", this.state.username);
      ok &= this.validate("email", this.state.email);
      ok &= this.validate("password", this.state.password);
      ok &= this.validate("password", this.state.password2);
    } else {
      ok &= !this.state.error_name;
      ok &= !this.state.error_username;
      ok &= !this.state.error_email;
      ok &= !this.state.error_password;
      ok &= !this.state.error_password2;
    }
    if (!ok) return;
    this.setState({loading: true});
    axios.post(urls().signup(), {
      name: this.state.name.trim(),
      username: this.state.username.trim(),
      email: this.state.email.trim(),
      password: this.state.password.trim(),
    }).then((response) => {
      this.setState({loading: false, name: "", username: "", email: "", password: "", password2: ""});
      if (response.data.code === 200) {
        Object.keys(timeout).map(field => clearTimeout(timeout[field]));
        this.setState({success: true});
      } else {
        this.setState({error_signup: "An error has occurred!"});
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_signup: "An error has occurred!"});
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <h1 style={style.h1}>Sign Up</h1>
          {this.state.success &&
          <Alert variant="success">
            <Alert.Heading>Welcome!</Alert.Heading>
            <p>
              Your account has been successfully registered! Check your email to verify your email before logging in.
            </p>
          </Alert>}
          {this.state.error_signup && <Alert variant={"danger"}>{this.state.error_signup}</Alert>}
          <Form noValidate onSubmit={this.submitForm}>
            <Form.Group>
              <Form.Label>Name</Form.Label>
              <Form.Control name={"name"} value={this.state.name} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"name"}
                            isInvalid={!!this.state.error_name}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_name}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group>
              <Form.Label>Username</Form.Label>
              <Form.Control name={"username"} value={this.state.username} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"name"}
                            isInvalid={!!this.state.error_username}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_username}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group>
              <Form.Label>Email</Form.Label>
              <Form.Control name={"email"} value={this.state.email} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"email"}
                            isInvalid={!!this.state.error_email}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_email}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group>
              <Form.Row>
                <Form.Group as={Col}>
                  <Form.Label>Password</Form.Label>
                  <Form.Control name={"password"} value={this.state.password} onBlur={this.handleChange}
                                onChange={this.handleChange}
                                type={"password"} isInvalid={!!this.state.error_password}/>
                  <Form.Control.Feedback type={"invalid"}>{this.state.error_password}</Form.Control.Feedback>
                </Form.Group>
                <Form.Group as={Col}>
                  <Form.Label>Re-enter Password</Form.Label>
                  <Form.Control name={"password2"} value={this.state.password2} onBlur={this.handleChange}
                                onChange={this.handleChange}
                                type={"password"} isInvalid={!!this.state.error_password2}/>
                  <Form.Control.Feedback type={"invalid"}>{this.state.error_password2}</Form.Control.Feedback>
                </Form.Group>
              </Form.Row>
              <ProgressBar variant={["danger", "warning", "info", "success"][this.state.strength - 1]}
                           label={["very weak", "weak", "medium", "strong"][this.state.strength - 1]}
                           now={25 * this.state.strength} className={"password-strength"}/>
            </Form.Group>
            <Button variant="primary" type="submit" block disabled={this.state.loading}>
              Sign Up{" "}
              {this.state.loading &&
              <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                       size="sm" role="status" aria-hidden="true"/>}
            </Button>
          </Form>
          <br/>
          <p style={{textAlign: "center"}}>Already have an account? <Link className={"inline"} to={'/login'}>
            Log in</Link> now!</p>
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
    maxWidth: 480
  },
};

export default SignUp
