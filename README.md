根据数据表 Model 来生成 Repository、Service、Handler文件。

## example
```model/user.go
type User struct {
  Name string
  Age int
}
```

> code-generator -m model

代码默认生成在当前路径:
```
internal
   handler
      user_handler.go
   service
      user_service.go
   repository
      user_repository.go
```
