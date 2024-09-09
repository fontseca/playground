const requestTarget = document.getElementById("http-request-target");

document.body.addEventListener('htmx:afterOnLoad', HandleAfterOnLoad);

document.addEventListener("DOMContentLoaded", function () {
  AppendQueryParameterRow("", "");
  AppendHeaderRow("", "");

  if (requestTarget) {
    requestTarget.addEventListener("keyup", ParseQueryParametersFromRequestBar);
  }

  requestTarget.value = localStorage.getItem("fontseca.dev/playground@http-request-target");
  ParseQueryParametersFromRequestBar();
});

function ParseHTTPResponse(httpResponseMessage) {
  const result = {
    proto: "",
    statusCode: 0,
    statusText: "",
    headers: {},
    body: ""
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
    result.headers[key] = value;
  }

  result.body = httpResponseMessage.substring(2 + endOfHeaders);
  return result;
}

function SubmitRequest(event) {
  if (event.ctrlKey && "Enter" === event.key) {
    event.preventDefault();
    htmx.trigger(document.forms["http-request-form"], "submit");
  }
}

function HandleAfterOnLoad(event) {
  const statusContainer = document.getElementById("http-response-status");
  const bodyContainer = document.getElementById("http-response-body");
  const responseHeadersTable = document.getElementById("http-response-headers");
  const httpResponseMessage = event.detail.xhr.responseText;
  const response = ParseHTTPResponse(httpResponseMessage);

  responseHeadersTable.innerHTML = "";
  bodyContainer.innerHTML = response.body;
  statusContainer.textContent = `${response.proto} ${response.statusCode} ${ response.statusText }`;

  Object.keys(response.headers).forEach(key => {
    const value = response.headers[key];
    const entry = responseHeadersTable.insertRow();
    entry.innerHTML = `
    <td>
      <input class="http-response-header-key"
             type="text"
             value="${key}"
             readonly/>
    </td>
    <td>
      <input class="http-response-header-value"
             type="text"
             value="${value}"
             style="width: 500px"
             readonly/>
    </td>
  `;
  })
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
