package ecpg

import (
	"testing"
)

func Test_checkDependencies(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "command exists",
			args: args{
				p: "ls",
			},
			want: true,
		},
		{
			name: "command does not exists",
			args: args{
				p: "someCommandThatDoesntExist",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDependencies(tt.args.p); got != tt.want {
				t.Errorf("checkDependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_executedECPGCommand tests a couple simple failures
func Test_executeECPGCommand(t *testing.T) {
	type args struct {
		stmt string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succesful query",
			args: args{
				stmt: `EXEC SQL SELECT COUNT(ID) FROM "SOME_SCHEMA"."SOME_TABLE";`,
			},
			wantErr: false,
		},
		{
			name: "bad query",
			args: args{
				stmt: `EXEC SQL SELECT COUNT(*) FROM "SOME"."TABLES" WHERE "SCHEMA_NAME" = ? AND "TABLE_NAME" = ? AND "IS_AWESOME" IN ('TRUE','FALSE');`,
			},
			wantErr: true,
		},
		{
			name: "bad query",
			args: args{
				stmt: `EXEC SQL SELECT COUNT(*) FROM "SOME"."TABLES" WHERE "SCHEMA_NAME" = "rawr";`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeECPGCommand(tt.args.stmt); (err != nil) != tt.wantErr {
				t.Errorf("executeECPGCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pullErrorData(t *testing.T) {
	type args struct {
		out string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple string",
			args: args{
				out: "stdin:1: ERROR: syntax error at or near \";\"",
			},
			want: "ERROR: syntax error at or near \";\"",
		},
		{
			name: "random (missing stdin:1:) string",
			args: args{
				out: "the brown dog jumps over the red fox",
			},
			want: "",
		},
		{
			name: "empty string",
			args: args{
				out: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pullErrorData(tt.args.out); got != tt.want {
				t.Errorf("pullErrorData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewECPG(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ECPG
		wantErr bool
	}{
		{
			name: "simple ECGP creation",
			args: args{
				c: &Config{
					AddSemiColon:      true,
					TrimWhiteSpace:    true,
					QuestionMarks:     true,
					DependencyChecker: func(s string) bool { return true },
				},
			},
			want: &ECPG{
				config: &Config{
					AddSemiColon:      true,
					TrimWhiteSpace:    true,
					QuestionMarks:     true,
					DependencyChecker: func(s string) bool { return true },
				},
			},
			wantErr: false,
		},
		{
			name: "simple ECGP creation",
			args: args{
				c: &Config{
					AddSemiColon:      true,
					TrimWhiteSpace:    true,
					QuestionMarks:     true,
					DependencyChecker: func(s string) bool { return false },
				},
			},
			want: &ECPG{
				config: &Config{
					AddSemiColon:      true,
					TrimWhiteSpace:    true,
					QuestionMarks:     true,
					DependencyChecker: func(s string) bool { return false },
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewECPG(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewECPG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr {
				return
			}

			if got.config.AddSemiColon != tt.want.config.AddSemiColon ||
				got.config.QuestionMarks != tt.want.config.QuestionMarks ||
				got.config.TrimWhiteSpace != tt.want.config.TrimWhiteSpace {
				t.Errorf("NewECPG() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestECPG_CheckStatement(t *testing.T) {
	type fields struct {
		config *Config
	}
	type args struct {
		stmt string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "simple inclusive case",
			fields: fields{
				config: &Config{
					AddSemiColon:   true,
					TrimWhiteSpace: true,
					QuestionMarks:  true,
				},
			},
			args: args{
				stmt: "SELECT name, value, type FROM MySchema.MyTable WHERE id = ?",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ECPG{
				config: tt.fields.config,
			}
			if err := e.CheckStatement(tt.args.stmt); (err != nil) != tt.wantErr {
				t.Errorf("ECPG.CheckStatement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
