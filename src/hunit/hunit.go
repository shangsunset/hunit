package hunit

import (
  "os"
  "io"
  "fmt"
  "time"
  "bytes"
  "net/http"
  "io/ioutil"
)

import (
  "gopkg.in/yaml.v2"
)

/**
 * HTTP client
 */
var client = http.Client{Timeout: time.Second * 30}

/**
 * A test context
 */
type Context struct {
  BaseURL   string
  Debug     bool
}

/**
 * A test request
 */
type Request struct {
  Method    string                `yaml:"method"`
  URL       string                `yaml:"url"`
  Headers   map[string]string     `yaml:"headers"`
  Entity    string                `yaml:"entity"`
}

/**
 * A test response
 */
type Response struct {
  Status    int                   `yaml:"status"`
  Headers   map[string]string     `yaml:"headers"`
  Entity    string                `yaml:"entity"`
}

/**
 * A test case
 */
type Case struct {
  Request   Request               `yaml:"request"`
  Response  Response              `yaml:"response"`
}

/**
 * Run a test case
 */
func (c Case) Run(context Context) (*Result, error) {
  
  method := c.Request.Method
  if method == "" {
    return nil, fmt.Errorf("Request requires a method (set 'method')")
  }
  
  url := joinPath(context.BaseURL, c.Request.URL)
  if url == "" {
    return nil, fmt.Errorf("Request requires a URL (set 'url')")
  }else if context.Debug {
    fmt.Printf("----> %v %v\n", method, url)
  }
  
  var entity io.Reader
  if c.Request.Entity != "" {
    entity = bytes.NewBuffer([]byte(c.Request.Entity))
  }
  
  req, err := http.NewRequest(method, url, entity)
  if err != nil {
    return nil, err
  }
  
  if c.Request.Headers != nil {
    for k, v := range c.Request.Headers {
      req.Header.Add(k, v)
    }
  }
  
  result := &Result{}
  
  rsp, err := client.Do(req)
  if rsp != nil && rsp.Body != nil {
    defer rsp.Body.Close()
  }
  if err != nil {
    return result.Error(err), nil
  }
  
  result.AssertEqual(c.Response.Status, rsp.StatusCode, "Unexpected status code")
  
  return result, nil
}

/**
 * A test suite
 */
type Suite struct {
  Cases []Case
}

/**
 * Load a test suite
 */
func LoadSuiteFromFile(p string) (*Suite, error) {
  
  file, err := os.Open(p)
  if err != nil {
    return nil, err
  }
  
  data, err := ioutil.ReadAll(file)
  if err != nil {
    return nil, err
  }
  
  return LoadSuite(data)
}

/**
 * Load a test suite
 */
func LoadSuite(data []byte) (*Suite, error) {
  
  var cases []Case
  err := yaml.Unmarshal(data, &cases)
  if err != nil {
    return nil, err
  }
  
  return &Suite{cases}, nil
}

/**
 * Run a test suite
 */
func (s *Suite) Run(context Context) ([]*Result, error) {
  results := make([]*Result, 0)
  
  for _, e := range s.Cases {
    r, err := e.Run(context)
    if err != nil {
      return nil, err
    }
    results = append(results, r)
  }
  
  return results, nil
}

/**
 * Join paths
 */
func joinPath(a, b string) string {
  if a == "" {
    return b
  }
  
  var i, j int
  for i = len(a) - 1; i >= 0; i-- {
    if a[i] != '/' {
      break
    }
  }
  
  for j = 0; j < len(b); j++ {
    if b[j] != '/' {
      break
    }
  }
  
  return a[:i + 1] +"/"+ b[j:]
}
