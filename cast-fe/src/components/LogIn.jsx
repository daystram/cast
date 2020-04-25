import React, {Component} from 'react';
import {Alert, Button, Container, Form, Spinner} from "react-bootstrap";
import axios from "axios";
import urls from "../helper/url";
import {Link, withRouter} from "react-router-dom";

class LogIn extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: "",
      remember: false,
      error_username: "",
      error_password: "",
      error_login: "",
      loading: false,
      verified: true
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
    this.setState({loading: true, verified: true});
    axios.post(urls().login(), {
      username: this.state.username.trim(),
      password: this.state.password.trim(),
      remember: this.state.remember
    }).then((response) => {
      this.setState({loading: false});
      switch (response.data.code) {
        case 200:
          const {from} = this.props.location.state || {from: false};
          if (from) this.props.history.push(from.pathname);
          else this.props.history.goBack();
          return;
        case 404:
          this.setState({error_username: "Username not registered"});
          return;
        case 403:
          this.setState({error_password: "Incorrect password"});
          return;
        case 406:
          this.setState({verified: false});
          return;
        default:
          this.setState({error_login: "An error has occurred!"});
          return;
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_login: "An error has occurred!"});
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <Container fluid style={style.content_container}>
        <h1 style={style.h1}>Log In</h1>
        {this.state.error_login && <Alert variant={"danger"}>{this.state.error_login}</Alert>}
        {!this.state.verified && <Alert variant={"warning"}>
          <Alert.Heading>Verify Your Email</Alert.Heading>
          <p>You need to verify your email before logging in. Check your email for an activation link.</p>
          <hr/>
          <div className="d-flex justify-content-end">
            <Button onClick={() => this.props.history.push("/verify")} variant="outline-warning">
              Resend Link
            </Button>
          </div>
        </Alert>}
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
          <Form.Group>
            <Form.Check label="Remember me" checked={this.state.remember}
                        onChange={(e) => this.setState({remember: e.target.checked})}/>
          </Form.Group>
          <Button variant="primary" type="submit" block disabled={this.state.loading}>
            Log In{" "}
            {this.state.loading &&
            <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                     size="sm" role="status" aria-hidden="true"/>}
          </Button>
        </Form>
        <br/>
        <p style={{textAlign: "center"}}><Link to={'/forgot'}>Forgot password?</Link></p>
        <p style={{textAlign: "center"}}>Don't have an account? <Link to={'/signup'}>Sign up</Link> now!</p>
      </Container>
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

export default withRouter(LogIn)
