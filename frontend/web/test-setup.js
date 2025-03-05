if (typeof this.global.TextEncoder === "undefined") {
  const { TextEncoder, TextDecoder } = require("util");
  this.global.TextEncoder = TextEncoder;
  this.global.TextDecoder = TextDecoder;
}
