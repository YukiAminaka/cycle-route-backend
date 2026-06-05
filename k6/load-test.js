import http from "k6/http";
import { check, fail } from "k6";
import { SharedArray } from "k6/data";
import { URL } from "https://jslib.k6.io/url/1.0.0/index.js";
import papaparse from "https://jslib.k6.io/papaparse/5.1.1/index.js";

const users = new SharedArray("users", function () {
  return papaparse.parse(open("./user.csv"), { header: true, skipEmptyLines: true }).data;
});

export const options = {
  thresholds: {
    http_req_duration: ["p(95)<1000"], // 95% のリクエストは 1000ms (1s) 以内に収める
  },
  scenarios: {
    contacts: {
      executor: "ramping-arrival-rate", // https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate
      exec: "load_test",

      gracefulStop: "10s",

      preAllocatedVUs: 50,
      stages: [
        // target: 1 秒あたりの load_test 関数の実行回数の目標値
        // duration: target 到達までにかかる時間
        { target: 20, duration: "1m" }, // 1分かけて target 20 に到達
        { target: 20, duration: "1m" }, // target 20 を1分間維持
      ],
    },
  },
};

/**
 * シナリオの内容
 * 1. /explore エンドポイントにリクエストを送る
 * 2. レスポンスからルートのIDを取得し、/routes/{routeId} エンドポイントにリクエストを送る
 * 4. 各リクエストのレスポンスコードが200であることを確認する
 */

export function load_test() {
  const exploreURL = new URL(`${__ENV.API_BASE_URL}/explore`);
  const offset = Math.floor(Math.random() * 100);
  exploreURL.searchParams.append(`limit`, `20`);
  exploreURL.searchParams.append(`offset`, `${offset}`);

  let res = http.get(exploreURL.toString());

  check_status_ok(res);

  let body = res.json();

  const routeRequests = Array();
  body.routes.forEach((route) => {
    const url = new URL(`${__ENV.API_BASE_URL}/routes/${route.id}`);
    routeRequests.push(["GET", url.toString()]);
  });

  let routeRes = http.batch(routeRequests);
  routeRes.forEach((res) => check_status_ok(res));
}

function check_status_ok(res) {
  const result = check(res, {
    "is status OK": (r) => r.status === 200,
  });
  if (!result) {
    fail(`status is not 200. response: ${JSON.stringify(res)}`);
  }
}

function login() {
  const user = users[__VU % users.length];

  // Step 1: ログインフローを作成する
  const flowRes = http.get(`${__ENV.KRATOS_PUBLIC_URL}/self-service/login/api`, {
    headers: { Accept: "application/json" },
  });

  check_status_ok(flowRes);

  const actionUrl = flowRes.json("ui.action");

  // Step 2: ログインフローを送信する
  const loginRes = http.post(
    actionUrl,
    JSON.stringify({
      method: "password",
      identifier: user.identifier,
      password: user.password,
    }),
    { headers: { "Content-Type": "application/json", Accept: "application/json" } },
  );

  check(loginRes, { "login successful": (r) => r.status === 200 });

  return loginRes.json("session_token");
}
