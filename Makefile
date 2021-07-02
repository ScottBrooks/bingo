deps:
	go get github.com/cortesi/modd/cmd/modd
	go get github.com/evanw/esbuild/cmd/esbuild
	curl -O https://unpkg.com/stimulus/dist/stimulus.umd.js 

docker:
	docker build -t bingo -t localhost:5000/bingo .

push-local: docker
	docker push localhost:5000/bingo:latest

deploy:
	kubectl apply -f k8s/deployment.yaml
