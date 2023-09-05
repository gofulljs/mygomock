package redis

import (
	"context"
	"mygomock/db/redis/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// go get github.com/golang/mock/mockgen
// go get github.com/golang/mock/mockgen/model
// 生成mock相关数据
// mockgen -destination=db/redis/mocks/mock_redis_cmdable.gen.go -package=mocks github.com/redis/go-redis/v9 Cmdable
func TestSet(t *testing.T) {
	type args struct {
		key        string
		mock       func(ctrl *gomock.Controller, key, value string, expiration time.Duration) (redis.Cmdable, *redis.StatusCmd)
		val        string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "set",
			args: args{
				key: "key1",
				val: "value1",
				mock: func(ctrl *gomock.Controller, key, value string, expiration time.Duration) (redis.Cmdable, *redis.StatusCmd) {
					cmd := mocks.NewMockCmdable(ctrl)
					status := redis.NewStatusCmd(context.Background())
					status.SetVal("OK")
					cmd.EXPECT().Set(context.Background(), key, value, expiration).Return(status)
					return cmd, status
				},
				expiration: time.Second,
			},
		},
		{
			name: "timeout",
			args: args{
				key: "key2",
				val: "value2",
				mock: func(ctrl *gomock.Controller, key, value string, expiration time.Duration) (redis.Cmdable, *redis.StatusCmd) {
					cmd := mocks.NewMockCmdable(ctrl)
					status := redis.NewStatusCmd(context.Background())
					status.SetErr(context.DeadlineExceeded)
					cmd.EXPECT().Set(context.Background(), key, value, expiration).Return(status)
					return cmd, status
				},
				expiration: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd, wantRes := tt.args.mock(ctrl, tt.args.key, tt.args.val, tt.args.expiration)
			status := cmd.Set(context.Background(), tt.args.key, tt.args.val, tt.args.expiration)
			assert.Equal(t, wantRes, status)
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		key  string
		val  string
		mock func(ctrl *gomock.Controller, key, value string) (redis.Cmdable, *redis.StringCmd)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "set",
			args: args{
				key: "key1",
				val: "value1",
				mock: func(ctrl *gomock.Controller, key, value string) (redis.Cmdable, *redis.StringCmd) {
					cmd := mocks.NewMockCmdable(ctrl)
					status := redis.NewStringCmd(context.Background())
					status.SetVal(value)
					cmd.EXPECT().Get(context.Background(), key).Return(status)
					return cmd, status
				},
			},
		},
		{
			name: "timeout",
			args: args{
				key: "key2",
				mock: func(ctrl *gomock.Controller, key, value string) (redis.Cmdable, *redis.StringCmd) {
					cmd := mocks.NewMockCmdable(ctrl)
					status := redis.NewStringCmd(context.Background())
					status.SetErr(context.DeadlineExceeded)
					cmd.EXPECT().Get(context.Background(), key).Return(status)
					return cmd, status
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd, wantRes := tt.args.mock(ctrl, tt.args.key, tt.args.val)
			status := cmd.Get(context.Background(), tt.args.key)
			assert.Equal(t, wantRes, status)
		})
	}
}
