## What is verrors?
verrors是Go1.13官方错误包的辅助库, 目的是让error能附带额外数据, 出于这个宗旨, 它的任何方法都双向兼容errors官方库, 如errors.Is, errors.Unwrap和fmt.Errorf("%w")都能正常使用.

verrors的特性有:
- 辅助库/低入侵: 非强制使用, 兼容errors官方库
- 支持给error添加任意Value, 
  现在就不再局限于[`errors.WithMessage()`](https://github.com/pkg/errors/blob/master/errors.go#L217)或者 [`errors.WithStack()`](https://github.com/pkg/errors/blob/master/errors.go#L145)
- 灵活, 可插拔, 可扩展.

欢迎尝鲜和讨论.
## Installation

```
go get github.com/zbysir/verrors
```

在国内, 推荐使用go mod与proxy: https://goproxy.cn, 使用更畅快.
## Getting Started

### New error
你应该直接使用官方库: 
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

和xerror类似, 信息包含了错误链上所有的错误, 由于错误链上只有一个错误, 所以你会看到一行信息:
```
- file not found [ code = 500; stack = TestX Z:/go/errors/verrors/errors_test.go:100 ]
```

你看到代码可能会困扰:
- 为什么New一个带code码的错误这么复杂? 我想要简单点的写法.
- 为什么打印错误的代码这么冗长? 我只想`logf("%+v", err)`.

别急, 这正是由于verrors足够灵活, 在介绍完基本使用方法后我们会解决这些问题.

### Wrapping
同样使用官方库
```
err := errors.New("file not found") 
err = fmt.Errorf("check health error: %w", err)
```

使用和上面同样的打印代码: `print(verrors.StdPackErrorsFormatter(verrors.Unpack(err)))`, 打印出来看看吧
```
- check health error: file not found
- file not found
```

emmm, 这好像没用上verrors呀, 现在我们就添加上, 和上面New Error一样, 需要给官方错误包裹上信息:
```
err := errors.New("file not found") 
err = WithStack(WithCode(fmt.Errorf("check health error: %w", err), 500))
```
打印
```
- check health error: file not found [ code = 500; stack = github.com/zbysir/verrors.TestReadMe Z:/go_project/verrors/errors_test.go:71 ]
- file not found
```
和官方代码对比, 这段代码仅仅添加verrors.WithCode和verrors.WithStack让错误信息丰富了许多.

那么官方的errors.Is还能使用吗? 答案是肯定的.
```
root := WithCode(errors.New("file not found"), 400)
err := WithStack(fmt.Errorf("check health error: %w", root))

print(errors.Is(err, root)) // true
```
errors.Unwrap行为也一致
```
root := WithCode(WithStack(errors.New("file not found")), 400)
err := WithStack(fmt.Errorf("check health error: %w", root))

print(errors.Unwrap(err) == root) // true
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
**集中精神, 重点来了**, 这里将会说明verrors如何易扩展.

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
这便是verrors的最终形态, 一次性解决了所有问题, 我们再来看看Errorfc如何实现, 在`extra.go`中有它的代码:
```
// Errorfc is shorthand for WithStack/WithCode/fmt.Errorf
func Errorfc(code int,format string, args ...interface{}) (r error) {
	return WithStack(WithCode(fmt.Errorf(format, args...), code), 2)
}
```
相反的, 它十分简单.

**其调用的所有方法都是你可以自行实现的, 并且可以随意组合**(如你不想要stack, 删除掉WithStack即可), 这就是verrors灵活可扩展的原因.

也许Errorfc方法并不能满足你的需求: 你可能不需要code 或者 stack, 所以它存放在`extra.go`文件中, 表示它仅仅是verror的扩展, 
实际上所有以`extra`开头的文件都只是verror内置的扩展(或者说是例子), 这意味着所有`extra`中的所有功能(包括WithCode, WithStack)都可以由你自己实现, 
至于如何使用verror内置的扩展, 随你喜好.

**如果你要自行实现WithXXX, 最简单的使用办法是copy`extra`中的代码到你的项目中, 并修改它们.**

### Print (打印)
如何打印实则和错误无关, 所以我们提供Unpack方法, 它可以将错误链中的信息格式化成为规整的结构体, 方便你自行实现打印.

verror内置了一个打印扩展: WithFormat(error) error. 它返回的错误会格式化错误链中所有错误和错误的值(WithValue). 如下
```
- do something err: file not found, fileName: /usr/xx.txt [ code = 400; stack = go.zhuzi.me/go/errors/verrors.TestErrorfc Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/extra_test.go:18 ]
- file not found
```

考虑到错误中包含任意Value, 所以就让所有Value都平铺在error信息后面(包括位置信息), 如果你觉得它不够好看, 你可以自定义格式化(format)逻辑:

和verrors中所有的InternalError实现原理一致, 只需要实现几个方法, 
- Unwrap()
- InternalError()
- Error()
- Format(f fmt.State, c rune)
```
type formatInternalError struct {
	err error
}

func (e formatInternalError) Unwrap() error {
	return Unwrap(e.err)
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

其中`verrors.Unpack(err)` 返回一个简单的结构体, 它足够整洁简单, 很容易编写格式化代码.

如果懒得写整个formatInternalError结构, 你还可以直接编写打印函数来覆盖掉默认的格式化函数`StdPackErrorsFormatter`, 下面是一例子
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
现在 格式化的错误信息如下
```
- [400] check health error: file not found >> go.zhuzi.me/go/errors/verrors.TestReadMe Z:/golang/go_path/src/go.zhuzi.me/go/errors/verrors/errors_test.go:109
- file not found
```
如你所见, code码被放在了前面, 位置信息也更好查看, 你可以根据项目需求调整它.

另外, 过多的WithFormat(error)会增加性能消耗, 所以为了WithStack和WithCode返回的错误也能打印出好看的信息, 我们将它们返回的错误也实现了fmt.Formatter接口, 而不是每次都WithFormat.

总结一下, 如果你想要实现自己的格式化方法, 有两种办法
- 在错误最外层包裹上自己实现了fmt.Formatter的错误, 如下面代码中的`verrors.WithFormat()`:
```
func Errorfc(code int, format string, args ...interface{}) (r error) {
	return verrors.WithFormat(verrors.WithStack(verrors.WithCode(fmt.Errorf(format, args...), code), 2))
}
```
- 替换掉verrors.StdPackErrorsFormatter函数.

推荐使用第二个办法.

> 当然, 我们现在提及的全部方法都只是verrors的"扩展", 如果需要, 你也可以全部自行实现他们.

## Proposal (建议)
虽然verrors可以开箱即用:
```
err := errors.New("file not found") 
err = verrors.WithVaule(fmt.Errorf("check health error: %w", err), "retry", true))
```
但这样的写法太冗长, 不过我们没办法提供简单写法, 因为你需要的错误信息是可定制的. 
所以使用verrors最好的办法是按照你的喜好在项目中组装verrors.

例如在你的项目中使用Code来标识错误并且喜爱使用位置信息, 你可以在项目中写下面的工具函数.
```
package myerrors

import "github.com/zbysir/verrors"

func NewCode(msg string, code int) error {
	return verrors.WithStack(verrors.WithCode(errors.New(msg), code))
}

func Errorfc(code int, format string, args ...interface{}) (r error) {
	return verrors.WithStack(verrors.WithCode(fmt.Errorf(format, args...), code), 2)
}
```

使用它
```
import "project/myerrors"

uid := 1
err := mysql.GetUser(uid)
err = myerrors.Errorfc(500, "GetUser error:%w, id: %v", err, uid)

log.Printf("%+v", err)
```

## How verrors work?
刚刚一直在说, 建议用户自行实现某某方法, 全都被用户实现了, 那verrors到底为我们提供了什么?

实际上verrors只提供了Unpack方法和它的思路, 这部分逻辑我们不希望用户去实现, 而是应该直接使用.

Unpack会解包一个错误, 和Unwrap不一样的是它可以在错误链中插入内部错误(InternalError)但不影响层级, 自定义层级与数据逻辑通过以下两个接口实现
- InternalError { InternalError() error }
- Setter { Set(Store) }

实现了 InternalError 的错误是一个内部错误, 它不会被放置到错误链中, 而是作为数据存储或者格式化时使用. 

实现了 Setter 的错误 会被当做数据层来实现WithValue.

如果你要自定义一个InternalError, 最好的办法就是参考`extra_stack_error.go`和`extra_value_error.go`.

## Obscure
此包还需要解决的问题:
- 过多的Wrap造成的性能消耗
- error本身应该是简单的, 引入自定义的(稍显复杂的)verrors是否本末倒置? 是否需要将verrors做得傻瓜化一些?
- 兼容官方库是否有意义? 有什么场景会同时使用两个库(即verrors和官方errors)?

作者也在思考上面的问题, 慢慢改进吧.
