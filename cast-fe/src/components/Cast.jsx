import React from "react";
import { Badge, Image } from "react-bootstrap";
import { Link } from "react-router-dom";
import { ProfileImage } from "./index";
import TimeAgo from "react-timeago";
import abbreviate from "../helper/abbreviate";
import api from "../apis/api";

function Cast(props) {
  return (
    <Link to={`/w/${props.video.hash}`} style={{ textDecoration: "none" }}>
      <div style={style.cast_card}>
        <Image
          src={api.cdn.thumbnail(props.video.hash)}
          style={style.cast_thumbnail}
        />
        <div style={style.cast_tag_bar}>
          {props.video.type === "live" && (
            <Badge pill style={style.cast_live_tag}>
              LIVE
            </Badge>
          )}
          {props.video.unlisted && (
            <Badge
              pill
              style={{ ...style.cast_tag, ...style.cast_tag_unlisted }}
            >
              <i className="fas fa-lock" /> Unlisted
            </Badge>
          )}
          <Badge pill style={style.cast_tag}>
            {abbreviate().number(props.video.views)}{" "}
            {props.video.type === "live" ? "viewers" : "views"}
          </Badge>
        </div>
        <div style={style.cast_detail}>
          <ProfileImage
            size={42}
            name={props.video.author.name}
            style={style.profile_image}
          />
          <div style={style.cast_author_details}>
            <p style={style.cast_title}>{props.video.title}</p>
            <p style={style.cast_author}>{props.video.author.name}</p>
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
    </Link>
  );
}

let style = {
  cast_card: {
    cursor: "pointer",
  },
  cast_thumbnail: {
    borderRadius: "8px 48px 8px 8px",
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
  cast_tag: {
    background: "#8B2803AA",
    color: "#DDD",
    borderRadius: 8,
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_tag_unlisted: {
    background: "rgb(3,69,139)",
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
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  cast_author: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0,
    overflow: "hidden",
    textOverflow: "ellipsis",
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
