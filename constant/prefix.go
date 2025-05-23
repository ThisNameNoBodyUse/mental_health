package constant

// 前缀常量

var RegisterPrefix string = "mental:register:"              // 注册前缀
var BlackListPrefix string = "mental:blacklist:"            // Redis黑名单前缀
var UserRolePrefix string = "mental:user_role:"             // 用户角色前缀（用户id对应的角色id集合）
var RolePermissionPrefix string = "mental:role_permission:" // 角色权限前缀（角色id对应的权限id集合）
var APIPermissionPrefix string = "mental:api_permission:"   // 接口权限前缀（权限id对应的可访问的接口路径集合）
