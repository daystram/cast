import React, {Component} from 'react';
import {Alert, Button, Card, Col, Container, Form, Image, Row, Spinner} from "react-bootstrap";
import Dropzone from "react-dropzone";
import SidebarProfile from "./SidebarProfile";
import axios from "axios";
import urls from "../helper/url";
import abbreviate from "../helper/abbreviate";
import auth from "../helper/auth";

let timeout = {};

class Profile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: auth().username(),
      subscribers: 0,
      views: 0,
      video_count: 0,
      name: "",
      email: "",
      before: {},
      error_name: "",
      error_email: "",
      loading_edit: false,
      loading_info: true,
      editing: false
    };
    this.fetchUser = this.fetchUser.bind(this);
    this.pressEdit = this.pressEdit.bind(this);
    this.pressDelete = this.pressDelete.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.checkAvailability = this.checkAvailability.bind(this);
    this.submitForm = this.submitForm.bind(this);
    this.fetchUser();
    document.title = "Profile | cast";
  }

  fetchUser() {
    axios.get(urls().user_info(), {}).then((response) => {
      this.setState({loading_info: false});
      if (response.data.code === 200) {
        let data = response.data.data;
        this.setState({
          name: data.name,
          email: data.email,
          subscribers: data.subscribers,
          views: data.views,
          video_count: data.uploads
        })
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading_info: false});
    });
  }

  handleChange(e) {
    this.setState({error_edit: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  pressEdit() {
    if (this.state.editing) {
      this.submitForm()
    } else {
      this.setState({
        before: {
          name: this.state.name,
          email: this.state.email,
        },
        new_profile: "",
        error_name: "",
        error_email: "",
        editing: true
      })
    }
  }

  pressDelete() {
    if (this.state.editing) {
      this.setState({
        name: this.state.before.name,
        email: this.state.before.email,
        new_profile: "",
        before: {},
        editing: false
      })
    }
  }

  validate(field, value) {
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
        if (value.trim().length > 22) {
          this.setState({error_name: "Name too long"});
          return false;
        }
        this.setState({error_name: ""});
        if (value.trim() !== this.state.before.name.trim()) this.checkAvailability(field, value.trim());
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
        if (value.trim() !== this.state.before.email.trim()) this.checkAvailability(field, value.trim());
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
        if (response.data.code !== 200) {
          this.setState({[`error_${field}`]: response.data.error});
        } else {
          this.setState({[`error_${field}`]: ""});
        }
      }).catch((error) => {
        console.log(error);
        this.setState({error_edit: "An error has occurred!"});
      });
    }, 400)
  }

  submitForm() {
    let ok = true;
    ok &= !this.state.error_name;
    ok &= !this.state.error_email;
    if (!ok) return;
    this.setState({error_edit: "", loading_edit: true});
    const form = new FormData();
    form.append("name", this.state.name);
    form.append("email", this.state.email);
    if (this.state.new_profile) form.append("profile", this.state.new_profile);
    axios.put(urls().edit_user(), form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data"
        },
      }
    ).then((response) => {
      if (response.data.code === 200) {
        this.setState({
          editing: false,
          loading_edit: false,
          before: {},
          error_name: "",
          error_email: "",
          new_profile: "",
        });
        if (this.state.new_profile) window.location.reload();
        console.log(response.data);
        this.setState({error_edit: response.data.error, loading_edit: false});
      }
    }).catch((error) => {
      console.log(error);
      this.setState({error_edit: "An error has occurred! Please try again.", loading_edit: false});
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
              <Row>
                <Col md sm={12} className={"responsive-fold"}>
                  <div style={style.profile_bar}>
                    {this.state.editing ?
                      <Dropzone accept={"image/*"} multiple={false}
                                onDrop={files => this.setState({new_profile: files[0]})}
                                disabled={this.state.loading_edit}>
                        {({getRootProps, getInputProps}) => (
                          <>
                            <Image
                              src={this.state.new_profile ? URL.createObjectURL(this.state.new_profile) : urls().profile(this.state.username)}
                              height={128} width={128} roundedCircle
                              style={{position: "absolute", zIndex: 0, objectFit: "cover"}}/>
                            <section
                              style={{...style.profile_upload, ...(this.state.new_profile ? style.profile_upload_modified : {})}}>
                              <div {...getRootProps()}>
                                <input {...getInputProps()} />
                                <p style={{textAlign: "center", lineHeight: "128px"}}>
                                  {!this.state.new_profile &&
                                  <span className="material-icons"
                                        style={{fontSize: 32, color: "dimgray"}}>publish</span>
                                  }
                                </p>
                              </div>
                            </section>
                          </>
                        )}
                      </Dropzone> :
                      <Image src={urls().profile(this.state.username)}
                             height={128} width={128} roundedCircle/>}
                    <div style={style.profile_name}>
                      {this.state.loading_info && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
                      {this.state.error_edit && <Alert variant={"danger"}>{this.state.error_edit}</Alert>}
                      {this.state.editing ?
                        <Form.Group>
                          <Form.Control name={"name"} value={this.state.name} onBlur={this.handleChange}
                                        onChange={this.handleChange} type={"text"} size={"lg"} style={style.name}
                                        isInvalid={this.state.error_name} placeholder={"Name"}/>
                          <Form.Control.Feedback type={"invalid"}>{this.state.error_name}</Form.Control.Feedback>
                        </Form.Group> :
                        <h1 style={style.name}>{this.state.name}</h1>
                      }
                      <h3 style={style.sub_count}>{abbreviate().number(this.state.subscribers)} subscribers</h3>
                    </div>
                  </div>

                  <h1 style={style.h1}>Showcase</h1>
                  <Card body style={style.profile_detail}>
                    <Row>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.subscribers}</p>
                        <p style={style.show_caption}>subscriber{this.state.subscribers === 1 ? "" : "s"}</p>
                      </Col>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.views}</p>
                        <p style={style.show_caption}>total view{this.state.views === 1 ? "" : "s"}</p>
                      </Col>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.video_count}</p>
                        <p style={style.show_caption}>cast{this.state.video_count === 1 ? "" : "s"} uploaded</p>
                      </Col>
                    </Row>
                  </Card>

                  <h1 style={style.h1}>Details</h1>
                  <Card body style={style.profile_detail}>
                    <Row>
                      <Col md={6} sm={12}>
                        <Form.Group>
                          <Form.Label>Email</Form.Label>
                          <Form.Control name={"email"} value={this.state.email} onBlur={this.handleChange}
                                        onChange={this.handleChange} type={"text"} style={style.email}
                                        isInvalid={this.state.editing ? this.state.error_email : false}
                                        disabled={!this.state.editing}/>
                          <Form.Control.Feedback type={"invalid"}>{this.state.error_email}</Form.Control.Feedback>
                        </Form.Group>
                      </Col>
                      <Col md={6} sm={12}>
                        <Form.Group>
                          <Form.Label>Username</Form.Label>
                          <Form.Control name={"Username"} value={this.state.username} type={"text"}
                                        style={style.email} disabled={true}/>
                        </Form.Group>
                      </Col>
                    </Row>
                  </Card>
                </Col>
                <Col md={"auto"} sm={12}>
                  <Button variant={"success"} block size={"sm"} style={style.button} onClick={this.pressEdit}
                          disabled={this.state.before && (this.state.name === this.state.before.name &&
                            this.state.email === this.state.before.email && !this.state.new_profile)}>
                    {this.state.loading_edit ?
                      <Spinner animation="grow" style={style.spinner}/> :
                      <span className="material-icons">{this.state.editing ? "check" : "edit"}</span>
                    }
                  </Button>
                  <Button variant={this.state.editing ? "info" : "danger"} block size={"sm"}
                          onClick={this.pressDelete} style={{...style.button, marginBottom: 32}}>
                    <span className="material-icons">{this.state.editing ? "clear" : "delete"}</span></Button>
                </Col>
              </Row>
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
  name: {
    margin: 0,
    fontSize: "3.5rem"
  },
  profile_upload: {
    width: 128,
    height: 128,
    background: "#f0f0f088",
    border: "3px dashed #ddddddaa",
    flexShrink: 0,
    borderRadius: 64,
    zIndex: 100,
  },
  profile_upload_modified: {
    background: "#f0f0f022",
    border: ""
  },
  profile_bar: {
    display: "flex",
    marginBottom: 48,
  },
  profile_name: {
    marginLeft: 16,
    alignSelf: "center"
  },
  sub_count: {
    fontSize: 22
  },
  profile_detail: {
    borderRadius: "8px 48px 8px 8px",
    marginBottom: 32
  },
  email: {
    fontSize: 18
  },
  show_count: {
    textAlign: "center",
    fontSize: 28,
    fontWeight: 800,
    marginBottom: 8
  },
  show_caption: {
    textAlign: "center",
    marginBottom: 0,
    fontWeight: 600
  }
};

export default Profile
