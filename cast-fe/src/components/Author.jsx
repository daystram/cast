import React from "react";
import { Button, Col, Image } from "react-bootstrap";
import abbreviate from "../helper/abbreviate";

function Author(props) {
  function handleSubscribe(e) {
    console.log("subscribe");
  }

  return (
    <>
      <Col xs={"auto"} style={{ display: "flex" }}>
        <Image
          src={""}
          width={128}
          height={128}
          style={style.profile_image}
          roundedCircle
        />
      </Col>
      <Col xl={5} lg={6} sm={6} xs={"auto"} style={{ display: "flex" }}>
        <div style={style.cast_author_details}>
          <p style={style.cast_author_name}>{props.author.name}</p>
          <p style={style.cast_author_sub}>
            {abbreviate().number(props.author.subscribers) || 0} subscriber
            {props.author.subscribers !== 1 && "s"}
          </p>
        </div>
      </Col>
      <Col
        sm={"auto"}
        xs={"none"}
        className={"xs-hidden"}
        style={{ display: "flex" }}
      >
        <Button
          style={style.sub_button}
          onClick={handleSubscribe}
          disabled={props.author.isSubscribed}
        >
          SUBSCRIBE
        </Button>
      </Col>
    </>
  );
}

let style = {
  sub_button: {
    fontWeight: 600,
    height: "fit-content",
    alignSelf: "center",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "start",
  },
  cast_author_details: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    marginLeft: 18,
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

export default Author;
