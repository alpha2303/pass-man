package connector

type FileConnector struct {
	filePath string
}

func CreateFileConnector(filePath string) FileConnector {
	return FileConnector{
		filePath: filePath,
	}
}

func (f FileConnector) Create() {

}

func (f FileConnector) Read() {

}

func (f FileConnector) Update() {

}

func (f FileConnector) Delete() {

}
