package base

import "context"

const InstanceIDKey = "iid-key"

func WithInstanceID(id string, ctx context.Context) context.Context {
	return context.WithValue(ctx, InstanceIDKey, id)
}

func GetInstanceID(ctx context.Context) string {
	v := ctx.Value(InstanceIDKey)
	res, _ := v.(string)
	return res
}

const ServiceNameKey = "sname-key"

func WithServiceName(name string, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceNameKey, name)
}

func GetServiceName(ctx context.Context) string {
	v := ctx.Value(ServiceNameKey)
	res, _ := v.(string)
	return res
}

const TraceIDKey = "tid-key"

func WithTraceID(id string, ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceIDKey, id)
}

func GetTraceID(ctx context.Context) string {
	v := ctx.Value(TraceIDKey)
	res, _ := v.(string)
	return res
}
