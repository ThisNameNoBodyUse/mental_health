package results

type Result[T any] struct { // Result是个泛型结构体
	Code int    `json:"code"` // 1成功，0失败
	Msg  string `json:"msg"`  // 错误信息
	Data T      `json:"data"` // 数据 是泛型
}

// Success 创建成功响应（无数据）
func Success[T any]() Result[T] {
	return Result[T]{Code: 1}
}

// SuccessWithData 创建成功响应（带数据）
func SuccessWithData[T any](data T) Result[T] {
	return Result[T]{Code: 1, Data: data}
}

// Error 创建错误响应
func Error[T any](msg string) Result[T] {
	return Result[T]{Code: 0, Msg: msg}
}
