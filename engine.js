const requestTabs = document.querySelectorAll("[data-tab-request-target]")
const requestTabContents = document.querySelectorAll("[data-tab-request-content]")

requestTabs.forEach(tab => {
  tab.addEventListener('click', () => {
    const target = document.querySelector(tab.dataset.tabRequestTarget)
    requestTabContents.forEach(tabContent => {
      tabContent.classList.remove('active')
    })

    requestTabs.forEach(tab => {
      tab.classList.remove('active')
    })

    tab.classList.add('active')
    target.classList.add('active')
  })
})

const responseTabs = document.querySelectorAll('[data-tab-response-target]')
const responseTabContents = document.querySelectorAll('[data-tab-response-content]')

responseTabs.forEach(tab => {
  tab.addEventListener('click', () => {
    const target = document.querySelector(tab.dataset.tabResponseTarget)

    if ("tab-response-body" === target.getAttribute("id")) {
      target.parentElement.classList.remove("with-diagonal");
    } else {
      target.parentElement.classList.add("with-diagonal");
    }

    responseTabContents.forEach(tabContent => {
      tabContent.classList.remove('active')
    })

    responseTabs.forEach(tab => {
      tab.classList.remove('active')
    })

    tab.classList.add('active')
    target.classList.add('active')
  })
})


let requestStarts;
const requestTarget = document.getElementById("http-request-target");

document.body.addEventListener('htmx:afterOnLoad', HandleAfterOnLoad);

document.addEventListener("DOMContentLoaded", function () {
  AppendQueryParameterRow("", "");
  AppendHeaderRow("", "");

  let alreadyLoaded = false;

  try {
    let target = new URLSearchParams(window.location.search).get("target");

    if (target) {
      if ("/" === target[0]) {
        target = `${window.location.protocol}//${window.location.host}${target}`
      }

      const url = new URL(target)
      requestTarget.value = url.toString();
      alreadyLoaded = true;
      new URLSearchParams(window.location.search).delete("target");
    }
  } catch (e) {
    console.error(e);
  }

  if (requestTarget) {
    requestTarget.addEventListener("keyup", ParseQueryParametersFromRequestBar);
  }

  if (!alreadyLoaded) {
    requestTarget.value = localStorage.getItem("fontseca.dev/playground@http-request-target");
  }

  ParseQueryParametersFromRequestBar();
});


function ParseCookie(str) {
  const obj = {};
  const pairs = str.split(/; */);

  let mainCookieSet = false; // Flag for tracking if Name/Value has been set.

  for (const pair of pairs) {
    const equalsIndex = pair.indexOf('=');

    let key = equalsIndex > -1 ? pair.slice(0, equalsIndex).trim() : pair.trim();
    let value = equalsIndex > -1 ? pair.slice(equalsIndex + 1).trim() : 'true';

    if (value.length > 1 && value[0] === '"' && value[value.length - 1] === '"') {
      value = value.slice(1, -1);
    }

    if (!mainCookieSet) {
      obj.Name = key;
      obj.Value = decodeURIComponent(value);
      mainCookieSet = true;
    } else {
      if (!Object.hasOwn(obj, key)) {
        obj[key] = decodeURIComponent(value);
      }
    }
  }

  return obj;
}


function ParseHTTPResponse(httpResponseMessage) {
  const result = {
    proto: "",
    statusCode: 0,
    statusText: "",
    headers: [],
    body: "",
    cookies: [],
  };

  const endOfStartLine = httpResponseMessage.indexOf("\n");
  const startLine = httpResponseMessage.substring(0, endOfStartLine);
  const endOfStartLineProto = startLine.indexOf(" ");
  result.proto = startLine.substring(0, endOfStartLineProto);

  const status = startLine.substring(1 + endOfStartLineProto);
  result.statusCode = parseInt(status.substring(0, status.indexOf(" ")));
  result.statusText = status.substring(1 + status.indexOf(" "));

  const endOfHeaders = httpResponseMessage.indexOf("\n\n");
  const headers = httpResponseMessage.substring(1 + endOfStartLine, endOfHeaders);

  for (const line of headers.split("\n")) {
    const [key, value] = line.split(": ");

    if (key.toLowerCase() === "set-cookie") {
      result.cookies.push(ParseCookie(value));
    }

    result.headers.push({key, value});
  }

  result.body = httpResponseMessage.substring(2 + endOfHeaders);
  return result;
}

document.getElementById("http-request-form").onsubmit = (ev) => {
  if (!ev.target.children[1].checkValidity()) {
    return
  }

  document.getElementById("http-response-body").innerHTML = "";
  document.getElementById("centered-label-response-body").classList.add("disable");
  const responseStatus = document.getElementById("response-status");
  const responseStats = document.getElementById("response-stats");
  responseStatus.getElementsByTagName("a")[0].textContent = `— —`;

  responseStats.getElementsByTagName("span")[0].textContent = `— MS`;
  responseStats.getElementsByTagName("span")[1].textContent = `— KB`;

  requestStarts = Date.now();
}

