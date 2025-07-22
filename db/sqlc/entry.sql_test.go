package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueries_CreateEntry(t *testing.T) {
	ctx := context.Background()
	acc, _ := createRandomAccount(t)

	type args struct {
		ctx context.Context
		arg CreateEntryParams
	}
	tests := []struct {
		name    string
		q       *Queries
		args    args
		check   func(t *testing.T, got Entry)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive_amount",
			q:    testQueries,
			args: args{ctx: ctx, arg: CreateEntryParams{acc.ID, acc.Balance}},
			check: func(t *testing.T, got Entry) {
				require.NotZero(t, got.ID, "o ID da entry não deve ser zero")
				require.NotZero(t, got.AccountID, "AccountID não deve ser zero")
				require.Equal(t, acc.Balance, got.Amount)
				require.WithinDuration(t, time.Now(), got.CreatedAt, time.Second*2)
			},
			wantErr: false,
		},
		{
			name: "zero_amount",
			q:    testQueries,
			args: args{ctx: ctx, arg: CreateEntryParams{acc.ID, 0}},
			check: func(t *testing.T, got Entry) {
				require.NotZero(t, got.ID)
				require.Equal(t, acc.ID, got.AccountID)
				require.Equal(t, int64(0), got.Amount)
			},
			wantErr: false,
		},
		{
			name:    "negative_amount_should_error",
			q:       testQueries,
			args:    args{ctx: ctx, arg: CreateEntryParams{acc.ID, -100}},
			check:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.q.db,
			}
			got, err := q.CreateEntry(tt.args.ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestQueries_GetEntry(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Entry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.GetEntry(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEntry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueries_ListEntries(t *testing.T) {
	type fields struct {
		db DBTX
	}
	type args struct {
		ctx context.Context
		arg ListEntriesParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Entry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				db: tt.fields.db,
			}
			got, err := q.ListEntries(tt.args.ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListEntries() got = %v, want %v", got, tt.want)
			}
		})
	}
}
