package playground

import (
  "fmt"
  "github.com/google/go-cmp/cmp"
  "reflect"
  "slices"
  "strings"
  "testing"
)

var collfile = strings.NewReader(`{"info":{"_postman_id":"c5d10f58-0959-4eba-901b-6d10088ec094","name":"TEST_COLL","schema":"https://schema.getpostman.com/json/collection/v2.1.0/collection.json","_exporter_id":"25152555"},"item":[{"name":"Posts","item":[{"name":"Post Comments","item":[{"name":"Get all post comments (v1)","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/posts/{{post_id}}/comments","host":["{{host}}"],"path":["posts","{{post_id}}","comments"]}},"response":[]},{"name":"Get all post comments (v2)","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/comments?postId={{post_id}}","host":["{{host}}"],"path":["comments"],"query":[{"key":"postId","value":"{{post_id}}"}]}},"response":[]}]},{"name":"Get all posts","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/posts","host":["{{host}}"],"path":["posts"]}},"response":[]},{"name":"Get one post","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/posts/{{post_id}}","host":["{{host}}"],"path":["posts","{{post_id}}"]}},"response":[]},{"name":"Create post","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"POST","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"},{"key":"Content-Type","value":"application/json","type":"text"}],"body":{"mode":"raw","raw":"{\n\t\"title\": \"sunt aut facere repellat\",\n\t\"body\": \"quia et suscipit quas totam\"\n}\n","options":{"raw":{"language":"json"}}},"url":{"raw":"{{host}}/posts","host":["{{host}}"],"path":["posts"]}},"response":[]}]},{"name":"Users","item":[{"name":"Get all users","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/users?page=1&limit=100","host":["{{host}}"],"path":["users"],"query":[{"key":"page","value":"1"},{"key":"limit","value":"100"}]}},"response":[]},{"name":"Get one user","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/users/{{user_id}}","host":["{{host}}"],"path":["users","{{user_id}}"]}},"response":[]}]},{"name":"Home","request":{"method":"GET","header":[],"url":{"raw":"{{host}}/","host":["{{host}}"],"path":[""]}},"response":[]},{"name":"Get all photos","protocolProfileBehavior":{"disabledSystemHeaders":{"connection":true,"accept-encoding":true,"accept":true}},"request":{"method":"GET","header":[{"key":"Connection","value":"keep-alive","type":"text"},{"key":"Accept-Encoding","value":"gzip, deflate","type":"text"},{"key":"Accept","value":"*/*","type":"text"}],"url":{"raw":"{{host}}/photos","host":["{{host}}"],"path":["photos"]}},"response":[]}],"event":[{"listen":"prerequest","script":{"type":"text/javascript","packages":{},"exec":[""]}},{"listen":"test","script":{"type":"text/javascript","packages":{},"exec":[""]}}],"variable":[{"key":"host","value":"https://jsonplaceholder.typicode.com","type":"string"},{"key":"post_id","value":"5","type":"string"},{"key":"user_id","value":"10","type":"string"}]}`)

