package core

import (
	"reflect"
	"testing"
)

func Test_splitPath(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			"",
			[]string{},
		},
		{
			"/",
			[]string{},
		},
		{
			"////",
			[]string{},
		},
		{
			"aaa/bbb",
			[]string{"aaa", "bbb"},
		},
		{
			"/aaa/bbb",
			[]string{"aaa", "bbb"},
		},
		{
			"/aaa/bbb/",
			[]string{"aaa", "bbb"},
		},
		{
			"//aaa//bbb//",
			[]string{"aaa", "bbb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitPath(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_PathParameters(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			name string
			want []PathParameter
		}{
			{
				"",
				[]PathParameter{},
			},
			{
				"/",
				[]PathParameter{},
			},
			{
				"/aaa",
				[]PathParameter{},
			},
			{
				"/aaa/{id}",
				[]PathParameter{
					{"aaa/{id}", "id"},
				},
			},
			{
				"/aaa/{id}/bbb/{some}",
				[]PathParameter{
					{"aaa/{id}", "id"},
					{"aaa/{id}/bbb/{some}", "some"},
				},
			},
			{
				"///aaa///{id}///",
				[]PathParameter{
					{"aaa/{id}", "id"},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := PathParameters(tt.name)
				if err != nil {
					t.Errorf("PathParameters() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("PathParameters() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []struct {
			name string
		}{
			{
				"/aaa/{id}/bbb/{id}",
			},
			{
				"/aaa/{}",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := PathParameters(tt.name)
				if err == nil {
					t.Error("PathParameters() want error")
					return
				}
			})
		}
	})
}

func Test_duplicatedPathParameters(t *testing.T) {
	tests := []struct {
		name string
		args []PathParameter
		want string
	}{
		{
			"empty",
			[]PathParameter{},
			"",
		},
		{
			"one",
			[]PathParameter{
				{"any string", "id"},
			},
			"",
		},
		{
			"two",
			[]PathParameter{
				{"any string", "id"},
				{"any string", "name"},
			},
			"",
		},
		{
			"three",
			[]PathParameter{
				{"any string", "id"},
				{"any string", "name"},
				{"any string", "id"},
			},
			"id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := duplicatedPathParameters(tt.args); got != tt.want {
				t.Errorf("HasDuplicationPathParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}
