package responseDataObject

// Web服务器的响应对象
type WebResponseObject struct {
	// 响应结果的状态值
	Code ResultStatus

	// 响应结果的状态值所对应的描述信息
	Message string

	// 响应结果的数据
	Data interface{}
}

func NewWebResponseObject() *WebResponseObject {
	return &WebResponseObject{
		Code:    Success,
		Message: "",
		Data:    nil,
	}
}

func (responseObject *WebResponseObject) SetDataError() {
	responseObject.SetResultStatus(DataError)
}

func (responseObject *WebResponseObject) SetAPIDataError() {
	responseObject.SetResultStatus(APIDataError)
}

func (responseObject *WebResponseObject) SetClientDataError() {
	responseObject.SetResultStatus(APIDataError)
}

func (responseObject *WebResponseObject) SetResultStatus(rs ResultStatus) {
	responseObject.Code = rs
	responseObject.Message = rs.String()
}

func (responseObject *WebResponseObject) SetData(data interface{}) {
	responseObject.Data = data
}
