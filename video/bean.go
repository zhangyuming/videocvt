package video

type mediaConvert interface {

	Convert(source string, dist string)(result string,err error);
	Teardown(key string);
	Reset(source string)(result string,err error);

}


