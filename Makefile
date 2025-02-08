# Run the app
.PHONY : run-app
run-app : 
	go run cmd/chronam-headliner/main.go

# Send the app a test URL
.PHONY : test-url
test-url :
	curl -X POST http://localhost:8050/chronamurl -H "Content-Type: application/json" \
	-d '{"url": "https://chroniclingamerica.loc.gov/lccn/sn83045462/1925-01-25/ed-1/seq-2/"}'
