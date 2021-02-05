import React, { Component } from "react";
import {
  Alert,
  Badge,
  Button,
  Card,
  Col,
  Form,
  Image,
  Modal,
  Row,
  Spinner,
} from "react-bootstrap";
import Dropzone from "react-dropzone";
import { currentHash } from "../helper/url";
import format from "../helper/format";
import { Prompt, withRouter } from "react-router-dom";
import { WithContext as ReactTags } from "react-tag-input";
import "../styles/tags.css";
import { THUMBNAIL_MAX_SIZE } from "../constants/file";
import {
  VIDEO_DESC_CHAR_LIMIT,
  VIDEO_TAG_CHAR_LIMIT,
  VIDEO_TAG_COUNT,
  VIDEO_TITLE_CHAR_LIMIT,
} from "../constants/video";
import api from "../apis/api";

const resolutions = ["Processing", "240p", "360p", "480p", "720p", "1080p"];
let timeout = {};

class CastEditable extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: this.props.video.title,
      tags: this.props.video.tags
        ? this.props.video.tags.map((tag) => {
            return { id: tag, text: tag };
          })
        : [],
      description: this.props.video.description,
      thumbnail: api.cdn.thumbnail(this.props.video.hash),
      error_title: "",
      error_tags: "",
      error_description: "",
      error_thumbnail: "",
      error_edit: "",
      error_delete: "",
      before: {},
      editing: false,
      loading_edit: false,
      loading_delete: false,
      prompt: false,
      updated: false,
    };
    this.pressEdit = this.pressEdit.bind(this);
    this.pressDelete = this.pressDelete.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.handleTagAdd = this.handleTagAdd.bind(this);
    this.handleTagDelete = this.handleTagDelete.bind(this);
    this.handleTagDrag = this.handleTagDrag.bind(this);
    this.submitForm = this.submitForm.bind(this);
    this.deleteVideo = this.deleteVideo.bind(this);
    this.openVideo = this.openVideo.bind(this);
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (
      !this.state.editing &&
      !this.state.updated &&
      (prevState.title !== this.props.video.title ||
        prevState.description !== this.props.video.description ||
        prevProps.video.tags !== this.props.video.tags)
    ) {
      this.setState({
        title: this.props.video.title,
        description: this.props.video.description,
        tags: this.props.video.tags
          ? this.props.video.tags.map((tag) => {
              return { id: tag, text: tag };
            })
          : [],
      });
    }
  }

  pressEdit() {
    if (this.state.editing) {
      this.submitForm();
    } else {
      this.setState({
        before: {
          title: this.state.title,
          tags: this.state.tags,
          description: this.state.description,
          thumbnail: this.state.thumbnail,
        },
        new_thumbnail: "",
        error_title: "",
        error_tags: "",
        error_description: "",
        error_thumbnail: "",
        editing: true,
      });
    }
  }

  pressDelete() {
    if (this.state.editing) {
      this.setState({
        title: this.state.before.title,
        tags: this.state.before.tags,
        description: this.state.before.description,
        thumbnail: this.state.before.thumbnail,
        new_thumbnail: "",
        before: {},
        error_title: "",
        error_tags: "",
        error_description: "",
        error_thumbnail: "",
        error_edit: "",
        editing: false,
      });
    } else {
      this.setState({ prompt: true });
    }
  }

  handleChange(e) {
    this.setState({ error_edit: "" });
    this.setState({ [e.target.name]: e.target.value });
    this.validate(e.target.name, e.target.value);
  }

  handleTagAdd(tag) {
    if (this.validate("tag", tag.text)) {
      this.setState((state) => ({ tags: [...state.tags, tag] }));
    }
  }

  handleTagDelete(i) {
    const { tags } = this.state;
    this.setState({
      tags: tags.filter((tag, index) => index !== i),
    });
  }

  handleTagDrag(tag, currPos, newPos) {
    const tags = [...this.state.tags];
    const newTags = tags.slice();
    newTags.splice(currPos, 1);
    newTags.splice(newPos, 0, tag);
    this.setState({ tags: newTags });
  }

  validate(field, value) {
    switch (field) {
      case "title":
        if (!value.trim()) {
          this.setState({ error_title: "Please enter video title" });
          return false;
        }
        if (value.trim().length > VIDEO_TITLE_CHAR_LIMIT) {
          this.setState({ error_title: "Title too long" });
          return false;
        }
        this.setState({ error_title: "" });
        if (
          value.toLowerCase().trim() !==
          this.state.before.title.toLowerCase().trim()
        )
          this.checkAvailability(value);
        return true;
      case "tag":
        let tagRe = /^[A-Za-z0-9]+$/;
        if (!tagRe.test(value.trim())) {
          this.setState({ error_tags: "Please insert alphanumeric tags" });
          return false;
        }
        this.setState({ error_tags: "" });
        return true;
      case "tags":
        if (value.length === 0) {
          this.setState({ error_tags: "Please enter video tags" });
          return false;
        }
        this.setState({ error_tags: "" });
        return true;
      case "description":
        if (!value.trim()) {
          this.setState({
            error_description: "Please enter video description",
          });
          return false;
        }
        if (value.trim().length > VIDEO_DESC_CHAR_LIMIT) {
          this.setState({ error_description: "Description too long" });
          return false;
        }
        this.setState({ error_description: "" });
        return true;
      case "thumbnail":
        if (value && value.size > THUMBNAIL_MAX_SIZE) {
          this.setState({
            new_thumbnail: null,
            error_thumbnail: "Maximum thumbnail size is 50 MB",
          });
          return false;
        }
        this.setState({ error_thumbnail: "" });
        return true;
      default:
        return false;
    }
  }

  checkAvailability(value) {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      api.cast
        .titleCheck(value.trim())
        .then((response) => {
          if (response.data.code !== 200 && this.state.editing) {
            this.setState({ error_title: response.data.error });
          } else {
            this.setState({ error_title: "" });
          }
        })
        .catch((error) => {
          console.log(error);
          this.setState({ error_edit: "An error has occurred!" });
        });
    }, 400);
  }

  submitForm() {
    let ok = true;
    if (!this.state.attempted) {
      this.setState({ error_edit: "", attempted: true });
      ok &= this.validate("title", this.state.title);
      ok &= this.validate("description", this.state.description);
      ok &= this.validate("tags", this.state.tags);
    } else {
      ok &= !this.state.error_title;
      ok &= !this.state.error_description;
      ok &= !this.state.error_tags;
    }
    if (!ok) return;
    this.setState({ error_edit: "", loading_edit: true });
    const form = new FormData();
    form.append("hash", this.props.video.hash);
    form.append("title", this.state.title);
    form.append("description", this.state.description);
    form.append("tags", this.state.tags.map((tag) => tag.text).join(","));
    if (this.state.new_thumbnail)
      form.append("thumbnail", this.state.new_thumbnail);
    api.cast
      .edit(form)
      .then((response) => {
        clearTimeout(timeout);
        if (response.data.code === 200) {
          if (this.state.new_thumbnail)
            this.setState({
              thumbnail: URL.createObjectURL(this.state.new_thumbnail),
            });
          this.setState({
            editing: false,
            loading_edit: false,
            before: {},
            error_title: "",
            error_tags: "",
            error_description: "",
            error_thumbnail: "",
            new_thumbnail: "",
            updated: true,
          });
        } else {
          this.setState({
            error_edit: response.data.error,
            loading_edit: false,
          });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({
          error_edit: "An error has occurred! Please try again.",
          loading_edit: false,
        });
      });
  }

  deleteVideo() {
    this.setState({ error_delete: "", loading_delete: true });
    api.cast
      .remove(this.props.video.hash)
      .then((response) => {
        if (response.data.code === 200) {
          this.setState({ loading_delete: false, prompt: false });
          this.props.onDelete();
        } else {
          this.setState({
            error_delete: response.data.error,
            loading_delete: false,
          });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({
          error_delete: "An error has occurred! Please try again.",
          loading_delete: false,
        });
      });
  }

  openVideo() {
    switch (this.props.video.type) {
      case "live":
        if (this.props.video.is_live && this.props.video.hash !== currentHash())
          this.props.history.push(`/w/${this.props.video.hash}`);
        break;
      case "vod":
        if (
          this.props.video.resolutions &&
          this.props.video.hash !== currentHash()
        )
          this.props.history.push(`/w/${this.props.video.hash}`);
        break;
      default:
        console.log("Cannot open cast!");
    }
  }

  render() {
    return (
      <Card body style={style.card}>
        <Row>
          <Col xl={3} lg={4} md={5} sm={12}>
            {this.state.editing ? (
              <>
                <Dropzone
                  accept={"image/*"}
                  multiple={false}
                  noDragEventsBubbling={true}
                  onDrop={(files) => {
                    this.setState({ new_thumbnail: files[0] });
                    this.validate("thumbnail", files[0]);
                  }}
                  disabled={false}
                >
                  {({ getRootProps, getInputProps }) => (
                    <section
                      style={{
                        ...style.thumbnail_upload,
                        ...(this.state.new_thumbnail &&
                          style.thumbnail_upload_modified),
                      }}
                    >
                      <div
                        {...getRootProps()}
                        style={style.thumbnail_container}
                      >
                        <input {...getInputProps()} />
                        <Image
                          src={
                            this.state.new_thumbnail
                              ? URL.createObjectURL(this.state.new_thumbnail)
                              : this.state.thumbnail
                          }
                          style={style.thumbnail}
                        />
                        <p
                          style={{
                            ...style.dropzone_icon,
                            ...(this.state.new_thumbnail &&
                              style.dropzone_icon_modified),
                          }}
                        >
                          {!this.state.new_thumbnail && (
                            <span
                              className="material-icons"
                              style={{ fontSize: 32, color: "dimgray" }}
                            >
                              publish
                            </span>
                          )}
                        </p>
                      </div>
                    </section>
                  )}
                </Dropzone>
                <div
                  className={"invalid-feedback"}
                  style={{ display: "block" }}
                >
                  {this.state.error_thumbnail}
                </div>
              </>
            ) : (
              <div style={style.thumbnail_container}>
                <Image
                  src={this.state.thumbnail}
                  style={style.thumbnail}
                  onClick={this.openVideo}
                />
              </div>
            )}
          </Col>
          <Col md sm={12} style={{ marginTop: 4 }}>
            {this.state.error_edit && (
              <Alert variant={"danger"}>{this.state.error_edit}</Alert>
            )}
            {this.state.editing ? (
              <Form autoComplete={"off"} onSubmit={(e) => e.preventDefault()}>
                <Form.Group style={{ marginBottom: 4 }}>
                  <Form.Control
                    name={"title"}
                    value={this.state.title}
                    onBlur={this.handleChange}
                    onChange={this.handleChange}
                    type={"text"}
                    size={"lg"}
                    style={style.title}
                    isInvalid={this.state.error_title}
                    placeholder={"Title"}
                  />
                  <Form.Control.Feedback type={"invalid"}>
                    {this.state.error_title}
                  </Form.Control.Feedback>
                </Form.Group>
              </Form>
            ) : (
              <p style={style.title}>{this.state.title}</p>
            )}
            {this.props.video.type === "vod" && (
              <>
                <p style={{ marginBottom: 0 }}>
                  {format().full_date(this.props.video.created_at)}
                </p>
                <div style={style.cast_tag_bar}>
                  <div>
                    <Badge pill style={style.cast_tag_resolution}>
                      {resolutions[this.props.video.resolutions]}{" "}
                      {this.props.video.resolutions === 5 || (
                        <Spinner animation="grow" style={style.spinner} />
                      )}
                    </Badge>
                    <Badge pill style={style.cast_tag}>
                      {`${this.props.video.views} view${
                        this.props.video.views === 1 ? "" : "s"
                      }`}
                    </Badge>
                    <Badge pill style={style.cast_tag}>
                      {`${this.props.video.likes} like${
                        this.props.video.likes === 1 ? "" : "s"
                      }`}
                    </Badge>
                  </div>
                </div>
              </>
            )}
            <div style={style.cast_tag_bar}>
              {this.state.editing ? (
                <Form.Group style={{ width: "100%", marginBottom: 4 }}>
                  <ReactTags
                    classNames={{
                      tags: this.state.error_tags
                        ? "ReactTags__tags__error"
                        : this.state.uploading
                        ? "ReactTags__tags__disabled"
                        : "ReactTags__tags",
                      tagInput:
                        this.state.tags.length === VIDEO_TAG_COUNT
                          ? "ReactTags__tagInput__disabled"
                          : "ReactTags__tagInput",
                    }}
                    tags={this.state.tags}
                    autofocus={false}
                    delimiters={[13, 32, 188]}
                    maxLength={VIDEO_TAG_CHAR_LIMIT}
                    placeholder={this.state.tags.length ? "" : "Tags"}
                    readOnly={this.state.loading_edit}
                    handleAddition={this.handleTagAdd}
                    handleDelete={this.handleTagDelete}
                    handleDrag={this.handleTagDrag}
                    handleInputChange={() => this.setState({ error_tags: "" })}
                    handleInputBlur={() =>
                      this.validate("tags", this.state.tags)
                    }
                    handleTagClick={this.handleTagClick}
                  />
                  {this.state.error_tags && (
                    <div style={style.invalidText}>{this.state.error_tags}</div>
                  )}
                </Form.Group>
              ) : (
                <div>
                  {this.state.tags &&
                    Object.values(this.state.tags).map((tag) => (
                      <Badge pill style={style.cast_tag} key={tag.text}>
                        {tag.text}
                      </Badge>
                    ))}
                </div>
              )}
            </div>
            {this.state.editing ? (
              <Form autoComplete={"off"} onSubmit={(e) => e.preventDefault()}>
                <Form.Group style={{ marginBottom: 4 }}>
                  <Form.Control
                    name={"description"}
                    value={this.state.description}
                    onBlur={this.handleChange}
                    onChange={this.handleChange}
                    as={"textarea"}
                    size={"lg"}
                    style={style.description}
                    isInvalid={this.state.error_description}
                    rows={5}
                    placeholder={"Description"}
                  />
                  <Form.Control.Feedback type={"invalid"}>
                    {this.state.error_description}
                  </Form.Control.Feedback>
                </Form.Group>
              </Form>
            ) : (
              <p style={style.description}>{this.state.description}</p>
            )}
          </Col>
          <Col md={"auto"} sm={12}>
            <Button
              variant={"success"}
              block
              size={"sm"}
              style={style.button}
              onClick={this.pressEdit}
              disabled={
                this.state.before &&
                this.state.title === this.state.before.title &&
                this.state.tags === this.state.before.tags &&
                this.state.description === this.state.before.description &&
                !this.state.new_thumbnail
              }
            >
              {this.state.loading_edit ? (
                <Spinner animation="grow" style={style.spinner} />
              ) : (
                <span className="material-icons">
                  {this.state.editing ? "check" : "edit"}
                </span>
              )}
            </Button>
            {(this.state.editing || this.props.deletable) && (
              <Button
                variant={this.state.editing ? "info" : "danger"}
                block
                size={"sm"}
                onClick={this.pressDelete}
                style={style.button}
              >
                <span className="material-icons">
                  {this.state.editing ? "clear" : "delete"}
                </span>
              </Button>
            )}
          </Col>
        </Row>
        <Modal show={this.state.prompt} size={"md"} centered>
          <Modal.Header>
            <Modal.Title id="contained-modal-title-vcenter">
              Delete Video?
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {this.state.error_delete && (
              <Alert variant={"danger"}>{this.state.error_delete}</Alert>
            )}
            <p>
              Are you sure you want to delete <b>{this.state.title}</b>?{" "}
              <b>This action cannot be undone.</b>
            </p>
          </Modal.Body>
          <Modal.Footer>
            <Button
              variant={"info"}
              onClick={() => this.setState({ prompt: false })}
            >
              Cancel
            </Button>
            <Button variant={"danger"} onClick={this.deleteVideo}>
              Delete{" "}
              {this.state.loading_delete ? (
                <Spinner animation="grow" style={style.spinner} />
              ) : (
                ""
              )}
            </Button>
          </Modal.Footer>
        </Modal>
        <Prompt
          when={this.state.editing || this.state.loading_edit}
          message={() => `Are you sure? You changes will be discarded.`}
        />
      </Card>
    );
  }
}

let style = {
  card: {
    borderRadius: "8px 8px 8px 8px",
    marginBottom: 16,
    width: "100%",
  },
  title: {
    lineHeight: 1.2,
    margin: 0,
    fontSize: "2rem",
    fontWeight: 600,
    marginBottom: 4,
  },
  thumbnail_upload: {
    background: "#f0f0f088",
    border: "3px dashed #ddddddaa",
    flexShrink: 0,
    borderRadius: "8px 48px 8px 8px",
    zIndex: 100,
  },
  thumbnail_upload_modified: {
    background: "#f0f0f022",
    border: "",
  },
  dropzone_icon: {
    borderRadius: "8px 48px 8px 8px",
    background: "#f0f0f088",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    width: "100%",
    height: "100%",
    position: "absolute",
    top: 0,
    margin: 0,
  },
  dropzone_icon_modified: {
    background: "#f0f0f022",
    border: "",
  },
  thumbnail_container: {
    height: 0,
    paddingBottom: "56.25%",
    position: "relative",
  },
  thumbnail: {
    borderRadius: "8px 48px 8px 8px",
    objectFit: "cover",
    width: "100%",
    height: "100%",
    position: "absolute",
    cursor: "pointer",
  },
  created_at: {
    fontSize: 16,
    fontWeight: 300,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  cast_tag_bar: {
    display: "flex",
    justifyContent: "space-between",
  },
  cast_tag: {
    background: "#8B2803AA",
    color: "#DDD",
    borderRadius: 8,
    // borderWidth: 1,
    // borderColor: "lightgray",
    // borderStyle: "solid",
    fontSize: 14,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_tag_resolution: {
    background: "#E84409",
    color: "white",
    marginRight: 8,
    marginBottom: 8,
    borderRadius: 8,
    fontSize: 14,
    fontWeight: 400,
  },
  description: {
    fontSize: 16,
    WebkitLineClamp: "3",
    overflow: "hidden",
    WebkitBoxOrient: "vertical",
    display: "-webkit-box",
  },
  spinner: {
    width: 10,
    height: 10,
    verticalAlign: "initial",
  },
  invalidInput: {
    borderColor: "#dc3545",
    paddingRight: "calc(1.5em + .75rem)",
    backgroundImage:
      "url(data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='12' hei…circle cx='6' cy='8.2' r='.6' fill='%23dc3545' stroke='none'/%3e%3c/svg%3e)",
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
  },
};

export default withRouter(CastEditable);
