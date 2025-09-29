import React from "react";
import "./Tooltip.css";

interface TooltipProps {
  text: string;
  children: React.ReactElement;
}

const Tooltip: React.FC<TooltipProps> = ({ text, children }) => {
  const child = React.Children.only(children);
  return React.cloneElement(child, {
    "data-tooltip": text,
    className: `${child.props.className || ""} tooltip`,
  });
};

export default Tooltip;
