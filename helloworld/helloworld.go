package helloworld

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	// Crash test: kill the worker during this durable timer, then start a new one.
	logger.Info("Sleeping 30s; kill the worker now")
	if err := workflow.Sleep(ctx, 30*time.Second); err != nil {
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)

	// Retry experiment: fail on attempts 1 and 2, succeed on attempt 3.
	for activity.GetInfo(ctx).Attempt < 3 {
		return "", fmt.Errorf("error due to attempt being set to: %v", activity.GetInfo(ctx).Attempt)
	}
	return "It worked this time", nil
}
