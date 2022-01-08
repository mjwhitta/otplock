package otplock

const advancedDashboard string = `
<html>
  <body>
    <center>
      <form method="post">
        <table style="width:100%">
          <col style="width:10%">
          <col style="width:90%">

          <tr>
            <td style="text-align:right; width=20">
              <label for="expires">Time (secs):</label>
            </td>
            <td>
              <input
                id="expires"
                name="expires"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="payload">Payload (hex):</label>
            </td>

            <td>
              <textarea
                id="payload"
                name="payload"
                rows="32"
                style="width:100%"
              ></textarea>
            </td>
          </tr>
        </table>

        <input
          hidden="true"
          id="type"
          name="type"
          type="text"
          value="simple"
        />
        <input type="submit" value="Submit"/>
      </form>
    </center>
  </body>
</html>
`

const advancedNew string = `
<html>
  <body>
    <center>` + navbar + `<br>
      <form method="post">
        <table style="width:100%">
          <col style="width:10%">
          <col style="width:90%">

          <tr>
            <td style="text-align:right; width=20">
              <label for="endpoint">Endpoint:</label>
            </td>
            <td>
              <input
                id="endpoint"
                name="endpoint"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="filename">Source filename:</label>
            </td>
            <td>
              <input
                id="filename"
                name="filename"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="compile">Compile command:</label>
            </td>
            <td>
              <input
                id="compile"
                name="compile"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="binary">Binary:</label>
            </td>
            <td>
              <input
                id="binary"
                name="binary"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="source">Template source:</label>
            </td>

            <td>
              <textarea
                id="source"
                name="source"
                rows="32"
                style="width:100%"
              ></textarea>
            </td>
          </tr>
        </table>

        <input
          hidden="true"
          id="type"
          name="type"
          type="text"
          value="advanced"
        />
        <input type="submit" value="Submit"/>
      </form>
    </center>
  </body>
</html>
`

// This is used for with Fprintf, so must escape %
const advancedResp string = `
<html>
  <body>
    <center>
      Save this URL for later use:<br><br>%s<br><br>
      <form action="%s" method="get">
        <input type="submit" value="Continue"/>
      </form>
    </center>
  </body>
</html>
`

// This is used for errors with Fprintf, so must escape %
const errPg string = `
<html>
  <body>
    <center>
      %s<br>
      <br>
      <form method="get">
        <input
          hidden="true"
          id="type"
          name="type"
          type="text"
          value="simple"
        />
        <input type="submit" value="Start over"/>
      </form>
    </center>
  </body>
</html>
`

const navbar string = `
<form method="get">
  <label for="level">Config level</label>
  <select id="level" name="level">
    <option value="advanced">Advanced</option>
    <option value="simple">Simple</option>
  </select>
  <input type="submit" value="Change"/>
</form>
<hr>
`

const notFound string = `
<html>
  <body>
    <h1>It works!</h1>
  </body>
</html>
`

const simpleDashboard string = `
<html>
  <body>
    <center>` + navbar + `<br>
      <form method="post">
        <table style="width:100%">
          <col style="width:10%">
          <col style="width:90%">

          <tr>
            <td style="text-align:right; width=20">
              <label for="endpoint">Endpoint:</label>
            </td>
            <td>
              <input
                id="endpoint"
                name="endpoint"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="expires">Time (secs):</label>
            </td>
            <td>
              <input
                id="expires"
                name="expires"
                style="width:100%"
                type="text"
              />
            </td>
          </tr>

          <tr>
            <td style="text-align:right; width=20">
              <label for="payload">Payload (hex):</label>
            </td>

            <td>
              <textarea
                id="payload"
                name="payload"
                rows="32"
                style="width:100%"
              ></textarea>
            </td>
          </tr>
        </table>

        <input
          hidden="true"
          id="type"
          name="type"
          type="text"
          value="simple"
        />
        <input type="submit" value="Submit"/>
      </form>
    </center>
  </body>
</html>
`

// This is used for with Fprintf, so must escape %
const simpleResp string = `
<html>
  <body>
    otpURL = %s<br>
    encHex = %s<br>
    <br>
    <center>
      <form method="get">
        <input
          hidden="true"
          id="type"
          name="type"
          type="text"
          value="simple"
        />
        <input type="submit" value="Start over"/>
      </form>
    </center>
  </body>
</html>
`

// Version is the package version
const Version = "1.2.1"
