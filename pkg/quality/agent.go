package quality

type QualityAgent interface {
	Pass(input interface{}) interface{}
}
