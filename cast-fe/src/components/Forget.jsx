import React, {Component} from 'react';
import {Alert, Button, Col, Container, Form, ProgressBar, Spinner} from "react-bootstrap";
import axios from "axios";
import urls from "../helper/url";
import {withRouter} from "react-router-dom";
import zxcvbn from "zxcvbn";

// TODO: WIP!
class Forget extends Component {
  constructor(props) {
    super(props);
    let params = new URLSearchParams(this.props.location.search);
    this.state = {
      key: params.get("key"),
      email: "",
      error_reset: "",
      password: "",
      password2: "",
      error_update: "",
      loading: params.get("key")
    };
    this.handleChange = this.handleChange.bind(this);
    this.sendLink = this.sendLink.bind(this);
    this.updatePassword = this.updatePassword.bind(this);
  }

  componentDidMount() {
    if (this.state.key) {
      axios.post(urls().check_reset(), {
        key: this.state.key
      }).then((response) => {
        this.setState({loading: false});
        if (response.data.code === 200) {
          this.setState({valid: true});
        } else {
          this.setState({error_reset: "Password reset link invalid!"});
        }
      }).catch((error) => {
        console.log(error);
        this.setState({error_reset: "An error has occurred!"});
        this.setState({loading: false});
      });
    }
  }

  handleChange(e) {
    this.setState({error_reset: "", error_update: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    switch (field) {
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
        return true;
      case "password":
        if (!value) {
          this.setState({error_password: "Please enter your password"});
          return false;
        }
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

  sendLink(e) {
    e.preventDefault();
    let ok = true;
    ok &= !this.state.error_email;
    if (!ok) return;
    this.setState({loading: true, success: false});
    axios.post(urls().reset_password(), {
      email: this.state.email.trim(),
    }).then((response) => {
      this.setState({loading: false});
      switch (response.data.code) {
        case 200:
          this.setState({success: true});
          break;
        case 403:
          this.setState({unverified: false});
          break;
        case 404:
          this.setState({error_reset: "Email not registered!"});
          break;
        default:
          this.setState({error_reset: "An error has occurred!"});
          break;
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_reset: "An error has occurred!", loading: false});
    });
  }

  updatePassword(e) {
    e.preventDefault();
    let ok = true;
    this.setState({success: false});
    if (!this.state.attempted) {
      this.setState({attempted: true});
      ok &= this.validate("password", this.state.password);
      ok &= this.validate("password2", this.state.password2);
    } else {
      ok &= !this.state.error_password;
      ok &= !this.state.error_password2;
    }
    if (!ok) return;
    this.setState({loading: true});
    axios.put(urls().update_password(), {
      key: this.state.key,
      password: this.state.password.trim(),
    }).then((response) => {
      this.setState({loading: false, password: "", password2: ""});
      if (response.data.code === 200) {
        this.setState({
          success: true,
          password: "",
          password2: "",
          error_password: "",
          error_password2: "",
        });
      } else {
        this.setState({error_update: "An error has occurred!"});
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_update: "An error has occurred!", loading: false});
    });
  }

  render() {
    let strength = zxcvbn(this.state.password).score;
    return (
      <Container fluid style={style.content_container}>
        <h1 style={style.h1}>Forget Password</h1>
        {this.state.key ?
          <>
            {this.state.error_reset &&
            <Alert variant={"danger"}>
              <Alert.Heading>Invalid Link!</Alert.Heading>
              <p>Your password reset link is invalid! Please request for a new password reset link.</p>
              <hr/>
              <div className="d-flex justify-content-end">
                <Button onClick={() => {
                  this.setState({key: "", error_reset: ""});
                  this.props.history.push("/forget")
                }} variant="outline-danger">
                  Re-send Link
                </Button>
              </div>
            </Alert>}
            {this.state.valid && <>
              {this.state.success ?
                <Alert variant="success">
                  <Alert.Heading>Welcome!</Alert.Heading>
                  <p>
                    Your password has been successfully reset. You can now log in with your new password!
                  </p>
                  <hr/>
                  <div className="d-flex justify-content-end">
                    <Button onClick={() => this.props.history.push("/login")} variant="outline-success">
                      Log In
                    </Button>
                  </div>
                </Alert> :
                <>
                  <p>Enter your new password below.</p>
                  {this.state.error_update && <Alert variant={"danger"}>{this.state.error_update}</Alert>}
                  <Form noValidate onSubmit={this.updatePassword}>
                    <Form.Group>
                      <Form.Row>
                        <Form.Group as={Col}>
                          <Form.Control name={"password"} value={this.state.password} onBlur={this.handleChange}
                                        onChange={this.handleChange} placeholder={"Password"}
                                        type={"password"} isInvalid={!!this.state.error_password}/>
                          <Form.Control.Feedback type={"invalid"}>{this.state.error_password}</Form.Control.Feedback>
                        </Form.Group>
                        <Form.Group as={Col}>
                          <Form.Control name={"password2"} value={this.state.password2} onBlur={this.handleChange}
                                        onChange={this.handleChange} placeholder={"Re-enter Password"}
                                        type={"password"} isInvalid={!!this.state.error_password2}/>
                          <Form.Control.Feedback type={"invalid"}>{this.state.error_password2}</Form.Control.Feedback>
                        </Form.Group>
                      </Form.Row>
                      <ProgressBar variant={["danger", "warning", "info", "success"][strength - 1]}
                                   label={["very weak", "weak", "medium", "strong"][strength - 1]}
                                   now={25 * strength} className={"password-strength"}/>
                    </Form.Group>
                    <Button variant="primary" type="submit" block disabled={this.state.loading}>
                      Reset Password{" "}
                      {this.state.loading &&
                      <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                               size="sm" role="status" aria-hidden="true"/>}
                    </Button>
                  </Form>
                </>}
            </>}
            {!this.state.valid && this.state.loading &&
            <Spinner style={style.spinner} animation="grow" variant="primary"/>}
          </> :
          <>
            <p>Enter your email below to get your password reset link.</p>
            {this.state.error_reset && <Alert variant={"danger"}>{this.state.error_reset}</Alert>}
            {this.state.success && <Alert variant={"success"}>Password reset link sent!</Alert>}
            {this.state.unverified &&
            <Alert variant={"warning"}>
              <Alert.Heading>Email Unverified!</Alert.Heading>
              <p>Your account email has not been verified. Please check your email for a verification link before
                resetting your password.</p>
              <hr/>
              <div className="d-flex justify-content-end">
                <Button onClick={() => {
                  this.setState({key: "", error_reset: ""});
                  this.props.history.push("/verify")
                }} variant="outline-warning">
                  Re-send Link
                </Button>
              </div>
            </Alert>}
            <Form.Group>
              <Form.Control name={"email"} value={this.state.email} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"email"} placeholder={"Email"}
                            isInvalid={!!this.state.error_email}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_email}</Form.Control.Feedback>
            </Form.Group>
            <Button variant="primary" type="submit" block disabled={this.state.loading}
                    onClick={this.sendLink}>
              Send Reset Link{" "}
              {this.state.loading &&
              <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                       size="sm" role="status" aria-hidden="true"/>}
            </Button>
          </>
        }
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
  spinner: {
    margin: "32px auto 64px auto",
    display: "block"
  },
};

export default withRouter(Forget)
