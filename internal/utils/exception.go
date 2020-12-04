package utils

type Exception struct {
	Id  int    // exception id
	Msg string // exception msg
}
type TryStruct struct {
	catches map[int]ExceptionHandler
	try     func()
}

func Try(tryHandler func()) *TryStruct {
	tryStruct := TryStruct{
		catches: make(map[int]ExceptionHandler),
		try:     tryHandler,
	}
	return &tryStruct
}

type ExceptionHandler func(Exception)

func (this *TryStruct) Catch(exceptionId int, catch func(Exception)) *TryStruct {
	this.catches[exceptionId] = catch
	return this
}

func (this *TryStruct) Finally(finally func()) {
	defer func() {
		if e := recover(); nil != e {
			exception := e.(Exception)
			if catch, ok := this.catches[exception.Id]; ok {
				catch(exception)
			}
			finally()
		}
	}()

	this.try()
}

func Throw(id int, msg string) Exception {
	panic(Exception{id, msg})
}

