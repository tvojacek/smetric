package smetric

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type NameBuilder struct {
	parameters  string
	fullName    string
	serviceName string
	lastSuffix  string
	errors      []error
}

var globalBuilder NameBuilder
var globalBuilderLock sync.RWMutex

func GetGlobalBuilder() NameBuilder {
	globalBuilderLock.RLock()
	defer globalBuilderLock.RUnlock()
	return globalBuilder
}

func NewNameBuilder(prefix string) NameBuilder {
	b := NameBuilder{
		fullName:    prefix,
		serviceName: prefix,
	}
	return b
}

func (b NameBuilder) LastSuffix() string {
	return b.lastSuffix
}

func (b NameBuilder) WithSuffix(name string) NameBuilder {
	if name == "" {
		b.errors = append(b.errors, fmt.Errorf("ignoring empty suffix after: %s", b.fullName))
		return b
	}
	if b.fullName != "" {
		b.fullName += "_" + name
	} else {
		b.fullName = name
	}
	b.lastSuffix = name
	return b
}

var EmptyParameterValue = strconv.Quote("1")

func (b NameBuilder) WithParameter(name string, value string) NameBuilder {
	if name == "" {
		b.errors = append(b.errors, fmt.Errorf("ignoring parameter with empty suffix after: %s", b.fullName))
		return b
	}
	if value == "" {
		value = EmptyParameterValue
	}
	if b.parameters != "" {
		b.parameters = b.parameters + "," + name + "=" + strconv.Quote(value)
	} else {
		b.parameters = name + "=" + strconv.Quote(value)
	}
	return b
}

func (b NameBuilder) Error() error {
	if b.errors != nil {
		return errors.Join(b.errors...)
	} else {
		return nil
	}
}

func (b NameBuilder) nameWithParameters() string {
	if len(b.parameters) == 0 {
		return b.fullName
	} else {
		return b.fullName + "{" + b.parameters + "}"
	}
}

func (b NameBuilder) String() (string, error) {
	return b.nameWithParameters(), b.Error()
}
