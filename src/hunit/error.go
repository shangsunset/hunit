package hunit

import (
  "reflect"
)

import (
  "github.com/kr/pretty"
  "github.com/davecgh/go-spew/spew"
)

/**
 * An assertion error
 */
type AssertionError struct {
  Expected    interface{}
  Actual      interface{}
  Message     string
}

/**
 * Error
 */
func (e AssertionError) Error() string {
  
  m := e.Message
  if e.Message != "" {
    m += ":\n"
  }
  
  _, ek := typeAndKind(e.Expected)
  if ek == reflect.String || ek == reflect.Struct || ek == reflect.Map || ek == reflect.Slice || ek == reflect.Array {
    d := pretty.Diff(e.Expected, e.Actual)
    if d != nil && len(d) > 0 {
      for _, e := range d {
        m += e +"\n"
      }
    }
  }else{
    m += "expected: "+ spew.Sdump(e.Expected)
    m += "  actual: "+ spew.Sdump(e.Actual)
  }
  
  return m
}

/**
 * Obtain a value's type and kind
 */
func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
  t := reflect.TypeOf(v)
  k := t.Kind()
  for k == reflect.Ptr {
    t = t.Elem()
    k = t.Kind()
  }
  return t, k
}
