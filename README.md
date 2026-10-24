# temporal-go-experiments

My first hands-on run with [Temporal](https://temporal.io) in Go: a hello-world Workflow changed to test retries and worker crashes.

Started from [`temporalio/samples-go/helloworld`](https://github.com/temporalio/samples-go/tree/main/helloworld) (MIT, see `LICENSE`).

## What I changed

- `helloworld/helloworld.go`: the Activity fails on attempts 1 and 2 and succeeds on attempt 3.
- The Workflow sleeps 30s (`workflow.Sleep`) after the Activity, so the worker can be killed mid-run.
- `scripts/crash-test.sh`: kills the worker during that sleep, waits 10s, starts a fresh worker.

## Run it

```bash
brew install go temporal
temporal server start-dev            # terminal 1, UI at http://localhost:8233
go run ./helloworld/worker           # terminal 2
go run ./helloworld/starter          # terminal 3
./scripts/crash-test.sh              # crash test (stop the terminal-2 worker first)
```

## What I observed

- In the event history - it only recorded the final attempt and the last failure, not each of the failed attempts. To find all the failures I had to look at the logs of the worker.
- When I ran temporal I expected it to note there are no workers initialized but instead it said the workflow was running even though nothing was actually picking it up and running it.
- I didn't set any retry policy and it looks like the default is backoff at 1s. Meaning attempt 2 waited 1s after attempt 1 failed and attempt 3 waited 2s
- I expected in the crash test that the new worker that will initialize after the previous one was killed to pick up the logs of the crashed worker (what i thought the point of temporal is). It actually did not because the sleep was running on the server and not the worker and it actually finished so to worker 2 worker 1 actually finished. Worker 2 rebuilt the Workflow’s state by replaying this history. It didn’t rerun the Activity, and replay doesn’t print logs, so its only new log line was ‘completed’
