#!/bin/bash

echo
echo "As process manager"
dist/hivemind --exit-with-highest-exit-code - <<PROCFILE
job1: echo "job1 is running"; sleep 0.5; exit 13; echo "Done"
job2: echo "job2 is running"; sleep 2; echo "Done"
PROCFILE

echo
echo "As Job Runner, all completing"
dist/hivemind --as-job-runner --exit-with-highest-exit-code - <<PROCFILE
job1: echo "short job1 is running"; sleep 0.5; echo "Done"
job2: echo "short job2 is running"; sleep 1.0; echo "Done"
job3: echo "long job3 is running"; sleep 2; echo "Done"
PROCFILE

echo
echo "As Job Runner, one exits and fail fast"
dist/hivemind  --as-job-runner --exit-with-highest-exit-code - <<PROCFILE
job1: echo "job1 is running"; sleep 0.5; echo "Fail"; exit 13
job2: echo "job2 is running"; sleep 0.5; echo 1; sleep 0.5; echo 2; sleep 0.5; echo 3; sleep 1; echo "Done"
PROCFILE
