import React, {Component} from 'react';
import {Alert, Button, Card, Col, Container, Form, ProgressBar, Row, Spinner} from "react-bootstrap";
import axios from "axios";
import bsCustomFileInput from 'bs-custom-file-input'
import {WithContext as ReactTags} from 'react-tag-input';
import SidebarProfile from "./SidebarProfile";
import urls from "../helper/url";
import CastEditable from "./CastEditable";
import {Prompt} from "react-router-dom";

import './tags.css'
import './file.css'
import auth from "../helper/auth";
import MediaQuery from "react-responsive";
import {MOBILE_BP} from "../constants/breakpoint";

let timeout = null;

class Manage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: "",
      description: "",
      tags: [],
      thumbnail: "",
      video: "",
      error_title: "",
      error_description: "",
      error_tags: "",
      error_thumbnail: "",
      error_video: "",
      error_upload: "",
      progress: 0,
      list: [],
      uploading: false,
      loading: true,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleTagAdd = this.handleTagAdd.bind(this);
    this.handleTagDelete = this.handleTagDelete.bind(this);
    this.handleTagDrag = this.handleTagDrag.bind(this);
    this.handleChangeFile = this.handleChangeFile.bind(this);
    this.submitForm = this.submitForm.bind(this);
    document.title = "Manage | cast";
  }

  componentDidMount() {
    bsCustomFileInput.init();
    this.fetchVideos()
  }

  fetchVideos() {
    axios.get(urls().list(), {
      params: {
        author: auth().username(),
        count: 8,
        offset: 0,
      }
    }).then((response) => {
      this.setState({loading: false});
      if (response.data.code === 200) {
        this.setState({list: response.data.data})
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: false});
    });
  }

  handleChange(e) {
    this.setState({error_upload: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  handleTagAdd(tag) {
    if (this.validate("tag", tag.text)) {
      this.setState(state => ({tags: [...state.tags, tag]}));
    }
  }

  handleTagDelete(i) {
    const {tags} = this.state;
    this.setState({
      tags: tags.filter((tag, index) => index !== i),
    });
  }

  handleTagDrag(tag, currPos, newPos) {
    const tags = [...this.state.tags];
    const newTags = tags.slice();
    newTags.splice(currPos, 1);
    newTags.splice(newPos, 0, tag);
    this.setState({tags: newTags});
  }

  handleChangeFile(e) {
    this.setState({error_upload: ""});
    this.setState({[e.target.name]: e.target.files[0]});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    switch (field) {
      case "title":
        if (!value.trim()) {
          this.setState({error_title: "Please enter video title"});
          return false;
        }
        this.setState({error_title: ""});
        this.checkAvailability(value);
        return true;
      case "description":
        if (!value.trim()) {
          this.setState({error_description: "Please enter video description"});
          return false;
        }
        this.setState({error_description: ""});
        return true;
      case "tag":
        let tagRe = /^[A-Za-z0-9]+$/;
        if (!tagRe.test(value.trim())) {
          this.setState({error_tags: "Please insert alphanumeric tags"});
          return false;
        }
        this.setState({error_tags: ""});
        return true;
      case "tags":
        if (value.length === 0) {
          this.setState({error_tags: "Please enter video tags"});
          return false;
        }
        this.setState({error_tags: ""});
        return true;
      case "thumbnail":
        if (!value) {
          this.setState({error_thumbnail: "Please select thumbnail image"});
          return false;
        }
        this.setState({error_thumbnail: ""});
        return true;
      case "video":
        if (!value) {
          this.setState({error_video: "Please select video file"});
          return false;
        }
        this.setState({error_video: ""});
        return true;
      default:
        return false;
    }
  }

  checkAvailability(value) {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      axios.get(urls().title_check(), {
        params: {
          title: value.trim()
        }
      }).then((response) => {
        if (response.data.code !== 200) {
          this.setState({error_title: response.data.error});
        } else {
          this.setState({error_title: ""});
        }
      }).catch((error) => {
        console.log(error);
        this.setState({error_upload: "An error has occurred!"});
      });
    }, 400)
  }

  submitForm(e) {
    e.preventDefault();
    let ok = true;
    if (!this.state.attempted) {
      this.setState({attempted: true});
      ok &= this.validate("title", this.state.title);
      ok &= this.validate("description", this.state.description);
      ok &= this.validate("tags", this.state.tags);
      ok &= this.validate("thumbnail", this.state.thumbnail);
      ok &= this.validate("video", this.state.video);
    } else {
      ok &= !this.state.error_title;
      ok &= !this.state.error_description;
      ok &= !this.state.error_tags;
      ok &= !this.state.error_thumbnail;
      ok &= !this.state.error_video;
    }
    if (!ok) return;

    this.setState({uploading: true});
    const form = new FormData();
    form.append("title", this.state.title);
    form.append("description", this.state.description);
    form.append("tags", this.state.tags.map(tag => tag.text).join(","));
    form.append("thumbnail", this.state.thumbnail);
    form.append("video", this.state.video);
    axios.post(urls().upload(), form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data"
        },
        onUploadProgress: (progress) => {
          this.setState({progress: (progress.loaded * 100) / progress.total})
        }
      }
    ).then((response) => {
      console.log(response);
      this.setState({
        title: "",
        description: "",
        tags: [],
        thumbnail: "",
        video: "",
        progress: 0,
        uploading: false,
        success: true
      });
      this.setState({loading: true, list: []});
      this.fetchVideos();
    }).catch((error) => {
      console.log(error);
      this.setState({error_upload: "An error has occurred! Please try again.", uploading: false});
    });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{marginBottom: 32}}>
                <Card body style={{borderRadius: "8px 48px 8px 8px"}}><SidebarProfile/></Card>
              </Col>
            </MediaQuery>
            <Col xl={10} sm={12} className={"mid-container-right"}>
              <h1 style={style.h1}>Upload</h1>
              {this.state.error_upload && <Alert variant={"danger"}>{this.state.error_upload}</Alert>}
              {this.state.success &&
              <Alert variant="success" onClose={() => this.setState({success: false})} dismissible>
                <Alert.Heading>Video Successfully Uploaded!</Alert.Heading>
                <p>
                  Your video is now currently being processed and will be available in a few minutes. Video may first
                  appear in lower quality, as higher resolutions are still being ingested.
                </p>
              </Alert>}
              <Form noValidate autocomplete={"off"} onSubmit={this.submitForm}>
                <Form.Row>
                  <Col md={6} sm={12}>
                    <Form.Group>
                      <Form.Label>Title</Form.Label>
                      <Form.Control name={"title"} value={this.state.title} onBlur={this.handleChange}
                                    onChange={this.handleChange} type={"text"} size={"lg"} style={{fontSize: "2rem"}}
                                    isInvalid={!!this.state.error_title} disabled={this.state.uploading}/>
                      <Form.Control.Feedback type={"invalid"}>{this.state.error_title}</Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Description</Form.Label>
                      <Form.Control name={"description"} value={this.state.description} onBlur={this.handleChange}
                                    onChange={this.handleChange} as={"textarea"} rows={7}
                                    isInvalid={!!this.state.error_description} disabled={this.state.uploading}/>
                      <Form.Control.Feedback type={"invalid"}>{this.state.error_description}</Form.Control.Feedback>
                    </Form.Group>
                  </Col>
                  <Col md={6} sm={12}>
                    <Form.Group>
                      <Form.Label>Tags</Form.Label>
                      <ReactTags
                        classNames={{
                          tags: this.state.error_tags ? "ReactTags__tags__error" : (this.state.uploading ? "ReactTags__tags__disabled" : "ReactTags__tags"),
                          tagInput: this.state.tags.length === 5 ? "ReactTags__tagInput__disabled" : "ReactTags__tagInput"
                        }}
                        tags={this.state.tags} autofocus={false} delimiters={[13, 32, 188]} maxLength={12}
                        placeholder={""} readOnly={this.state.uploading} handleAddition={this.handleTagAdd}
                        handleDelete={this.handleTagDelete} handleDrag={this.handleTagDrag}
                        handleInputChange={() => this.setState({error_tags: ""})}
                        handleInputBlur={() => this.validate("tags", this.state.tags)}
                        handleTagClick={this.handleTagClick}/>
                      {this.state.error_tags &&
                      <div style={style.invalidText}>{this.state.error_tags}</div>}
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Thumbnail</Form.Label>
                      <div className={"custom-file"}>
                        <Form.Control name={"thumbnail"} type={"file"} className={"custom-file-input"} accept="image/*"
                                      onChange={this.handleChangeFile} onBlur={this.handleChangeFile}
                                      isInvalid={!!this.state.error_thumbnail} disabled={this.state.uploading}/>
                        <Form.Control.Feedback type={"invalid"}>{this.state.error_thumbnail}</Form.Control.Feedback>
                        <Form.Label style={{overflow: "hidden"}} className="custom-file-label">
                          {this.state.thumbnail && this.state.thumbnail.name}</Form.Label>
                      </div>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Video</Form.Label>
                      <div className={"custom-file"}>
                        <Form.Control name={"video"} type={"file"} className={"custom-file-input"} accept="video/*"
                                      onChange={this.handleChangeFile} onBlur={this.handleChangeFile}
                                      isInvalid={!!this.state.error_video} disabled={this.state.uploading}/>
                        <Form.Control.Feedback type={"invalid"}>{this.state.error_video}</Form.Control.Feedback>
                        <Form.Label style={{overflow: "hidden"}} className="custom-file-label">
                          {this.state.video && this.state.video.name}</Form.Label>
                      </div>
                    </Form.Group>
                    <ProgressBar animated now={this.state.progress} style={style.progress}/>
                    <Button variant="primary" type="submit" block disabled={this.state.uploading}>
                      Upload{" "}
                      {this.state.uploading &&
                      <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                               size="sm" role="status" aria-hidden="true"/>}
                    </Button>
                  </Col>
                </Form.Row>
              </Form>
              <hr/>
              <h1 style={style.h1}>Manage</h1>
              <Col>
                {!this.state.loading && (this.state.list ? this.state.list.map(video =>
                  <Row key={video.hash}>
                    <CastEditable video={video} deletable={true} onDelete={() => this.fetchVideos()}/>
                  </Row>
                ) : <h5 style={style.h5}>No videos uploaded yet!</h5>)}
                {this.state.loading && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
              </Col>
            </Col>
          </Row>
        </Container>
        <Prompt
          when={this.state.uploading}
          message={() => `Are you sure? Your current upload will be cancelled.`}/>
      </>
    )
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
  },
  content_container: {},
  progress: {
    marginBottom: 16
  },
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16
  },
  spinner: {
    margin: "32px auto 64px auto",
    display: "block"
  },
  invalidInput: {
    borderColor: "#dc3545",
    paddingRight: "calc(1.5em + .75rem)",
    backgroundImage: "url(data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='12' hei…circle cx='6' cy='8.2' r='.6' fill='%23dc3545' stroke='none'/%3e%3c/svg%3e)",
    backgroundRepeat: "no-repeat",
    backgroundPosition: "right calc(.375em + .1875rem) center",
    backgroundSize: "calc(.75em + .375rem) calc(.75em + .375rem)",
  },
  invalidText: {
    display: "block",
    width: "100%",
    marginTop: ".25rem",
    fontSize: "80%",
    color: "#dc3545",
  }
};

export default Manage
