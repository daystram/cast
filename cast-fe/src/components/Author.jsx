import React from 'react';
import {Button, Col, Image} from "react-bootstrap";
import abbreviate from '../helper/abbreviate'
import urls from "../helper/url";

function Author(props) {

  function handleSubscribe(e) {
    console.log("subscribe")
  }

  return (
    <>
      <Col xs={"auto"} style={{display: "flex"}}>
        <Image src={urls().profile(props.author.username)}
               width={128} height={128} style={style.profile_image} roundedCircle/>
      </Col>
      <Col xl={5} lg={6} sm={6} xs={"auto"} style={{display: "flex"}}>
        <div style={style.cast_author_details}>
          <p style={style.cast_author_name}>{props.author.name}</p>
          <p style={style.cast_author_sub}>
            {(abbreviate().number(props.author.subscribers)) || 0} subscriber{props.author.subscribers !== 1 && "s"}
          </p>
        </div>
      </Col>
      <Col sm={"auto"} xs={"none"} className={"xs-hidden"} style={{display: "flex"}}>
        <Button style={style.sub_button} onClick={handleSubscribe}
                disabled={props.author.isSubscribed}>SUBSCRIBE</Button>
      </Col>
    </>
  )
}

let style = {
  cast_card: {
    cursor: "pointer"
  },
  cast_thumbnail: {
    borderRadius: "8px 48px 8px 8px",
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    objectFit: "cover",
    width: "100%",
  },
  cast_tag_bar: {
    marginTop: 8,
  },
  cast_live_tag: {
    background: "red",
    color: "white",
    marginRight: 8,
    marginBottom: 8,
    borderRadius: 8,
    fontSize: 16,
    fontWeight: 400
  },
  cast_viewer_tag: {
    background: "white",
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_detail: {
    display: "flex",
  },
  // cast_author_details: {
  //   display: "grid"
  // },
  sub_button: {
    fontWeight: 600,
    height: "fit-content",
    alignSelf: "center"
  },
  cast_title: {
    fontSize: 18,
    fontWeight: 600,
    margin: 0
  },
  cast_author: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0
  },
  cast_duration: {
    fontSize: 16,
    fontWeight: 300,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "start"
  },

  author_bar: {
    display: "flex",
    justifyContent: "space-between"
  },
  author_profile: {
    display: "flex",
    width: "100%"
  },
  cast_author_details: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    marginLeft: 18
  },
  cast_author_name: {
    fontSize: 28,
    fontWeight: 600,
    margin: 0,
  },
  cast_author_sub: {
    fontSize: 22,
    fontWeight: 400,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
};

export default Author
