import React from "react";
import { Badge, Image } from "react-bootstrap";
import { useHistory } from "react-router-dom";
import abbreviate from "../helper/abbreviate";
import urls from "../helper/url";
import TimeAgo from "react-timeago";

function Cast(props) {
  const history = useHistory();

  function playVideo() {
    if (props.onClick) props.onClick(props.video.type, props.video.hash);
    if (props.video.hash !== urls().current_hash())
      history.push(`/w/${props.video.hash}`);
  }

  function viewAuthor(e) {
    e.stopPropagation();
    console.log(`View: ${props.video.author.name}`);
  }

  return (
    <div style={style.cast_card} onClick={playVideo}>
      <Image
        src={urls().thumbnail(props.video.hash)}
        style={style.cast_thumbnail}
      />
      <div style={style.cast_tag_bar}>
        {props.video.type === "live" && (
          <Badge pill style={style.cast_live_tag}>
            LIVE
          </Badge>
        )}
        <Badge pill style={style.cast_viewer_tag}>
          {abbreviate().number(props.video.views)}{" "}
          {props.video.type === "live" ? "viewers" : "views"}
        </Badge>
      </div>
      <div style={style.cast_detail}>
        <Image
          src={urls().profile(props.video.author.username)}
          height={42}
          width={42}
          style={style.profile_image}
          onClick={viewAuthor}
          roundedCircle
        />
        <div style={style.cast_author_details}>
          <p style={style.cast_title}>{props.video.title}</p>
          <p style={style.cast_author} onClick={viewAuthor}>
            {props.video.author.name}
          </p>
          <p style={style.cast_duration}>
            {props.video.type === "live" ? (
              <TimeAgo
                date={props.video.created_at}
                formatter={(value, unit, _) => {
                  return `Streaming for ${value} ${unit}${
                    value === 1 ? "" : "s"
                  }`;
                }}
              />
            ) : (
              <TimeAgo date={props.video.created_at} />
            )}
          </p>
        </div>
      </div>
    </div>
  );
}

let style = {
  cast_card: {
    cursor: "pointer",
  },
  cast_thumbnail: {
    borderRadius: "8px 48px 8px 8px",
    // borderWidth: 1,
    // borderColor: "lightgray",
    // borderStyle: "solid",
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
    fontWeight: 400,
  },
  cast_viewer_tag: {
    background: "#8B2803AA",
    color: "#DDD",
    borderRadius: 8,
    // borderWidth: 1,
    // borderColor: "lightgray",
    // borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_detail: {
    display: "flex",
  },
  cast_author_details: {
    display: "grid",
  },
  cast_title: {
    fontSize: 18,
    fontWeight: 600,
    margin: 0,
  },
  cast_author: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0,
  },
  cast_duration: {
    fontSize: 16,
    fontWeight: 300,
    margin: 0,
    opacity: 0.7,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "start",
  },
};

export default Cast;
