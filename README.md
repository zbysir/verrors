## What is verrors?
verrors是Go1.13官方错误包的辅助库, 目的是让error能附带额外数据, 出于这个宗旨, 它的任何方法都双向兼容errors官方库, 如errors.Is, errors.Unwrap和fmt.Errorf("%w")都能正常使用.

verrors的特性有:
- 辅助库/低入侵: 非强制使用, 兼容errors官方库, 代码改动少
- 支持给error添加任意Value, 
  现在就不再局限于[`errors.WithMessage()`](https://github.com/pkg/errors/blob/master/errors.go#L217)或者 [`errors.WithStack()`](https://github.com/pkg/errors/blob/master/errors.go#L145)
- 灵活, 可插拔, 可扩展.

## Installation

```
go get github.com/zbysir/verrors
```

在国内, 推荐使用go mod与proxy: https://goproxy.cn, 使用更畅快.
## Getting Started

### New error
建议直接使用官方库: 
```
err := errors.New("file not found")
```

如果你想要给这个错误附加上位置或者Code码, 你可以这样:

```
// WithStack会给错误添加位置信息(location information)
// WithCode会给错误添加code信息
err = verrors.WithStack(verrors.WithCode(errors.New("file not found"), 500))

// 打印它
print(verrors.StdPackErrorsFormatter(verrors.Unpack(err)))
```

你会看到两行信息, 和xerror类似, 信息包含了错误链上所有的错误.
```
- file not found [ code = 500; stack = TestX Z:/go/errors/verrors/errors_test.go:100 ]
- file not found
```

你看到代码可能会困扰:
- 为什么New一个带code码的错误这么复杂? 我想要简单点的写法.
- 为什么打印错误的代码这么冗长? 我只想`logf("%+v", err)`.
- 错误链上包含了两个错误, 但我只想要一个包含所有信息的错误应该怎么做?

别急, 这正是由于verrors足够灵活, 在介绍完基本使用方法后我们会解决这些问题.

### Wrapping
使用官方库
```
err := errors.New("file not found") 
err = fmt.Errorf("check health error: %w", err)
```

使用和上面同样的打印代码: `print(verrors.StdPackErrorsFormatter(verrors.Unpack(err)))`, 打印出来看看吧
```
- check health error: file not found
- file not found
```

emmm, 这好像没用上verrors呀, 现在我们就添加上.

```
err := errors.New("file not found") 
err = fmt.Errorf("check health error: %w", verrors.WithCode(verrors.WithStack(err), 400))
```
打印
```
- check health error: file not found [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestReadMe Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/errors_test.go:107 ]
- file not found
```
和官方代码对比, 这段代码仅仅添加verrors.WithCode和verrors.WithStack让错误信息丰富了许多.

那么官方的errors.Is还能使用吗? 答案是肯定的.
```
root := errors.New("file not found") 
err := fmt.Errorf("check health error: %w", verrors.WithCode(verrors.WithStack(root), 400))

print(errors.Is(err, root)) // true
```

不过和上面的问题一样, 这段代码太长了, 再等等, 稍后我们会简化它.

### WithValue
你可能需要为错误添加更多的信息, 如code, stack, need-retry, 这十分简单:

```
err := errors.New("file not found") 
err = fmt.Errorf("check health error: %w", verrors.WithVaule(err, "retry", true))
```
打印
```
- check health error: file not found [ retry = true ]
- file not found
```

实际上`WithCode`也只是WithValue的速记写法.

## Shorthand (简化写法)
收集上面提到的问题, 现在我们来一一解答

- 为什么New一个带code码的错误这么复杂? 我想要简单点的写法.
- 为什么打印错误的代码这么冗长? 我只想`logf("%+v", err)`.
- 错误链上包含了两个错误, 但我只想要一个包含所有信息的错误应该怎么做?

为了简化`fmt.Errorf("do something error: %w", verrors.WithStack(verrors.WithCode(err, 500)))`, 我们提供了Errorfc方法, 先来看它的使用方法:

```
err := errors.New("file not found")
testFileName := "/usr/xx.txt"

err = verrors.Errorfc(400, "do something err: %w, fileName: %s", err, testFileName)
fmt.Printf("\n%+v", err)
```
将会打印如下
```
- do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestErrorfc Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/extra_test.go:18 ]
- file not found
```

现在一次性解决了所有问题, 我们来看看Errorfc如何实现, 在`extra.go`中有它的代码:
```
// Errofc is shorthand for NewFormatError/WithStack/WithCode/fmt.Errorf
func Errorfc(code int,format string, args ...interface{}) (r error) {
	return WithStack(WithCode(NewToInternalError(fmt.Errorf(format, args...)), code), 2)
}
```
它十分简单.

**其调用的所有方法都是导出的, 并且可以随你组合**, 这就是verrors灵活的原因.

但Errorfc方法并不能满足你的需求: 你可能不需要code 或者 stack. 所以它存放在`extra.go`文件中, 表示它仅仅是verror的扩展, 你自己也可以实现这部分的扩展.

实际上所有以`extra`开头的文件都只是verror的扩展(或者说是例子), 你虽然可以使用它们, 但我还是建议自行实现它们以满足项目中个性化的需求.

### 打印
可能需要另起一个段落来说明打印.

verror内置了一个打印扩展: NewFormatError(error) error. 它返回的错误会格式化错误链中所有错误和错误的值(WithValue). 如下
```
- do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestErrorfc Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/extra_test.go:18 ]
- file not found
```

考虑到错误中包含任意Value, 所以就让所有Value都平铺在error信息后面(包括位置信息), 如果你觉得它不够好看, 你可以自定义格式化(format)逻辑:

和NewFormatError实现原理一致, 只需要实现几个方法, 
- Unwrap() 
- InternalError() 
- Error() 
- Format(f fmt.State, c rune)
```
type formatInternalError struct {
	err error
}

func (e formatInternalError) Unwrap() error {
	if u, ok := e.err.(Wrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func (e formatInternalError) InternalError() error {
	return e.err
}

func (e formatInternalError) Error() string {
	return e.err.Error()
}

// 简单的打印错误, 只是为了方便临时查看, 建议用户实现自己的formatInternalError打印方法.
// use %+v to print more info.
func (e formatInternalError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			_, _ = f.Write([]byte(StdPackErrorsFormatter(Unpack(e))))
			return
		}
	}
	_, _ = f.Write([]byte(e.Error()))
}
```

其中`verrors.Unpack(err)` 返回的信息是一个简单的结构体, 它足够简单, 很容易编写格式化代码.

如果懒得写整个formatInternalError结构, 你可以直接编写打印函数来覆盖掉StdPackErrorsFormatter, 下面是一例子
```
verrors.StdPackErrorsFormatter = 
  func (ps PackErrors) string {
    var s strings.Builder
    for _, v := range ps {
        if s.Len() != 0 {
            s.WriteString("\n")
        }
        s.WriteString("- ")
 
        code, codeExist := v.Get("code")
        if codeExist {
            s.WriteString(fmt.Sprintf("[%v] ", code))
        }
 
        s.WriteString(fmt.Sprintf("%v", v.Cause()))
 
        loc, locExist := v.Get("stack")
        if locExist {
            s.WriteString(fmt.Sprintf(" >> %s", loc))
        }
    }
 
    return s.String()
  }
```
现在formatInternalError格式化的错误信息如下
```
- [400] check health error: file not found >> go.zhuzi.me/go/errors/verrors.TestReadMe Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/errors_test.go:109
- file not found
```

另外, 过多的NewFormatError(error) error会增加性能消耗, 所以为了WithStack和WithCode返回的错误也能打印出好看的信息, 我们将它们返回的错误也实现了fmt.Formatter接口, 
如果你想要实现自己的格式化方法, 记得在错误最外层包裹上自己实现了fmt.Formatter的错误, 如下面代码中的`verrors.NewFormatError()`:
```
func Errorfc(code int, format string, args ...interface{}) (r error) {
	return verrors.NewFormatError(verrors.WithStack(verrors.WithCode(verrors.NewToInternalError(fmt.Errorf(format, args...)), code), 2))
}
```

## Proposal (建议)
按照你的喜好来组装verrors.

例如在你的项目中使用Code来标识错误并且喜爱使用位置信息, 你可以在项目中写上下面的工具函数
```
package myerrors

import "github.com/zbysir/verrors"

func NewCode(msg string, code int) error {
	return verrors.WithStack(verrors.WithCode(verrors.NewToInternalError(errors.New(msg)), code))
}

func Errorfc(code int, format string, args ...interface{}) (r error) {
	return verrors.WithStack(verrors.WithCode(verrors.NewToInternalError(fmt.Errorf(format, args...)), code), 2)
}
```

## verrors如何工作?

同样十分简单, 小小的理解下面的接口或者方法:

- Interface: InternalError { InternalError() error }
- Interface: Setter { Set(Store) }
- Func: Unpack(err) PackErrors

Todo ...