package protogen

import (
	"flag"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	flags flag.FlagSet
)

func (file *File) SetLocation(location Location) {
	file.location = location
}

func NewPlugin() *Plugin {
	gen := &Plugin{
		FilesByPath:    make(map[string]*File),
		fileReg:        new(protoregistry.Files),
		enumsByName:    make(map[protoreflect.FullName]*Enum),
		messagesByName: make(map[protoreflect.FullName]*Message),
		opts: Options{
			ParamFunc: flags.Set,
		},
	}

	return gen
}

func (gen *Plugin) GetFileReg() *protoregistry.Files {
	return gen.fileReg
}

func (f *File) ResolveProtoFile(gen *Plugin, desc protoreflect.FileDescriptor) error {
	for i, eds := 0, desc.Enums(); i < eds.Len(); i++ {
		f.Enums = append(f.Enums, newEnum(gen, f, nil, eds.Get(i)))
	}
	for i, mds := 0, desc.Messages(); i < mds.Len(); i++ {
		f.Messages = append(f.Messages, newMessage(gen, f, nil, mds.Get(i)))
	}
	for i, xds := 0, desc.Extensions(); i < xds.Len(); i++ {
		f.Extensions = append(f.Extensions, newField(gen, f, nil, xds.Get(i)))
	}
	for i, sds := 0, desc.Services(); i < sds.Len(); i++ {
		f.Services = append(f.Services, newService(gen, f, sds.Get(i)))
	}
	for _, message := range f.Messages {
		if err := message.resolveDependencies(gen); err != nil {
			return err
		}
	}
	for _, extension := range f.Extensions {
		if err := extension.resolveDependencies(gen); err != nil {
			return err
		}
	}
	for _, service := range f.Services {
		for _, method := range service.Methods {
			if err := method.resolveDependencies(gen); err != nil {
				return err
			}
		}
	}

	return nil
}
