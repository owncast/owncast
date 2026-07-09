var _defs;
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }
import * as React from "react";
var SvgAirplay = function SvgAirplay(props) {
  return /*#__PURE__*/React.createElement("svg", _extends({
    xmlns: "http://www.w3.org/2000/svg",
    xmlnsXlink: "http://www.w3.org/1999/xlink",
    xmlSpace: "preserve",
    viewBox: "0 0 46 42"
  }, props), _defs || (_defs = /*#__PURE__*/React.createElement("defs", null, /*#__PURE__*/React.createElement("path", {
    id: "airplay_svg__a",
    d: "M22.2 24.2 8.5 39.9c-.5.6-.1 1.5.7 1.5h27.5c.8 0 1.2-.9.7-1.5L23.8 24.2c-.2-.2-.5-.4-.8-.4s-.6.1-.8.4M6.5.6c-2.3 0-3.1.2-3.9.7-.8.4-1.5 1.1-1.9 1.9C.2 4 0 4.9 0 7.1v17.5c0 2.3.2 3.1.7 3.9.4.8 1.1 1.5 1.9 1.9s1.7.7 3.9.7h5.2l2.2-2.6H5.8c-1.1 0-1.6-.1-2-.3s-.7-.5-1-1c-.2-.4-.3-.8-.3-2v-19c0-1.1.1-1.6.3-2s.5-.7 1-1c.4-.2.8-.3 2-.3h34.3c1.1 0 1.6.1 2 .3s.7.5 1 1c.2.4.3.8.3 2v19c0 1.1-.1 1.6-.3 2s-.5.7-1 1c-.4.2-.8.3-2 .3H32l2.2 2.6h5.2c2.3 0 3.1-.2 3.9-.7.8-.4 1.5-1.1 1.9-1.9s.7-1.7.7-3.9V7.1c0-2.3-.2-3.1-.7-3.9-.4-.8-1.1-1.5-1.9-1.9S41.6.6 39.4.6z"
  }))), /*#__PURE__*/React.createElement("clipPath", {
    id: "airplay_svg__b"
  }, /*#__PURE__*/React.createElement("use", {
    xlinkHref: "#airplay_svg__a",
    style: {
      overflow: "visible"
    }
  })), /*#__PURE__*/React.createElement("path", {
    d: "M-6.4-5.8h58.7v53.6H-6.4z",
    style: {
      clipPath: "url(#airplay_svg__b)"
    }
  }));
};
export default SvgAirplay;