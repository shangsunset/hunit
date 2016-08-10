package hunit

import (
  "fmt"
  "github.com/bww/epl"
)

/**
 * Interpolate expressions in a string
 */
func interpolate(s, pre, suf string, context interface{}) (string, error) {
  if len(pre) < 1 || len(suf) < 1 {
    return "", fmt.Errorf("Invalid variable prefix/suffix")
  }
  
  fp := pre[0]
  fs := suf[0]
  
  var out string
  var i, esc int
  for {
    if i >= len(s) {
      break
    }
    
    if s[i] == '\\' {
      esc++
      if (esc % 2) == 0 {
        out += "\\"
      }
      i++
      continue
    }else{
      esc = 0
    }
    
    if s[i] == fp && matchAhead(s[i:], pre) {
      i += len(pre)
      start := i
      for {
        if i >= len(s) {
          return "", fmt.Errorf("Unexpected end-of-input")
        }
        if s[i] == fs && matchAhead(s[i:], suf) {
          
          prg, err := epl.Compile(s[start:i])
          if err != nil {
            return "", err
          }
          
          res, err := prg.Exec(context)
          if err != nil {
            return "", err
          }
          
          switch v := res.(type) {
            case string:
              out += v
            default:
              out += fmt.Sprintf("%v", v)
          }
          
          i += len(suf)
          break
        }else{
          i++
        }
      }
    }else{
      out += string(s[i])
      i++
    }
    
  }
  
  return out, nil
}

/**
 * Match ahead
 */
func matchAhead(s, x string) bool {
  if len(s) < len(x) {
    return false
  }
  for i := 0; i < len(x); i++ {
    if s[i] != x[i] {
      return false
    }
  }
  return true
}