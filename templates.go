package main

// repository模板 - 将接口和实现放在同一个文件中
const repositoryTemplate = `
package repository

import (
    "context"
    "daily/internal/repository/model"
    "daily/internal/errcode"
    "gorm.io/gorm"
)

// {{.Name}}Repository 仓储接口
type {{.Name}}Repository interface {
    Create(ctx context.Context, m *model.{{.Name}}) error
    BatchCreate(ctx context.Context, ms []*model.{{.Name}}) error
    Update(ctx context.Context, id int64, m *model.{{.Name}}) error
    UpdateByCondition(ctx context.Context, condition map[string]interface{}, m *model.{{.Name}}) error
    Delete(ctx context.Context, id int64) error
    DeleteByCondition(ctx context.Context, condition map[string]interface{}) error
    Get(ctx context.Context, id int64) (*model.{{.Name}}, error)
    Find(ctx context.Context, condition map[string]interface{}) ([]*model.{{.Name}}, error)
}

// {{.Name | toLower}}Repository 仓储实现
type {{.Name | toLower}}Repository struct {
    db *gorm.DB
}

// New{{.Name}}Repository 创建仓储实例
func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
    return &{{.Name | toLower}}Repository{
        db: db,
    }
}

// Create 创建单条记录
func (r *{{.Name | toLower}}Repository) Create(ctx context.Context, m *model.{{.Name}}) error {
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
        return errcode.Wrap(err, "创建{{.Name}}记录失败")
    }
    return nil
}

// BatchCreate 批量创建记录
func (r *{{.Name | toLower}}Repository) BatchCreate(ctx context.Context, ms []*model.{{.Name}}) error {
    if err := r.db.WithContext(ctx).Create(ms).Error; err != nil {
        return errcode.Wrap(err, "批量创建{{.Name}}记录失败")
    }
    return nil
}

// Update 更新单条记录
func (r *{{.Name | toLower}}Repository) Update(ctx context.Context, id int64, m *model.{{.Name}}) error {
    result := r.db.WithContext(ctx).Model(&model.{{.Name}}{}).Where("id = ?", id).Updates(m)
    if err := result.Error; err != nil {
        return errcode.Wrap(err, "更新{{.Name}}记录失败")
    }
    if result.RowsAffected == 0 {
        return errcode.Wrap(gorm.ErrRecordNotFound, "{{.Name}}记录不存在")
    }
    return nil
}

// UpdateByCondition 条件更新
func (r *{{.Name | toLower}}Repository) UpdateByCondition(ctx context.Context, condition map[string]interface{}, m *model.{{.Name}}) error {
    if err := r.db.WithContext(ctx).Model(&model.{{.Name}}{}).Where(condition).Updates(m).Error; err != nil {
        return errcode.Wrap(err, "条件更新{{.Name}}记录失败")
    }
    return nil
}

// Delete 删除单条记录
func (r *{{.Name | toLower}}Repository) Delete(ctx context.Context, id int64) error {
    result := r.db.WithContext(ctx).Delete(&model.{{.Name}}{}, id)
    if err := result.Error; err != nil {
        return errcode.Wrap(err, "删除{{.Name}}记录失败")
    }
    if result.RowsAffected == 0 {
        return errcode.Wrap(gorm.ErrRecordNotFound, "{{.Name}}记录不存在")
    }
    return nil
}

// DeleteByCondition 条件删除
func (r *{{.Name | toLower}}Repository) DeleteByCondition(ctx context.Context, condition map[string]interface{}) error {
    if err := r.db.WithContext(ctx).Where(condition).Delete(&model.{{.Name}}{}).Error; err != nil {
        return errcode.Wrap(err, "条件删除{{.Name}}记录失败")
    }
    return nil
}

// Get 获取单条记录
func (r *{{.Name | toLower}}Repository) Get(ctx context.Context, id int64) (*model.{{.Name}}, error) {
    var m model.{{.Name}}
    if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errcode.Wrap(err, "{{.Name}}记录不存在")
        }
        return nil, errcode.Wrap(err, "获取{{.Name}}记录失败")
    }
    return &m, nil
}

// Find 查询多条记录
func (r *{{.Name | toLower}}Repository) Find(ctx context.Context, condition map[string]interface{}) ([]*model.{{.Name}}, error) {
    var ms []*model.{{.Name}}
    if err := r.db.WithContext(ctx).Where(condition).Find(&ms).Error; err != nil {
        return nil, errcode.Wrap(err, "查询{{.Name}}记录失败")
    }
    return ms, nil
}
`

