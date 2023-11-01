const http = require("k6/http");
const { check, sleep } = require("k6");

module.exports.options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: "constant-arrival-rate",
      rate: 100,
      timeUnit: "1s",
      duration: "20s",
      preAllocatedVUs: 50,
      maxVUs: 1200,
    },
  },
};

// test HTTP
module.exports.default = function () {
  const res = http.get("http://localhost:8080/api");
  check(res, { "status was 200": (r) => r.status === 200 });
  sleep(1);
};
