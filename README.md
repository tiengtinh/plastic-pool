go build && PORT=5000 ./plastic-pool -max_workers 5

for i in {1..15}; do curl localhost:5000/work -d name=$USER$i -d delay=$(expr $i % 9 + 1)s; done