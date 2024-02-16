package must

func Succeed(err error) {
	panicIfNonNil(err)
}

func Return1[T any](returnValue T, err error) T {
	panicIfNonNil(err)
	return returnValue
}

func Return2[T any, U any](returnValue1 T, returnValue2 U, err error) (T, U) {
	panicIfNonNil(err)
	return returnValue1, returnValue2
}

func panicIfNonNil(err error) {
	if err != nil {
		panic(err)
	}
}