function SubmitRequest(event) {
  if (event.ctrlKey && "Enter" === event.key) {
    event.preventDefault();
    htmx.trigger(document.forms["http-request-form"], "submit");
  }
}

function HandleAfterOnLoad(event) {
  const bodyContainer = document.getElementById("http-response-body");
  const responseHeadersTable = document.getElementById("http-response-headers");
  const responseCookiesTable = document.getElementById("http-response-cookies");
  const httpResponseMessage = event.detail.xhr.responseText;
  const response = ParseHTTPResponse(httpResponseMessage);

  responseHeadersTable.innerHTML = "";
  responseCookiesTable.innerHTML = "";
  bodyContainer.innerHTML = response.body;
  const responseStatus = document.getElementById("response-status");
  const responseStats = document.getElementById("response-stats");

  responseStatus.classList.add("active");
  responseStats.classList.add("active");

  const statusAnchor = responseStatus.getElementsByTagName("a")[0];
  statusAnchor.setAttribute("href", `https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/${response.statusCode}`)
  statusAnchor.setAttribute("title", `Read more about the \`${response.statusCode} ${response.statusText}\` response.`)
  statusAnchor.textContent = `${response.statusCode} ${response.statusText}`;

  responseStats.getElementsByTagName("span")[0].textContent = `${Date.now() - requestStarts} MS`;
  responseStats.getElementsByTagName("span")[1].textContent = `${response.body.length / 1000} KB`;

  document.querySelector("li[data-tab-response-target='#tab-response-headers']").textContent = `Headers (${response.headers.length})`;

  document.querySelectorAll("#tab-response-headers .disable").forEach(e => e.classList.remove("disable"));
  document.getElementById("centered-label-response-headers").classList.add("disable");

  response.headers.forEach(header => {
    const entry = responseHeadersTable.insertRow();
    entry.innerHTML = `
    <td>
      <input class="http-response-header-key"
             type="text"
             value="${header.key}"
             readonly/>
    </td>
    <td>
      <input class="http-response-header-value"
             type="text"
             value="${header.value}"
             readonly/>
    </td>
  `;
  })

  if (response.cookies.length > 0) {
    document.querySelectorAll("#tab-response-cookies .disable").forEach(e => e.classList.remove("disable"));
    document.getElementById("centered-label-response-cookies").classList.add("disable");
    document.querySelector("li[data-tab-response-target='#tab-response-cookies']").textContent = `Cookies (${response.cookies.length})`;
    response.cookies.forEach(cookie => {
      const entry = responseCookiesTable.insertRow();
      entry.innerHTML = `
      <td><input type="text" value="${cookie.Name}" readonly/></td>
      <td><input type="text" value="${cookie.Value}" readonly/></td>
      <td><input type="text" value="${cookie.Domain ?? "—"}" readonly/></td>
      <td><input type="text" value="${cookie.Path ?? "—"}" readonly/></td>
      <td><input type="text" value="${cookie.Expires ?? "Session"}" readonly/></td>
      <td><input type="text" value="${cookie.HttpOnly ?? "—"}" readonly/></td>
      <td><input type="text" value="${cookie.Secure ?? "—"}" readonly/></td>
    `;
    })
  } else {
    document.querySelector("li[data-tab-response-target='#tab-response-cookies']").textContent = `Cookies`;
    document.querySelector("#tab-response-cookies h3").classList.add("disable");
    document.querySelector("#tab-response-cookies table").classList.add("disable");
    document.getElementById("centered-label-response-cookies").classList.remove("disable");
    document.getElementById("centered-label-response-cookies").children[0].textContent = "No cookies received from the server."
  }
}

function StoreRequestTargetURL() {
  localStorage.setItem("fontseca.dev/playground@http-request-target", requestTarget.value.trim());
}

function GetQueryParametersTable() {
  return document.getElementById("http-request-query-parameters");
}

function AppendQueryParameterRow(key, value) {
  const tbody = GetQueryParametersTable();
  const entry = tbody.insertRow();
  entry.innerHTML = `
    <td>
      <input class="http-request-query-parameter-key"
             type="text"
             placeholder="Key"
             form="http-request-form"
             minlength="1"
             maxlength="64"
             value="${key}" />
    </td>
    <td>
      <input class="http-request-query-parameter-value"
             type="text"
             form="http-request-form"
             placeholder="Value"
             minlength="1"
             maxlength="64"
             value="${value}" />
    </td>
  `;

  entry
    .querySelector(".http-request-query-parameter-key")
    .addEventListener("keyup", InterceptQueryParameterEntry);

  entry
    .querySelector(".http-request-query-parameter-value")
    .addEventListener("keyup", UpdateQueryParametersInTarget);
}

