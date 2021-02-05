import React from "react";

function ProfileImage(props) {
  return (
    <div
      className="text-center"
      style={{
        width: props.size,
        height: props.size,
        flexShrink: 0,
        borderRadius: props.size / 2,
        background: "gray",
        color: "white",
        fontSize: `${props.size / 2}px`,
        lineHeight: `${props.size}px`,
        textAlign: "center",
        textTransform: "capitalize",
        ...props.style,
      }}
      onClick={props.onClick}
    >
      {props.name
        ?.split(" ")
        .slice(0, 2)
        .map((t) => t[0]?.toUpperCase())
        .join("")}
    </div>
  );
}

export default ProfileImage;
