import React, {Component} from 'react';
import {Alert, Button, Container, Form, Spinner} from "react-bootstrap";
import axios from "axios";
import urls from "../helper/url";
import {Redirect} from "react-router-dom";

class Verify extends Component {
  constructor(props) {
    super(props);
    this.state = {
      key: this.props.match.params.key,
      email: "",
      error_verify: "",
      loading: this.props.match.params.key
    };
    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  componentDidMount() {
    if (this.state.key) {
      axios.post(urls().verify(), {
        key: this.state.key
      }).then((response) => {
        this.setState({loading: false});
        if (response.data.code === 200) {
          this.setState({success: true});
        } else {
          this.setState({error_verify: "Verification link invalid!"});
        }
      }).catch((error) => {
        console.log(error);
        this.setState({error_verify: "An error has occurred!"});
        this.setState({loading: false});
      });
    }
  }

  handleChange(e) {
    this.setState({error_verify: ""});
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
      default:
        return false;
    }
  }

  submitForm(e) {
    e.preventDefault();
    let ok = true;
    ok &= !this.state.error_email;
    if (!ok) return;
    this.setState({loading: true, verified: false, success: false});
    axios.post(urls().resend_verify(), {
      email: this.state.email.trim(),
    }).then((response) => {
      this.setState({loading: false});
      switch (response.data.code) {
        case 200:
          this.setState({success: true});
          break;
        case 404:
          this.setState({error_verify: "Email not registered!"});
          break;
        case 409:
          this.setState({verified: true});
          break;
        default:
          this.setState({error_verify: "An error has occurred!"});
          break;
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_verify: "An error has occurred!"});
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <Container fluid style={style.content_container}>
        {this.state.redirectLogin ? <Redirect to={"/login"}/> : ""}
        {this.state.redirectVerify ? <Redirect to={"/verify"}/> : ""}
        <h1 style={style.h1}>Verify Email</h1>
        {this.state.key ?
          <>
            {this.state.error_verify &&
            <Alert variant={"danger"}>
              <Alert.Heading>Invalid Link!</Alert.Heading>
              <p>Your verification link is invalid! Please request for a new verification link sent to your email.</p>
              <hr/>
              <div className="d-flex justify-content-end">
                <Button onClick={() => this.setState({redirectVerify: true, key: "", error_verify: ""})}
                        variant="outline-danger">
                  Re-send Link
                </Button>
              </div>
            </Alert>}
            {this.state.success &&
            <Alert variant="success">
              <Alert.Heading>Welcome!</Alert.Heading>
              <p>
                Your email has been successfully verified. You can now log in, start sharing videos and stream live!
              </p>
              <hr/>
              <div className="d-flex justify-content-end">
                <Button onClick={() => this.setState({redirectLogin: true})} variant="outline-success">
                  Log In
                </Button>
              </div>
            </Alert>}
            {this.state.loading && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
          </> :
          <>
            <p>Enter your email below to re-send verification link.</p>
            {this.state.error_verify && <Alert variant={"danger"}>{this.state.error_verify}</Alert>}
            {this.state.success && <Alert variant={"success"}>Verification link sent!</Alert>}
            {this.state.verified &&
            <Alert variant="success">
              <Alert.Heading>Welcome!</Alert.Heading>
              <p>
                Your email has already been verified. You can now log in, start sharing videos and stream live!
              </p>
              <hr/>
              <div className="d-flex justify-content-end">
                <Button onClick={() => this.setState({redirectLogin: true})} variant="outline-success">
                  Log In
                </Button>
              </div>
            </Alert>}
            <Form.Group>
              <Form.Label>Email</Form.Label>
              <Form.Control name={"email"} value={this.state.email} onBlur={this.handleChange}
                            onChange={this.handleChange} type={"email"}
                            isInvalid={!!this.state.error_email}/>
              <Form.Control.Feedback type={"invalid"}>{this.state.error_email}</Form.Control.Feedback>
            </Form.Group>
            <Button variant="primary" type="submit" block disabled={this.state.loading}
                    onClick={this.submitForm}>
              Send Email{" "}
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
    padding: "36px 0 0 0",
    maxWidth: 480
  },
  spinner: {
    margin: "32px auto 64px auto",
    display: "block"
  },
};

export default Verify