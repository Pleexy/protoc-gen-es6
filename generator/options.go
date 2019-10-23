package generator

type Options struct {
	UsePrivateFields bool
	GenerateGettersSetters bool
	ValidateOnSet bool
	Flow bool
	ESModules bool
	ConvertString bool
	Grpc bool
}