function ParseQueryParametersFromRequestBar() {
  StoreRequestTargetURL();

  const url = requestTarget.value;
  const at = url.indexOf("?");
  const parametersTable = GetQueryParametersTable();

  if (!~at) { // no query parameters at all
    if (parametersTable.children.length > 1) { // clear superfluous parameters
      while (parametersTable.firstChild && parametersTable.firstChild !== parametersTable.lastChild) {
        parametersTable.removeChild(parametersTable.firstChild);
      }
    }

    return;
  }

  const parameters = url.substring(1 + at).split("&");

  for (let i = 0; i < parameters.length; i++) {
    const parameter = parameters[i].trim();
    const key = parameter.split("=")[0];
    const index = parameter.indexOf("=");
    let value = "";

    if (index > 0) {
      value = parameter.substring(parameter.indexOf("=") + 1);
    }

    if (i < parametersTable.children.length) { // reuse current rows in table
      const row = parametersTable.children.item(i);
      row.querySelector(".http-request-query-parameter-key").value = key;
      row.querySelector(".http-request-query-parameter-value").value = value;
      continue;
    }

    AppendQueryParameterRow(key, value);
  }

  if (parameters.length === parametersTable.children.length) { // query parameters table needs a new row
    AppendQueryParameterRow("", "");
  }

  let i = 0;
  if (0 !== parameters.length) {
    i = parameters.length;
  }

  for (; i < parametersTable.children.length - 1; i++) { // discard last rows
    parametersTable.removeChild(parametersTable.children.item(i));
  }
}

function UpdateQueryParametersInTarget() {
  StoreRequestTargetURL();

  const parametersTable = GetQueryParametersTable();
  let queryParameters = "?";
  let willClearParameters = true;

  for (const row of parametersTable.children) {
    const key = row.querySelector(".http-request-query-parameter-key").value.trim();

    if ("" === key) {
      continue;
    }

    willClearParameters = false;

    const value = row.querySelector(".http-request-query-parameter-value").value.trim();

    if (queryParameters !== "?") {
      queryParameters += "&";
    }

    queryParameters += key + "=" + value;
  }

  let url = requestTarget.value.trim();
  const base = url.split("?")[0];

  if (willClearParameters) {
    requestTarget.value = base;
  }

  if ("?" !== queryParameters) {
    requestTarget.value = base + queryParameters;
  }
}

function InterceptQueryParameterEntry(event) {
  UpdateQueryParametersInTarget();

  const parameterKeyInputElement = event.target;
  const parametersTable = GetQueryParametersTable();
  const parametersCount = parametersTable.rows.length;
  const parameterKey = parameterKeyInputElement.value.trim();
  const indexOfCurrentParamRow = parameterKeyInputElement
    .parentElement
    .parentElement
    .rowIndex;

  if (indexOfCurrentParamRow === parametersCount - 1) { // we're at the penultimate row
    const parameterValue = parameterKeyInputElement
      .parentElement
      .parentElement
      .querySelector(".http-request-query-parameter-value")
      .value
      .trim();

    if ("" === parameterKey && "" === parameterValue) {
      parametersTable.removeChild(parametersTable.lastChild);
    }
  }

  const needsNewParamRow = parametersCount < 1 + indexOfCurrentParamRow;
  if ("" !== parameterKey && needsNewParamRow) {
    AppendQueryParameterRow("", "");
  }
}

function GetHeadersTable() {
  return document.getElementById("http-request-headers");
}

function AppendHeaderRow(key, value) {
  const tbody = GetHeadersTable();
  const entry = tbody.insertRow();
  entry.innerHTML = `
    <td>
      <input class="http-request-header-key"
             type="text"
             placeholder="Key"
             form="http-request-form"
             name="header-key"
             value="${key}" />
    </td>
    <td>
      <input class="http-request-header-value"
             type="text"
             form="http-request-form"
             name="header-value"
             placeholder="Value"
             value="${value}" />
    </td>
  `;

  entry
    .querySelector(".http-request-header-key")
    .addEventListener("keyup", InterceptHeaderEntry);

  entry.querySelector(".http-request-header-value");
}

function InterceptHeaderEntry(event) {
  const headerKeyInputElement = event.target;
  const headersTable = GetHeadersTable();
  const headersCount = headersTable.rows.length;
  const headerKey = headerKeyInputElement.value.trim();
  const indexOfCurrentHeaderRow = headerKeyInputElement
    .parentElement
    .parentElement
    .rowIndex;

  if (indexOfCurrentHeaderRow === headersCount - 1) { // we're at the penultimate row
    const headerValue = headerKeyInputElement
      .parentElement
      .parentElement
      .querySelector(".http-request-header-value")
      .value
      .trim();

    if ("" === headerKey && "" === headerValue) {
      headersTable.removeChild(headersTable.lastChild);
    }
  }

  const needsNewHeaderRow = headersCount < 1 + indexOfCurrentHeaderRow;
  if ("" !== headerKey && needsNewHeaderRow) {
    AppendHeaderRow("", "");
  }
}