// service模板 - 将接口和实现放在同一个文件中
const serviceTemplate = `
package service

import (
    "context"
    "daily/internal/repository"
    "daily/internal/repository/model"
    "daily/internal/errcode"
)

// {{.Name}}Service 服务接口
type {{.Name}}Service interface {
    Create(ctx context.Context, req *{{.Name}}CreateRequest) error
    Update(ctx context.Context, id int64, req *{{.Name}}UpdateRequest) error
    Delete(ctx context.Context, id int64) error
    Get(ctx context.Context, id int64) (*{{.Name}}Response, error)
    List(ctx context.Context, req *{{.Name}}ListRequest) ([]*{{.Name}}Response, error)
}

// 请求响应结构体
type (
    {{.Name}}CreateRequest struct {
        // TODO: 添加创建请求字段
    }
    {{.Name}}UpdateRequest struct {
        // TODO: 添加更新请求字段
    }
    {{.Name}}ListRequest struct {
        // TODO: 添加查询条件字段
    }
    {{.Name}}Response struct {
        // TODO: 添加响应字段
    }    
)

// {{.Name | toLower}}Service 服务实现
type {{.Name | toLower}}Service struct {
    repo repository.{{.Name}}Repository
}

// New{{.Name}}Service 创建服务实例
func New{{.Name}}Service(repo repository.{{.Name}}Repository) {{.Name}}Service {
    return &{{.Name | toLower}}Service{
        repo: repo,
    }
}

// Create 创建
func (s *{{.Name | toLower}}Service) Create(ctx context.Context, req *{{.Name}}CreateRequest) error {
    model := &model.{{.Name}}{
        // TODO: 从请求中复制字段到model
    }
    if err := s.repo.Create(ctx, model); err != nil {
        return errcode.Wrap(err, "创建{{.Name}}失败")
    }
    return nil
}

// Update 更新
func (s *{{.Name | toLower}}Service) Update(ctx context.Context, id int64, req *{{.Name}}UpdateRequest) error {
    model := &model.{{.Name}}{
        // TODO: 从请求中复制字段到model
    }
    if err := s.repo.Update(ctx, id, model); err != nil {
        return errcode.Wrap(err, "更新{{.Name}}失败")
    }
    return nil
}

// Delete 删除
func (s *{{.Name | toLower}}Service) Delete(ctx context.Context, id int64) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return errcode.Wrap(err, "删除{{.Name}}失败")
    }
    return nil
}

// Get 获取详情
func (s *{{.Name | toLower}}Service) Get(ctx context.Context, id int64) (*{{.Name}}Response, error) {
    _, err := s.repo.Get(ctx, id)
    if err != nil {
        return nil, errcode.Wrap(err, "获取{{.Name}}详情失败")
    }

    return &{{.Name}}Response{
        // TODO: 从model复制字段到响应
    }, nil
}

// List 获取列表
func (s *{{.Name | toLower}}Service) List(ctx context.Context, req *{{.Name}}ListRequest) ([]*{{.Name}}Response, error) {
    condition := map[string]interface{}{
        // TODO: 从请求中构造查询条件
    }
    models, err := s.repo.Find(ctx, condition)
    if err != nil {
        return nil, errcode.Wrap(err, "查询{{.Name}}列表失败")
    }

    resp := make([]*{{.Name}}Response, 0, len(models))
    for _, _ = range models {
        resp = append(resp, &{{.Name}}Response{
            // TODO: 从model复制字段到响应
        })
    }
    return resp, nil
}
`

// handler模板
const handlerTemplate = `
package handler

import (
    "daily/internal/response"
    "daily/internal/service"
    "github.com/gin-gonic/gin"
    "strconv"
)

type {{.Name}}Handler struct {
    svc service.{{.Name}}Service
}

func New{{.Name}}Handler(svc service.{{.Name}}Service) *{{.Name}}Handler {
    return &{{.Name}}Handler{
        svc: svc,
    }
}

// Route implements the Router interface
func (h *{{ .Name }}Handler) Route(e *gin.Engine) {
	g := e.Group("/{{ .Name | toLower }}")
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("/:id", h.Get)
	g.GET("", h.List)
}

// Create 创建
func (h *{{.Name}}Handler) Create(c *gin.Context) {
    var req service.{{.Name}}CreateRequest
    if err := c.BindJSON(&req); err != nil {
        response.ParamError(c, "参数错误")
        return
    }

    if err := h.svc.Create(c.Request.Context(), &req); err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil)
}

// Update 更新
func (h *{{.Name}}Handler) Update(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ParamError(c, "ID参数错误")
        return
    }

    var req service.{{.Name}}UpdateRequest
    if err := c.BindJSON(&req); err != nil {
        response.ParamError(c, "参数错误")
        return
    }

    if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil)
}

// Delete 删除
func (h *{{.Name}}Handler) Delete(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ParamError(c, "ID参数错误")
        return
    }

    if err := h.svc.Delete(c.Request.Context(), id); err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil)
}

// Get 获取详情
func (h *{{.Name}}Handler) Get(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.ParamError(c, "ID参数错误")
        return
    }

    data, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, data)
}

// List 获取列表
func (h *{{.Name}}Handler) List(c *gin.Context) {
    var req service.{{.Name}}ListRequest
    if err := c.BindQuery(&req); err != nil {
        response.ParamError(c, "参数错误")
        return
    }
    data, err := h.svc.List(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, err)
        return
    }
    response.Success(c, data)
}
`
