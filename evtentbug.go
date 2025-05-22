package share

import (
	"reflect"
	"sync"
)

var (
	lock        sync.RWMutex
	subscribers map[EventType][][]any
	once        sync.Once
)

func init() {
	once.Do(func() {
		subscribers = make(map[EventType][][]interface{})
	})
}

// values 的第一个参数是回调函数，后面是回调函数的参数
func EvtBusSubscribe(eventType EventType, values ...any) {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := subscribers[eventType]; !exists {
		subscribers[eventType] = [][]any{}
	}

	callback := reflect.ValueOf(values[0])
	args := make([]reflect.Value, len(values[1:]))
	for i, arg := range values[1:] {
		args[i] = reflect.ValueOf(arg)
	}

	// 将准备好的回调函数和参数存储
	subscribers[eventType] = append(subscribers[eventType], []any{callback, args})
}

// 查找事件回调函数并执行，直接调用预先准备好的反射函数和参数
func EvtBusPublish(event EventType, values ...any) {
	lock.RLock()
	defer lock.RUnlock()

	if events, ok := subscribers[event]; ok {
		for _, event := range events {
			
			callback := event[0].(reflect.Value)
			preparedArgs := event[1].([]reflect.Value)

			// 拼接注册参数和事件参数
			extraArgs := make([]reflect.Value, len(values))
			for i, arg := range values {
				extraArgs[i] = reflect.ValueOf(arg)
			}
			callback.Call(append(preparedArgs, extraArgs...))
		}
	}
}

// func handleLogin(user string, password string, arg1 string, arg2 string) {
// 	println("handleLogin:", user, password, arg1, arg2)
// }

// func handleLogout(user string, password string, arg1 string, arg2 string) {
// 	println("handleLogout:", user, password, arg1, arg2)
// }

// func handleMessage(user string, password string, arg1 string, arg2 string) {
// 	println("handleMessage:", user, password, arg1, arg2)
// }

// func handleAck(user string, password string, arg1 string, arg2 string) {
// 	println("handleAck:", user, password, arg1, arg2)
// }

// // handleFileReq processes a file request with user credentials and arguments.
// // It takes a username, password, and two additional string arguments
// // and prints them to standard output for debugging purposes.
// func handleFileReq(user string, password string, arg1 string, arg2 string) {
// 	println("handleFileReq:", user, password, arg1, arg2)
// }

// // 以上注册是执行帮我写个使用示例，以下是使用示例
// func main() {
// 	// 注册事件处理函数，第一个参数是事件类型，后面是回调函数和回调函数的参数
// 	Subscribe(1, handleLogin, "user1", "password1")
// 	Subscribe(2, handleLogout, "user2", "password2")
// 	Subscribe(3, handleMessage, "user3", "password3")
// 	Subscribe(4, handleAck, "user4", "password4")
// 	Subscribe(5, handleFileReq, "user5", "password5")

// 	// 发布事件，第一个参数是事件类型，后面是事件参数
// 	Publish(1, "test1", "password1")
// 	Publish(2, "test2", "password2")
// 	Publish(3, "test3", "password3")
// 	Publish(4, "test4", "password4")
// 	Publish(5, "test5", "password5")
// }
