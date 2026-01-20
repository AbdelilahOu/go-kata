package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type Option func(*UserAggregator)

type UserAggregator struct {
	timeout time.Duration
	logger  *slog.Logger
}

func WithTimeout(t time.Duration) Option {
	return func(a *UserAggregator) {
		a.timeout = t
	}
}

func WithLogger() Option {
	return func(a *UserAggregator) {
		a.logger = slog.Default()
	}
}

func NewUserAggregator(options ...Option) *UserAggregator {
	userAggr := &UserAggregator{}

	for _, opt := range options {
		opt(userAggr)
	}

	return userAggr
}

func (a *UserAggregator) Aggregate(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	var profile, orders string
	var err error
	g.Go(func() error { profile, err = fetchProfile(ctx, id); return err })
	g.Go(func() error { orders, err = fetchOrder(ctx, id); return err })

	if err := g.Wait(); err != nil {
		a.logger.Error(err.Error())
		return "", err
	}

	a.logger.Info("fetched seccess")
	return fmt.Sprintf("%s | %s", profile, orders), nil
}

func fetchProfile(ctx context.Context, id int) (string, error) {
	fmt.Printf("fetching profile %x \n", id)
	select {
	case <-time.After(2 * time.Second):
		return "Name: Alice", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}

}

func fetchOrder(ctx context.Context, id int) (string, error) {
	fmt.Printf("fetching order %x \n", id)
	select {
	case <-time.After(500 * time.Millisecond):
		return "Order: 5", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}

}

func main() {
	a := NewUserAggregator(WithLogger(), WithTimeout(1*time.Second))
	result, err := a.Aggregate(context.Background(), 1)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	fmt.Print(result)
}