func Test_parseColl(t *testing.T) {
  got, _ := parseColl(collfile)

  want := &coll{
    Info: struct {
      Name string "json:\"name\""
    }{Name: "TEST_COLL"},
    Item: []collItem{
      {
        Name: "Posts",
        Item: []collItem{
          {
            Name: "Post Comments",
            Item: []collItem{
              {
                Name: "Get all post comments (v1)",
                Item: []collItem(nil),
                Request: &collRequest{
                  URL: collURL{
                    Raw:      "{{host}}/posts/{{post_id}}/comments",
                    Protocol: "",
                    Host:     []string{"{{host}}"},
                    Path:     []string{"posts", "{{post_id}}", "comments"},
                    Port:     "",
                    Query:    []collQueryParam(nil),
                  },
                  Method: "GET",
                  Header: []collHeader{
                    {Key: "Connection", Value: "keep-alive", Disabled: false},
                    {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                    {Key: "Accept", Value: "*/*", Disabled: false},
                  },
                },
              },
              {
                Name: "Get all post comments (v2)",
                Item: []collItem(nil),
                Request: &collRequest{
                  URL: collURL{
                    Raw:      "{{host}}/comments?postId={{post_id}}",
                    Protocol: "",
                    Host:     []string{"{{host}}"},
                    Path:     []string{"comments"},
                    Port:     "",
                    Query:    []collQueryParam{{Key: "postId", Value: "{{post_id}}"}},
                  },
                  Method: "GET",
                  Header: []collHeader{
                    {Key: "Connection", Value: "keep-alive", Disabled: false},
                    {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                    {Key: "Accept", Value: "*/*", Disabled: false},
                  },
                },
              },
            },
          },
          {
            Name: "Get all posts",
            Item: []collItem(nil),
            Request: &collRequest{
              URL: collURL{
                Raw:      "{{host}}/posts",
                Protocol: "",
                Host:     []string{"{{host}}"},
                Path:     []string{"posts"},
                Port:     "",
                Query:    []collQueryParam(nil)},
              Method: "GET",
              Header: []collHeader{
                {Key: "Connection", Value: "keep-alive", Disabled: false},
                {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                {Key: "Accept", Value: "*/*", Disabled: false},
              },
            },
          },
          {
            Name: "Get one post",
            Item: []collItem(nil),
            Request: &collRequest{
              URL: collURL{
                Raw:      "{{host}}/posts/{{post_id}}",
                Protocol: "",
                Host:     []string{"{{host}}"},
                Path:     []string{"posts", "{{post_id}}"},
                Port:     "",
                Query:    []collQueryParam(nil)},
              Method: "GET",
              Header: []collHeader{
                {Key: "Connection", Value: "keep-alive", Disabled: false},
                {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                {Key: "Accept", Value: "*/*", Disabled: false},
              },
            },
          },
          {
            Name: "Create post",
            Item: []collItem(nil),
            Request: &collRequest{
              URL: collURL{
                Raw:      "{{host}}/posts",
                Protocol: "",
                Host:     []string{"{{host}}"},
                Path:     []string{"posts"},
                Port:     "",
                Query:    []collQueryParam(nil)},
              Method: "POST",
              Header: []collHeader{
                {Key: "Connection", Value: "keep-alive", Disabled: false},
                {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                {Key: "Accept", Value: "*/*", Disabled: false},
                {Key: "Content-Type", Value: "application/json", Disabled: false},
              },
              Body: &collBody{
                Mode:       "raw",
                Raw:        "{\n\t\"title\": \"sunt aut facere repellat\",\n\t\"body\": \"quia et suscipit quas totam\"\n}\n",
                URLEncoded: []collURLEncodedParameter(nil),
              },
            },
          },
        },
      },
      {
        Name: "Users",
        Item: []collItem{
          {
            Name: "Get all users",
            Item: []collItem(nil),
            Request: &collRequest{
              URL: collURL{
                Raw:      "{{host}}/users?page=1&limit=100",
                Protocol: "",
                Host:     []string{"{{host}}"},
                Path:     []string{"users"},
                Port:     "",
                Query: []collQueryParam{
                  {Key: "page", Value: "1"},
                  {Key: "limit", Value: "100"},
                },
              },
              Method: "GET",
              Header: []collHeader{
                {Key: "Connection", Value: "keep-alive", Disabled: false},
                {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                {Key: "Accept", Value: "*/*", Disabled: false},
              },
            },
          },
          {
            Name: "Get one user",
            Item: []collItem(nil),
            Request: &collRequest{
              URL: collURL{
                Raw:      "{{host}}/users/{{user_id}}",
                Protocol: "",
                Host:     []string{"{{host}}"},
                Path:     []string{"users", "{{user_id}}"},
                Port:     "",
                Query:    []collQueryParam(nil)},
              Method: "GET",
              Header: []collHeader{
                {Key: "Connection", Value: "keep-alive", Disabled: false},
                {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
                {Key: "Accept", Value: "*/*", Disabled: false},
              },
            },
          },
        },
      },
      {
        Name: "Home",
        Item: []collItem(nil),
        Request: &collRequest{
          URL: collURL{
            Raw:      "{{host}}/",
            Protocol: "",
            Host:     []string{"{{host}}"},
            Path:     []string{""},
            Port:     "",
            Query:    []collQueryParam(nil)},
          Method: "GET",
          Header: []collHeader{},
        },
      },
      {Name: "Get all photos",
        Item: []collItem(nil),
        Request: &collRequest{
          URL: collURL{
            Raw:      "{{host}}/photos",
            Protocol: "",
            Host:     []string{"{{host}}"},
            Path:     []string{"photos"},
            Port:     "",
            Query:    []collQueryParam(nil)},
          Method: "GET",
          Header: []collHeader{
            {Key: "Connection", Value: "keep-alive", Disabled: false},
            {Key: "Accept-Encoding", Value: "gzip, deflate", Disabled: false},
            {Key: "Accept", Value: "*/*", Disabled: false},
          },
        },
      },
    },
    Variable: []collVariable{
      {Key: "host", Value: "https://jsonplaceholder.typicode.com", Type: "string", Name: "", Disabled: false},
      {Key: "post_id", Value: "5", Type: "string", Name: "", Disabled: false},
      {Key: "user_id", Value: "10", Type: "string", Name: "", Disabled: false},
    },
  }

  if !reflect.DeepEqual(want, got) {
    t.Fatal(cmp.Diff(want, got))
  }
}

func Test_walk(t *testing.T) {
  collfile.Seek(0, 0)
  c, err := parseColl(collfile)

  if nil != err {
    fmt.Println(c)
    t.Fatalf("unexpected error: %s", err)
  }

  var (
    requestsArrayBuilder = &strings.Builder{}
    dirtreeBuilder       = &strings.Builder{}
    variables            = make(map[string]collVariable)
  )

  for v := range slices.Values(c.Variable) {
    variables[v.Key] = v
  }

  newString = func() string { return "714b9856-cac2-4a77-a149-ca1a797918cb" }
  walk(variables, requestsArrayBuilder, dirtreeBuilder, "", c.Item)

  got := requestsArrayBuilder.String()
  want := `{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all post comments (v1)","full_name":"Posts / Post Comments / Get all post comments (v1)","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts/{{post_id}}/comments","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts/5/comments",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all post comments (v2)","full_name":"Posts / Post Comments / Get all post comments (v2)","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/comments?postId={{post_id}}","url_port":"","url_protocol":"","url_query":[{"key":"postId","value":"{{post_id}}",},],"url_resolved":"https://jsonplaceholder.typicode.com/comments?postId=5",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all posts","full_name":"Posts / Get all posts","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get one post","full_name":"Posts / Get one post","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts/{{post_id}}","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts/5",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Create post","full_name":"Posts / Create post","request_method":"POST","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},{"key":"Content-Type","value":"application/json",},],"request_body_mode":"raw","request_body_raw":"{\n\t\"title\": \"sunt aut facere repellat\",\n\t\"body\": \"quia et suscipit quas totam\"\n}\n","url_raw":"{{host}}/posts","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all users","full_name":"Users / Get all users","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/users?page=1&limit=100","url_port":"","url_protocol":"","url_query":[{"key":"page","value":"1",},{"key":"limit","value":"100",},],"url_resolved":"https://jsonplaceholder.typicode.com/users?page=1&limit=100",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get one user","full_name":"Users / Get one user","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/users/{{user_id}}","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/users/10",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Home","full_name":"Home","request_method":"GET","request_header":[],"url_raw":"{{host}}/","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all photos","full_name":"Get all photos","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/photos","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/photos",},`

  if !reflect.DeepEqual(want, got) {
    t.Fatal(cmp.Diff(want, got))
  }
}

func Test_collGen(t *testing.T) {
  collfile.Seek(0, 0)
  newString = func() string { return "714b9856-cac2-4a77-a149-ca1a797918cb" }
  collsrc, colldirtree, _ := collGen(collfile)

  want := `<script>const requests = [{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all post comments (v1)","full_name":"Posts / Post Comments / Get all post comments (v1)","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts/{{post_id}}/comments","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts/5/comments",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all post comments (v2)","full_name":"Posts / Post Comments / Get all post comments (v2)","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/comments?postId={{post_id}}","url_port":"","url_protocol":"","url_query":[{"key":"postId","value":"{{post_id}}",},],"url_resolved":"https://jsonplaceholder.typicode.com/comments?postId=5",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all posts","full_name":"Posts / Get all posts","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get one post","full_name":"Posts / Get one post","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/posts/{{post_id}}","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts/5",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Create post","full_name":"Posts / Create post","request_method":"POST","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},{"key":"Content-Type","value":"application/json",},],"request_body_mode":"raw","request_body_raw":"{\n\t\"title\": \"sunt aut facere repellat\",\n\t\"body\": \"quia et suscipit quas totam\"\n}\n","url_raw":"{{host}}/posts","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/posts",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all users","full_name":"Users / Get all users","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/users?page=1&limit=100","url_port":"","url_protocol":"","url_query":[{"key":"page","value":"1",},{"key":"limit","value":"100",},],"url_resolved":"https://jsonplaceholder.typicode.com/users?page=1&limit=100",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get one user","full_name":"Users / Get one user","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/users/{{user_id}}","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/users/10",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Home","full_name":"Home","request_method":"GET","request_header":[],"url_raw":"{{host}}/","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/",},{"id":"714b9856-cac2-4a77-a149-ca1a797918cb","name":"Get all photos","full_name":"Get all photos","request_method":"GET","request_header":[{"key":"Connection","value":"keep-alive",},{"key":"Accept-Encoding","value":"gzip, deflate",},{"key":"Accept","value":"*/*",},],"url_raw":"{{host}}/photos","url_port":"","url_protocol":"","url_query":[],"url_resolved":"https://jsonplaceholder.typicode.com/photos",},];</script>`

  if !reflect.DeepEqual(want, collsrc) {
    t.Fatal(cmp.Diff(want, collsrc))
  }

  want = `<header><h3>TEST_COLL</h3></header><div class="item folder"><span class="name">Posts</span><div class="item folder"><span class="name">Post Comments</span><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get all post comments (v1)</span></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get all post comments (v2)</span></div></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get all posts</span></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get one post</span></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Create post</span></div></div><div class="item folder"><span class="name">Users</span><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get all users</span></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get one user</span></div></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Home</span></div><div class="item"><span data-id="714b9856-cac2-4a77-a149-ca1a797918cb" class="name">Get all photos</span></div>`

  if !reflect.DeepEqual(want, colldirtree) {
    t.Fatal(cmp.Diff(want, colldirtree))
  }
}
