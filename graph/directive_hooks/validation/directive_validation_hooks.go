package directiveValidationHooks

import (
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

var directives = map[string]func(directive *ast.Directive, f *modelgen.Field){
	"maxLength":      maxLengthDirective,
	"minLength":      minLengthDirective,
	"maxValue":       maxValueDirective,
	"minValue":       minValueDirective,
	"notEmptyString": notEmptyStringDirective,
}

func CustomFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}

	for directiveName, directiveFunction := range directives {
		c := fd.Directives.ForName(directiveName)

		if c != nil {
			directiveFunction(c, f)

		}
	}

	return f, nil
}

func maxLengthDirective(directive *ast.Directive, f *modelgen.Field) {
	formatConstraint := directive.Arguments.ForName("value")

	if formatConstraint != nil {
		f.Tag += " maxLength:" + formatConstraint.Value.String()
	}
}

func minLengthDirective(directive *ast.Directive, f *modelgen.Field) {
	formatConstraint := directive.Arguments.ForName("value")

	if formatConstraint != nil {
		f.Tag += " minLength:" + formatConstraint.Value.String()
	}
}

func maxValueDirective(directive *ast.Directive, f *modelgen.Field) {
	formatConstraint := directive.Arguments.ForName("value")

	if formatConstraint != nil {
		f.Tag += " maxValue:" + formatConstraint.Value.String()
	}
}

func minValueDirective(directive *ast.Directive, f *modelgen.Field) {
	formatConstraint := directive.Arguments.ForName("value")

	if formatConstraint != nil {
		f.Tag += " minValue:" + formatConstraint.Value.String()
	}
}

func notEmptyStringDirective(directive *ast.Directive, f *modelgen.Field) {
	formatConstraint := directive.Arguments.ForName("state")

	if formatConstraint != nil {
		f.Tag += " notEmptyString:true"
	}
}
