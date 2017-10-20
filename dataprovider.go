package main

//
//
//
type IDataProvider interface {
	GetData([]byte) []byte
	GetCSVHeader([]byte) []byte
}

//
//
//
type GetDate func([]byte) []byte

//
//
//
type CounterDataProvider struct {
	testName  []byte
	csvHeader []byte
	ca        *cntArr
}

//
//
//
func NewCounterDataProvider() (*CounterDataProvider, error) {

	return &CounterDataProvider{
		csvHeader: []byte("Testname;D;C;B;A"),
		testName:  []byte("Sawtooth"),
		ca:        &cntArr{},
	}, nil
}

//
// GetCSVHeader
//
func (p *CounterDataProvider) GetCSVHeader(b []byte) []byte {
	return append(b, p.csvHeader...)
}

//
// GetData
//
func (p *CounterDataProvider) GetData(b []byte) []byte {

	// test name
	b = append(b, p.testName...)

	// data
	incCounter(p.ca)
	b = append(b, ';', p.ca.d+48, ';', p.ca.c+48, ';', p.ca.b+48, ';', p.ca.a+48)
	p.ca.a++

	return b

}
