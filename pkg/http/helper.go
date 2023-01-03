package http

func Ok(body any) *HttpResponse {
	return &HttpResponse{
		StatusCode: 200,
		Body:       body,
	}
}

func BadRequest(e error) *HttpResponse {
	return &HttpResponse{
		StatusCode: 400,
		Errors:     append(make([]error, 0), e),
	}
}

func BadValidation(e []error) *HttpResponse {
	return &HttpResponse{
		StatusCode: 400,
		Errors:     e,
	}
}

//export const noContent = (): HttpResponse => ({
//statusCode: 204,
//});
//
//export const badRequest = (error: Error): HttpResponse<Error> => ({
//statusCode: 400,
//body: error,
//});
//
//export const unauthorized = (error: Error): HttpResponse<Error> => ({
//statusCode: 401,
//body: error,
//});
//
//export const forbidden = (error: Error): HttpResponse<Error> => ({
//statusCode: 403,
//body: error,
//});
//
//export const notFound = (error: Error): HttpResponse<Error> => ({
//statusCode: 404,
//body: error,
//});
//
//export const serverError = (error?: Error | unknown): HttpResponse<Error> => {
//const stack = error instanceof Error ? error.stack : undefined;
//return {
//statusCode: 500,
//body: new ServerError(stack),
//};
//};
