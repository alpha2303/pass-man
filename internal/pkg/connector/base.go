package connector

type IConnector interface {
	Create()
	Read()
	Update()
	Delete()
}
