package main

import (
	"context"
	"os"
)

func main() {
	ctx := context.Background()
  stdio, err:=os.StartProcess("docker", []string{"system", "stdio"}, &os.ProcAttr{})
	//docker system dial-stdio
  if err != nil {
    print(err)
  } else {
  print(ctx)
  print(stdio)
  }
}
